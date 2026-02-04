package common

import (
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	common_service "github.com/shuTwT/hoshikuzu/internal/services/system/common"

	"github.com/gofiber/fiber/v2"
)

type CommonHandler interface {
	GetHomeStatistics(c *fiber.Ctx) error
}

type CommonHandlerImpl struct {
	commonService common_service.CommonService
}

func NewCommonHandlerImpl(commonService common_service.CommonService) *CommonHandlerImpl {
	return &CommonHandlerImpl{
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
func (h *CommonHandlerImpl) GetHomeStatistics(c *fiber.Ctx) error {
	return c.JSON(model.NewSuccess("统计信息获取成功", h.commonService.GetHomeStatistic(c.Context())))
}
