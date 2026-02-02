package route

import (
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type RouteHandler interface {
	GetRoutes(c *fiber.Ctx) error
}

type RouteHandlerImpl struct {
}

func NewRouteHandlerImpl() *RouteHandlerImpl {
	return &RouteHandlerImpl{}
}

func (h *RouteHandlerImpl) GetRoutes(c *fiber.Ctx) error {
	return c.JSON(model.NewSuccess("success", []string{}))
}
