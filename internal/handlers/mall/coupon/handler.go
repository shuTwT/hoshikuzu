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

func (h *CouponHandlerImpl) ListCoupons(c *fiber.Ctx) error {
	coupons, err := h.couponService.ListAllCoupons(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", coupons))
}

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
