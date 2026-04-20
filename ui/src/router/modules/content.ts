import type { RouteRecordRaw } from 'vue-router'

export default [
    {
        path: "/content",
        name:"ContentManagement",
        redirect: "/content/post",
        meta: {
            title: '内容管理',
            showLink: true,
            roles: ['admin', 'common'],
            rank: 10
        },
        children: [
            {
                path: '/content/post',
                name: 'PostManagement',
                component: () => import('@/views/content/post/index.vue'),
                meta: {
                    title: '文章管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/post/edit',
                name: 'PostEditor',
                component: () => import('@/views/content/post/editor/index.vue'),
                meta: {
                    title: '文章编辑',
                    showLink: false,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/category',
                name: 'CategoryManagement',
                component: () => import('@/views/content/category/index.vue'),
                meta: {
                    title: '分类管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/tag',
                name: 'TagManagement',
                component: () => import('@/views/content/tag/index.vue'),
                meta: {
                    title: '标签管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/album',
                name: 'AlbumManagement',
                component: () => import('@/views/content/album/index.vue'),
                meta: {
                    title: '相册管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/comment',
                name: 'CommentManagement',
                component: () => import('@/views/content/comment/index.vue'),
                meta: {
                    title: '评论管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/flink',
                name: 'FlinkManagement',
                component: () => import('@/views/content/flink/index.vue'),
                meta: {
                    title: '友链管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/flink-application',
                name: 'FlinkApplicationManagement',
                component: () => import('@/views/content/flinkapplication/index.vue'),
                meta: {
                    title: '友链申请管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/friend-circle',
                name: 'FriendCircleRecord',
                component: () => import('@/views/content/friend-circle/index.vue'),
                meta: {
                    title: '朋友圈',
                    showLink: true,
                    roles: ['admin'],
                },
            },
            {
                path: '/content/essay',
                name: 'EssayManagement',
                component: () => import('@/views/content/essay/index.vue'),
                meta: {
                    title: '说说',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
            {
                path: '/content/menu',
                name: 'MenuManagement',
                component: () => import('@/views/content/menu/index.vue'),
                meta: {
                    title: '前台菜单管理',
                    showLink: true,
                    roles: ['admin', 'common'],
                },
            },
        ]
    }
] satisfies RouteRecordRaw[]