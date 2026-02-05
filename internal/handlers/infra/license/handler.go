package license

import (
	"strconv"
	"time"

	"github.com/shuTwT/hoshikuzu/internal/services/infra/license"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
	"github.com/shuTwT/hoshikuzu/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type LicenseHandler interface {
	ListLicensePage(c *fiber.Ctx) error
	QueryLicense(c *fiber.Ctx) error
	CreateLicense(c *fiber.Ctx) error
	UpdateLicense(c *fiber.Ctx) error
	DeleteLicense(c *fiber.Ctx) error
	VerifyLicense(c *fiber.Ctx) error
}

type LicenseHandlerImpl struct {
	licenseService license.LicenseService
}

func NewLicenseHandlerImpl(licenseService license.LicenseService) *LicenseHandlerImpl {
	return &LicenseHandlerImpl{licenseService: licenseService}
}

// @Summary 获取授权分页列表
// @Description 获取所有授权的分页列表
// @Tags 后台管理接口/授权
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[model.LicenseResp]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/page [get]
func (h *LicenseHandlerImpl) ListLicensePage(c *fiber.Ctx) error {
	pageQuery := model.PageQuery{}
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	count, licenses, err := h.licenseService.ListLicensePage(c.Context(), pageQuery.Page, pageQuery.Size)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	licenseResps := make([]*model.LicenseResp, 0, len(licenses))
	for _, licenseEntity := range licenses {
		licenseResps = append(licenseResps, &model.LicenseResp{
			ID:           licenseEntity.ID,
			CreatedAt:    (model.LocalTime)(licenseEntity.CreatedAt),
			UpdatedAt:    (model.LocalTime)(licenseEntity.UpdatedAt),
			Domain:       licenseEntity.Domain,
			LicenseKey:   licenseEntity.LicenseKey,
			CustomerName: licenseEntity.CustomerName,
			ExpireDate:   (model.LocalTime)(licenseEntity.ExpireDate),
			Status:       licenseEntity.Status,
		})
	}

	pageResult := model.PageResult[*model.LicenseResp]{
		Total:   int64(count),
		Records: licenseResps,
	}
	return c.JSON(model.NewSuccess("授权列表获取成功", pageResult))
}

// @Summary 查询授权
// @Description 查询指定授权的信息
// @Tags 后台管理接口/授权
// @Accept json
// @Produce json
// @Param id path int true "授权ID"
// @Success 200 {object} model.HttpSuccess{data=model.LicenseResp}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/query/{id} [get]
func (h *LicenseHandlerImpl) QueryLicense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	licenseEntity, err := h.licenseService.QueryLicense(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	licenseResp := &model.LicenseResp{
		ID:           licenseEntity.ID,
		CreatedAt:    (model.LocalTime)(licenseEntity.CreatedAt),
		UpdatedAt:    (model.LocalTime)(licenseEntity.UpdatedAt),
		Domain:       licenseEntity.Domain,
		LicenseKey:   licenseEntity.LicenseKey,
		CustomerName: licenseEntity.CustomerName,
		ExpireDate:   (model.LocalTime)(licenseEntity.ExpireDate),
		Status:       licenseEntity.Status,
	}
	return c.JSON(model.NewSuccess("授权查询成功", licenseResp))
}

// @Summary 创建授权
// @Description 创建一个新的授权
// @Tags 后台管理接口/授权
// @Accept json
// @Produce json
// @Param req body model.LicenseCreateReq true "授权创建请求"
// @Success 200 {object} model.HttpSuccess{data=model.LicenseResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/create [post]
func (h *LicenseHandlerImpl) CreateLicense(c *fiber.Ctx) error {
	var req *model.LicenseCreateReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	licenseKey, err := utils.GenerateLicenseKey(req.Domain)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	newLicense, err := h.licenseService.CreateLicense(c.Context(), req.Domain, licenseKey, req.CustomerName, req.ExpireDate)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	licenseResp := &model.LicenseResp{
		ID:           newLicense.ID,
		CreatedAt:    (model.LocalTime)(newLicense.CreatedAt),
		UpdatedAt:    (model.LocalTime)(newLicense.UpdatedAt),
		Domain:       newLicense.Domain,
		LicenseKey:   newLicense.LicenseKey,
		CustomerName: newLicense.CustomerName,
		ExpireDate:   (model.LocalTime)(newLicense.ExpireDate),
		Status:       newLicense.Status,
	}
	return c.JSON(model.NewSuccess("授权创建成功", licenseResp))
}

// @Summary 更新授权
// @Description 更新指定授权的信息
// @Tags 后台管理接口/授权
// @Accept json
// @Produce json
// @Param id path int true "授权ID"
// @Param req body model.LicenseUpdateReq true "授权更新请求"
// @Success 200 {object} model.HttpSuccess{data=model.LicenseResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/update/{id} [put]
func (h *LicenseHandlerImpl) UpdateLicense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var req *model.LicenseUpdateReq
	if err = c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedLicense, err := h.licenseService.UpdateLicense(c.Context(), id, req.Domain, req.LicenseKey, req.CustomerName, req.ExpireDate, req.Status)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	licenseResp := &model.LicenseResp{
		ID:           updatedLicense.ID,
		CreatedAt:    (model.LocalTime)(updatedLicense.CreatedAt),
		UpdatedAt:    (model.LocalTime)(updatedLicense.UpdatedAt),
		Domain:       updatedLicense.Domain,
		LicenseKey:   updatedLicense.LicenseKey,
		CustomerName: updatedLicense.CustomerName,
		ExpireDate:   (model.LocalTime)(updatedLicense.ExpireDate),
		Status:       updatedLicense.Status,
	}
	return c.JSON(model.NewSuccess("授权更新成功", licenseResp))
}

// @Summary 删除授权
// @Description 删除指定授权
// @Tags 后台管理接口/授权
// @Accept json
// @Produce json
// @Param id path int true "授权ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/delete/{id} [delete]
func (h *LicenseHandlerImpl) DeleteLicense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	err = h.licenseService.DeleteLicense(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("授权删除成功", nil))
}

// @Summary 验证授权
// @Description 验证指定域名的授权是否有效
// @Tags 公开接口/授权
// @Accept json
// @Produce json
// @Param req body model.LicenseVerifyReq true "授权验证请求"
// @Success 200 {object} model.HttpSuccess{data=model.LicenseVerifyResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/license/verify [post]
func (h *LicenseHandlerImpl) VerifyLicense(c *fiber.Ctx) error {
	var req *model.LicenseVerifyReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	licenseEntity, err := h.licenseService.VerifyLicense(c.Context(), req.Domain)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if licenseEntity == nil {
		return c.JSON(model.NewSuccess("授权验证", &model.LicenseVerifyResp{
			Valid:   false,
			Message: "未找到有效的授权信息",
		}))
	}

	return c.JSON(model.NewSuccess("授权验证", &model.LicenseVerifyResp{
		Valid:        true,
		CustomerName: licenseEntity.CustomerName,
		ExpireDate:   licenseEntity.ExpireDate.Format(time.RFC3339),
		Message:      "授权有效",
	}))
}
