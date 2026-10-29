import { usePermissionStore, usePermissionStoreHook } from "@/stores/modules/permission";
import type { RouteComponent, Router,RouteRecordRaw } from "vue-router";
import { cloneDeep, intersection, isEmpty } from "lodash-es"
import router from ".";
import { apiClient, useApi } from "@/api";
import { useStorageLocal } from "@/utils/utils";
import { userKey, type DataInfo } from "@/utils/auth";
import { buildHierarchyTree } from "@/utils/tree";
import { loadPluginRoutes } from "@/plugins/pluginLoader";

const modulesRoutes = import.meta.glob("/src/views/**/*.{vue,tsx}");

function isAllEmpty(value:any){
  return value === undefined || value === null || value === ''
}

function handRank(routeInfo: any) {
  const { name, path, parentId, meta } = routeInfo;
  return isAllEmpty(parentId)
    ? isAllEmpty(meta?.rank) ||
      (meta?.rank === 0 && name !== "Home" && path !== "/")
      ? true
      : false
    : false;
}

/** 按照路由中meta下的rank等级升序来排序路由 */
function ascending(arr: any[]) {
  arr.forEach((v, index) => {
    // 当rank不存在时，根据顺序自动创建，首页路由永远在第一位
    if (handRank(v)){
       v.meta.rank = index + 2;
      }
  });
  return arr.sort(
    (a: { meta: { rank: number } }, b: { meta: { rank: number } }) => {
      return a?.meta.rank - b?.meta.rank;
    }
  );
}

/** 过滤meta中showLink为false的菜单 */
function filterTree(data: RouteComponent[]) {
  const newTree = (cloneDeep(data) as any).filter(
    (v: { meta: { showLink: boolean } }) => v.meta?.showLink !== false
  );
  newTree.forEach(
    (v: { children:any }) => v.children && (v.children = filterTree(v.children))
  );
  return newTree;
}

/**
 * 判断两个数组是否存在相同的基本数据类型元素
 * @param arr1 第一个数组
 * @param arr2 第二个数组
 * @returns 存在相同值则返回 true，否则 false
 */
function isOneOfArray(arr1?:any[], arr2?:any[]):boolean {
  return Array.isArray(arr1) && Array.isArray(arr2)?
    !isEmpty(intersection(arr1, arr2))
  : false
}


function handleAsyncRoutes(routeList: any[]){
  if(routeList.length === 0){
    usePermissionStoreHook().handleWholeMenus(routeList)
  }else{
    addAsyncRoutes(routeList as Array<RouteRecordRaw>)?.map((v)=>{
      // 防止重复添加路由
      if(
        router.options.routes[0].children?.findIndex(
          value =>value.path === v.path
        ) !== -1
      ){
        return;
      }else{
        router.options.routes[0].children?.push(v)

        ascending(router.options.routes[0].children)
        if(!router.hasRoute(v.name as string)){
          router.addRoute(v)
        }
        const flattenRoutes:any = router.getRoutes().find(n=>n.path==="/");
        router.addRoute(flattenRoutes)
      }
    })
    usePermissionStoreHook().handleWholeMenus(routeList)
  }
  addPathMatch()
}

/** 过滤children长度为0的的目录，当目录下没有菜单时，会过滤此目录，目录没有赋予roles权限，当目录下只要有一个菜单有显示权限，那么此目录就会显示 */
function filterChildrenTree(data: RouteComponent[]) {
  const newTree = cloneDeep(data).filter((v: any) => v?.children?.length !== 0);
  newTree.forEach(
    (v: any) => v.children && (v.children = filterChildrenTree(v.children))
  );
  return newTree;
}

/** 从localStorage里取出当前登录用户的角色roles，过滤无权限的菜单 */
function filterNoPermissionTree(data: RouteComponent[]) {
  const currentRoles =
    useStorageLocal().getItem<DataInfo<number>>(userKey)?.roles ?? [];
  const newTree = cloneDeep(data).filter((v: any) =>{
    return isOneOfArray(v.meta?.roles, currentRoles)
    }
  );

  newTree.forEach(
    (v: any) => v.children && (v.children = filterNoPermissionTree(v.children))
  );
  return filterChildrenTree(newTree);
}

/** 通过指定 `key` 获取父级路径集合，默认 `key` 为 `path` */
function getParentPaths(value: string, routes: RouteRecordRaw[], key = "path") {
  // 深度遍历查找
  function dfs(routes: RouteRecordRaw[], value: string, parents: string[]) {
    for (let i = 0; i < routes.length; i++) {
      const item = routes[i];
      // 返回父级path
      if ((item as any)[key] === value) return parents;
      // children不存在或为空则不递归
      if (!item.children || !item.children.length) continue;
      // 往下查找时将当前path入栈
      parents.push(item.path);

      if (dfs(item.children, value, parents).length) return parents;
      // 深度遍历查找未找到时当前path 出栈
      parents.pop();
    }
    // 未找到时返回空数组
    return [];
  }

  return dfs(routes, value, []);
}

/** 查找对应 `path` 的路由信息 */
function findRouteByPath(path: string, routes: RouteRecordRaw[]):RouteRecordRaw|undefined {
  let res = routes.find((item: { path: string }) => item.path == path);
  if (res) {
    return isProxy(res) ? toRaw(res) : res;
  } else {
    for (let i = 0; i < routes.length; i++) {
      const childRoutes = routes[i].children;
      if (
        childRoutes instanceof Array  &&  childRoutes.length > 0
      ) {
        res = findRouteByPath(path, childRoutes);
        if (res) {
          return isProxy(res) ? toRaw(res) : res;
        }
      }
    }
    return undefined;
  }
}

function addPathMatch() {
  if (!router.hasRoute("pathMatch")) {
    router.addRoute({
      path: "/:pathMatch(.*)",
      name: "pathMatch",
      redirect: "/error/404"
    });
  }
}

function initRouter():Promise<Router>{
  return new Promise(resolve=>{
    Promise.all([
      useApi(apiClient.api.v1RoutesList),
      loadPluginRoutes(),
    ]).then(([{data}, pluginRoutes])=>{
      handleAsyncRoutes(cloneDeep(data).concat(pluginRoutes as any))
      resolve(router)
    })
  })
}

function addAsyncRoutes(arrRoutes:Array<RouteRecordRaw>){
  if(!arrRoutes || !arrRoutes.length) return
  const moduleRoutesKeys = Object.keys(modulesRoutes);
  arrRoutes.forEach((v)=>{
    v.meta!.backstage = true
    if(v.children&&v.children.length&&!v.redirect){
      v.redirect = v.children[0].path
    }
    if(v.children&&v.children.length&&!v.name){
      v.name = (v.children[0].name as string)+"Parent"
    }
    if(v.children&&v.children.length){
      addAsyncRoutes(v.children)
    }
  })
  return arrRoutes;
}

/**
 * 将多级嵌套路由处理成一维数组
 * @param routesList 传入路由
 * @returns 返回处理后的一维路由
 */
function formatFlatteningRoutes(routesList: RouteRecordRaw[]) {
  if (routesList.length === 0) return routesList;
  let hierarchyList = buildHierarchyTree(routesList);
  for (let i = 0; i < hierarchyList.length; i++) {
    if (hierarchyList[i].children) {
      hierarchyList = hierarchyList
        .slice(0, i + 1)
        .concat(hierarchyList[i].children, hierarchyList.slice(i + 1));
    }
  }
  return hierarchyList;
}

/**
 * 一维数组处理成多级嵌套数组（三级及以上的路由全部拍成二级，keep-alive 只支持到二级缓存）
 * https://github.com/pure-admin/vue-pure-admin/issues/67
 * @param routesList 处理后的一维路由菜单数组
 * @returns 返回将一维数组重新处理成规定路由的格式
 */
function formatTwoStageRoutes(routesList: RouteRecordRaw[]) {
  if (routesList.length === 0) return routesList;
  const newRoutesList: RouteRecordRaw[] = [];
  routesList.forEach((v: RouteRecordRaw) => {
    if (v.path === "/") {
      newRoutesList.push({
        component: v.component,
        name: v.name,
        path: v.path,
        redirect: v.redirect,
        meta: v.meta,
        children: []
      });
    } else {
      newRoutesList[0]?.children?.push({ ...v });
    }
  });
  return newRoutesList;
}

export {
  ascending,
  filterTree,
  isOneOfArray,
  initRouter,
  addAsyncRoutes,
  filterNoPermissionTree,
  formatFlatteningRoutes,
  formatTwoStageRoutes,
  findRouteByPath,
  getParentPaths
}
