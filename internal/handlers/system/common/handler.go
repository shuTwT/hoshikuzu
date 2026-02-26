package common

import (
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	common_service "github.com/shuTwT/hoshikuzu/internal/services/system/common"

	"github.com/gofiber/fiber/v2"
)

type CommonHandler struct {
	commonService common_service.CommonService
}

func NewCommonHandler(commonService common_service.CommonService) *CommonHandler {
	return &CommonHandler{
		commonService: commonService,
	}
}

// @Summary 获取首页统计信息
// @Description 获取首页统计信息，包括文章数量、评论数量、相册数量、照片数量
// @Tags 后台管理接口/统计
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=model.HomeStatistic}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/common/statistic [get]
func (h *CommonHandler) GetHomeStatistics(c *fiber.Ctx) error {
	return c.JSON(model.NewSuccess("统计信息获取成功", h.commonService.GetHomeStatistic(c.Context())))
}
