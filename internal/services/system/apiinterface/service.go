package apiinterface

import (
	"context"
	"log"
	"strings"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/apiperms"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type ApiInterfaceService interface {
	GetAppRoutes(app *fiber.App) []model.ApiRoute
	SyncRoutes(app *fiber.App)
	ListApiRoutesPage(page int, limit int) (int, []*ent.ApiPerms, error)
}

type ApiInterfaceServiceImpl struct {
	client *ent.Client
}

func NewApiInterfaceServiceImpl(client *ent.Client) *ApiInterfaceServiceImpl {
	return &ApiInterfaceServiceImpl{client: client}
}

func (s *ApiInterfaceServiceImpl) GetAppRoutes(app *fiber.App) []model.ApiRoute {
	routes := app.GetRoutes()
	// 筛选出 get,post,put,delete 方法的路由
	routeList := []model.ApiRoute{}
	for _, route := range routes {
		// 过滤掉不是 /api/v1 开头的路由
		if !strings.HasPrefix(route.Path, "/api/v1") || route.Name == "" {
			continue
		}
		if route.Method == "GET" || route.Method == "POST" || route.Method == "PUT" || route.Method == "DELETE" {
			routeList = append(routeList, model.ApiRoute{
				Name:           route.Name,
				Path:           route.Path,
				Method:         route.Method,
				Desc:           route.Path,
				PermissionType: model.PermissionType.Public,
				Roles:          []string{"default"},
				Status:         "active",
			})
		}
	}

	return routeList
}

/**
 * 启动时同步
 * @Description: 同步路由
 */
func (s *ApiInterfaceServiceImpl) SyncRoutes(app *fiber.App) {
	routeList := s.GetAppRoutes(app)
	for _, route := range routeList {
		// 检查路由是否已存在
		exists, err := s.client.ApiPerms.Query().Where(apiperms.Path(route.Path), apiperms.Method(route.Method)).Exist(context.Background())
		if err != nil {
			log.Printf("SyncRoutes error: %v", err)
			continue
		}
		if exists {
			continue
		}
		_, err = s.client.ApiPerms.Create().
			SetName(route.Name).
			SetPath(route.Path).
			SetMethod(route.Method).
			SetDesc(route.Desc).
			SetPermissionType(route.PermissionType).
			SetRoles(route.Roles).
			SetStatus(route.Status).
			Save(context.Background())
		if err != nil {
			log.Printf("SyncRoutes error: %v", err)
		}
	}
}

func (s *ApiInterfaceServiceImpl) ListApiRoutesPage(page int, limit int) (int, []*ent.ApiPerms, error) {

	offset := (page - 1) * limit

	// 统计总数
	count, err := s.client.ApiPerms.Query().Count(context.Background())
	if err != nil {
		return 0, nil, err
	}

	apiRoutes, err := s.client.ApiPerms.Query().
		Offset(offset).
		Limit(limit).
		All(context.Background())
	if err != nil {
		return 0, nil, err
	}

	return count, apiRoutes, nil
}
