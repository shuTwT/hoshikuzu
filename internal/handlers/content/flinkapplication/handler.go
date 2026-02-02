package flinkapplication

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	flinkapplication_service "github.com/shuTwT/hoshikuzu/internal/services/content/flinkapplication"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type FlinkApplicationHandler interface {
	CreateFlinkApplication(c *fiber.Ctx) error
	ListFlinkApplicationPage(c *fiber.Ctx) error
	QueryFlinkApplication(c *fiber.Ctx) error
	ApproveFlinkApplication(c *fiber.Ctx) error
}

type FlinkApplicationHandlerImpl struct {
	client                  *ent.Client
	flinkApplicationService flinkapplication_service.FlinkApplicationService
}

func NewFlinkApplicationHandlerImpl(client *ent.Client, flinkApplicationService flinkapplication_service.FlinkApplicationService) *FlinkApplicationHandlerImpl {
	return &FlinkApplicationHandlerImpl{
		client:                  client,
		flinkApplicationService: flinkApplicationService,
	}
}

func (h *FlinkApplicationHandlerImpl) CreateFlinkApplication(c *fiber.Ctx) error {
	var createReq *model.FlinkApplicationCreateReq
	if err := c.BodyParser(&createReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	application, err := h.flinkApplicationService.CreateFlinkApplication(c.Context(), createReq)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", application))
}

func (h *FlinkApplicationHandlerImpl) ListFlinkApplicationPage(c *fiber.Ctx) error {
	var pageQuery model.FlinkApplicationPageReq
	if err := c.QueryParser(&pageQuery); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	applications, count, err := h.flinkApplicationService.ListFlinkApplicationPage(c.Context(), pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	records := []model.FlinkApplicationResp{}
	for _, application := range applications {
		records = append(records, model.FlinkApplicationResp{
			ID:                 application.ID,
			CreatedAt:          &application.CreatedAt,
			UpdatedAt:          &application.UpdatedAt,
			WebsiteURL:         application.WebsiteURL,
			ApplicationType:    application.ApplicationType,
			WebsiteName:        application.WebsiteName,
			WebsiteLogo:        application.WebsiteLogo,
			WebsiteDescription: application.WebsiteDescription,
			ContactEmail:       application.ContactEmail,
			SnapshotURL:        application.SnapshotURL,
			OriginalWebsiteURL: application.OriginalWebsiteURL,
			ModificationReason: application.ModificationReason,
			Status:             application.Status,
			RejectReason:       application.RejectReason,
		})
	}
	pageResult := model.PageResult[model.FlinkApplicationResp]{
		Total:   int64(count),
		Records: records,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

func (h *FlinkApplicationHandlerImpl) QueryFlinkApplication(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	application, err := h.flinkApplicationService.QueryFlinkApplication(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", application))
}

func (h *FlinkApplicationHandlerImpl) ApproveFlinkApplication(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	var updateReq *model.FlinkApplicationUpdateReq
	if err = c.BodyParser(&updateReq); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	application, err := h.flinkApplicationService.ApproveFlinkApplication(c.Context(), id, updateReq.Status, updateReq.RejectReason)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", application))
}
