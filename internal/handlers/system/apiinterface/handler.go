package apiinterface

import (
	"github.com/shuTwT/hoshikuzu/ent"
	api_interface_service "github.com/shuTwT/hoshikuzu/internal/services/system/apiinterface"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type ApiInterfaceHandler struct {
	apiInterfaceService api_interface_service.ApiInterfaceService
}

func NewApiInterfaceHandler(apiInterfaceService api_interface_service.ApiInterfaceService) *ApiInterfaceHandler {
	return &ApiInterfaceHandler{
		apiInterfaceService: apiInterfaceService,
	}
}

// @Summary 查询API路由分页列表
// @Description 查询所有API路由的分页列表
// @Tags 后台管理接口/api路由
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.ApiRoute]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/api-interface/page [get]
func (h *ApiInterfaceHandler) ListApiRoutesPage(ctx *fiber.Ctx) error {
	var req model.PageQuery
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	count, apiRoutes, err := h.apiInterfaceService.ListApiRoutesPage(req.Page, req.Size)
	if err != nil {
		return ctx.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.ApiPerms]{
		Total:   int64(count),
		Records: apiRoutes,
	}

	return ctx.JSON(model.NewSuccess("获取路由列表成功", pageResult))
}
