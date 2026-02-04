package coupon

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	coupon_service "github.com/shuTwT/hoshikuzu/internal/services/mall/coupon"

	"github.com/gofiber/fiber/v2"
)

type CouponHandler interface {
	ListCoupons(c *fiber.Ctx) error
	ListCouponsPage(c *fiber.Ctx) error
	CreateCoupon(c *fiber.Ctx) error
	UpdateCoupon(c *fiber.Ctx) error
	QueryCoupon(c *fiber.Ctx) error
	DeleteCoupon(c *fiber.Ctx) error
	BatchUpdateCoupons(c *fiber.Ctx) error
	BatchDeleteCoupons(c *fiber.Ctx) error
	SearchCoupons(c *fiber.Ctx) error
}

type CouponHandlerImpl struct {
	couponService coupon_service.CouponService
}

func NewCouponHandlerImpl(couponService coupon_service.CouponService) *CouponHandlerImpl {
	return &CouponHandlerImpl{
		couponService: couponService,
	}
}

// @Summary 获取所有优惠券
// @Description 获取所有优惠券列表
// @Tags 优惠券
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.Coupon}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/list [get]
func (h *CouponHandlerImpl) ListCoupons(c *fiber.Ctx) error {
	coupons, err := h.couponService.ListAllCoupons(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", coupons))
}

// @Summary 获取优惠券分页列表
// @Description 获取优惠券分页列表
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Coupon]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/page [get]
func (h *CouponHandlerImpl) ListCouponsPage(c *fiber.Ctx) error {
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

	coupons, total, err := h.couponService.ListCoupons(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.Coupon]{
		Total:   int64(total),
		Records: coupons,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 创建优惠券
// @Description 创建优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param coupon_create_req body model.CouponCreateReq true "优惠券创建请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.Coupon}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/create [post]
func (h *CouponHandlerImpl) CreateCoupon(c *fiber.Ctx) error {
	var req *model.CouponCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	coupon, err := h.couponService.CreateCoupon(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", coupon))
}

// @Summary 更新优惠券
// @Description 更新优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param id path int true "优惠券 ID"
// @Param coupon_update_req body model.CouponUpdateReq true "优惠券更新请求体"
// @Success 200 {object} model.HttpSuccess{data=ent.Coupon}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/update/{id} [put]
func (h *CouponHandlerImpl) UpdateCoupon(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req *model.CouponUpdateReq
	if err = c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	coupon, err := h.couponService.UpdateCoupon(c.Context(), id, req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", coupon))
}

// @Summary 查询优惠券
// @Description 查询优惠券详情
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param id path int true "优惠券 ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Coupon}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/query/{id} [get]
func (h *CouponHandlerImpl) QueryCoupon(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	coupon, err := h.couponService.GetCoupon(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", coupon))
}

// @Summary 删除优惠券
// @Description 删除优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param id path int true "优惠券 ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/delete/{id} [delete]
func (h *CouponHandlerImpl) DeleteCoupon(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	if err := h.couponService.DeleteCoupon(c.Context(), id); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量更新优惠券
// @Description 批量更新优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param coupon_batch_update_req body model.CouponBatchUpdateReq true "优惠券批量更新请求体"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/batch/update [put]
func (h *CouponHandlerImpl) BatchUpdateCoupons(c *fiber.Ctx) error {
	var req *model.CouponBatchUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.couponService.BatchUpdateCoupons(c.Context(), req.IDs, req); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量删除优惠券
// @Description 批量删除优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param coupon_batch_delete_req body model.CouponBatchDeleteReq true "优惠券批量删除请求体"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/batch/delete [delete]
func (h *CouponHandlerImpl) BatchDeleteCoupons(c *fiber.Ctx) error {
	var req *model.CouponBatchDeleteReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if err := h.couponService.BatchDeleteCoupons(c.Context(), req.IDs); err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 搜索优惠券
// @Description 搜索优惠券
// @Tags 优惠券
// @Accept json
// @Produce json
// @Param name query string false "优惠券名称"
// @Param status query int false "优惠券状态"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Coupon]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/coupon/search [get]
func (h *CouponHandlerImpl) SearchCoupons(c *fiber.Ctx) error {
	var req model.CouponSearchReq
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	coupons, total, err := h.couponService.SearchCoupons(c.Context(), req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	pageResult := model.PageResult[*ent.Coupon]{
		Total:   int64(total),
		Records: coupons,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}
