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

func (h *CouponUsageHandlerImpl) ListCouponUsages(c *fiber.Ctx) error {
	usages, err := h.couponUsageService.ListAllCouponUsages(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(model.NewSuccess("success", usages))
}

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
