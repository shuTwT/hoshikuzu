package couponusage

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	couponusage_service "github.com/shuTwT/hoshikuzu/internal/services/mall/couponusage"

	"github.com/gofiber/fiber/v2"
)

type CouponUsageHandler interface {
	ListCouponUsages(c *fiber.Ctx) error
	ListCouponUsagesPage(c *fiber.Ctx) error
	CreateCouponUsage(c *fiber.Ctx) error
	UpdateCouponUsage(c *fiber.Ctx) error
	QueryCouponUsage(c *fiber.Ctx) error
	DeleteCouponUsage(c *fiber.Ctx) error
	BatchUpdateCouponUsages(c *fiber.Ctx) error
	BatchDeleteCouponUsages(c *fiber.Ctx) error
	SearchCouponUsages(c *fiber.Ctx) error
}

type CouponUsageHandlerImpl struct {
	couponUsageService couponusage_service.CouponUsageService
}

func NewCouponUsageHandlerImpl(couponUsageService couponusage_service.CouponUsageService) *CouponUsageHandlerImpl {
	return &CouponUsageHandlerImpl{
		couponUsageService: couponUsageService,
	}
}

// @Summary 获取所有优惠券使用记录
// @Description 获取所有优惠券使用记录列表
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.CouponUsage}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/list [get]
func (h *CouponUsageHandlerImpl) ListCouponUsages(c *fiber.Ctx) error {
	usages, err := h.couponUsageService.ListAllCouponUsages(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", usages))
}

// @Summary 获取优惠券使用记录列表分页
// @Description 获取优惠券使用记录列表分页
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.CouponUsage]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/page [get]
func (h *CouponUsageHandlerImpl) ListCouponUsagesPage(c *fiber.Ctx) error {
	var pageQuery model.PageQuery

	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if pageQuery.Page < 1 {
		pageQuery.Page = 1
	}
	if pageQuery.Size < 1 {
		pageQuery.Size = 10
	}

	usages, total, err := h.couponUsageService.ListCouponUsages(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.CouponUsage]{
		Total:   int64(total),
		Records: usages,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建优惠券使用记录
// @Description 创建优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param req body model.CouponUsageCreateReq true "优惠券使用记录创建请求"
// @Success 200 {object} model.HttpSuccess{data=ent.CouponUsage}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/create [post]
func (h *CouponUsageHandlerImpl) CreateCouponUsage(c *fiber.Ctx) error {
	var req *model.CouponUsageCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	usage, err := h.couponUsageService.CreateCouponUsage(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", usage))
}

// @Summary 更新优惠券使用记录
// @Description 更新优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param id path int true "优惠券使用记录 ID"
// @Param req body model.CouponUsageUpdateReq true "优惠券使用记录更新请求"
// @Success 200 {object} model.HttpSuccess{data=ent.CouponUsage}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/update/{id} [put]
func (h *CouponUsageHandlerImpl) UpdateCouponUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req *model.CouponUsageUpdateReq
	if err = c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	usage, err := h.couponUsageService.UpdateCouponUsage(c.Context(), id, req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", usage))
}

// @Summary 查询优惠券使用记录
// @Description 查询优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param id path int true "优惠券使用记录 ID"
// @Success 200 {object} model.HttpSuccess{data=ent.CouponUsage}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/query/{id} [get]
func (h *CouponUsageHandlerImpl) QueryCouponUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	usage, err := h.couponUsageService.GetCouponUsage(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", usage))
}

// @Summary 删除优惠券使用记录
// @Description 删除优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param id path int true "优惠券使用记录 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/delete/{id} [delete]
func (h *CouponUsageHandlerImpl) DeleteCouponUsage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	if err := h.couponUsageService.DeleteCouponUsage(c.Context(), id); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量更新优惠券使用记录
// @Description 批量更新优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param req body model.CouponUsageBatchUpdateReq true "优惠券使用记录批量更新请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/batch/update [put]
func (h *CouponUsageHandlerImpl) BatchUpdateCouponUsages(c *fiber.Ctx) error {
	var req *model.CouponUsageBatchUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.couponUsageService.BatchUpdateCouponUsages(c.Context(), req.IDs, req); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量删除优惠券使用记录
// @Description 批量删除优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param req body model.CouponUsageBatchDeleteReq true "优惠券使用记录批量删除请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/batch/delete [delete]
func (h *CouponUsageHandlerImpl) BatchDeleteCouponUsages(c *fiber.Ctx) error {
	var req *model.CouponUsageBatchDeleteReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.couponUsageService.BatchDeleteCouponUsages(c.Context(), req.IDs); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 搜索优惠券使用记录
// @Description 搜索优惠券使用记录
// @Tags 优惠券使用记录
// @Accept json
// @Produce json
// @Param req query model.CouponUsageSearchReq true "优惠券使用记录搜索请求"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.CouponUsage]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/couponusage/search [get]
func (h *CouponUsageHandlerImpl) SearchCouponUsages(c *fiber.Ctx) error {
	var req model.CouponUsageSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	usages, total, err := h.couponUsageService.SearchCouponUsages(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.CouponUsage]{
		Total:   int64(total),
		Records: usages,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}
