package route

import (
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type RouteHandler struct {
}

// @Summary 获取路由
// @Description 获取系统中所有路由的列表
// @Tags 后台管理接口/路由
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]string}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/routes [get]
func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

func (h *RouteHandler) GetRoutes(c *fiber.Ctx) error {
	return c.JSON(model.NewSuccess("success", []string{}))
}
