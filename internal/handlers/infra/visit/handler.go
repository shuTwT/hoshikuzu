package visit

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/internal/services/infra/visit"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type VisitHandler interface {
	HandleVisitor(c *fiber.Ctx) error
	ListVisitLogPage(c *fiber.Ctx) error
	QueryVisitLog(c *fiber.Ctx) error
	DeleteVisitLog(c *fiber.Ctx) error
	BatchDeleteVisitLog(c *fiber.Ctx) error
}

type VisitHandlerImpl struct {
	visitService visit.VisitService
}

func NewVisitHandlerImpl(visitService visit.VisitService) VisitHandler {
	return &VisitHandlerImpl{visitService: visitService}
}

func (h *VisitHandlerImpl) HandleVisitor(c *fiber.Ctx) error {
	var req model.VisitLogReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	return h.visitService.CreateVisitLog(c.Context(), req)
}

func (h *VisitHandlerImpl) ListVisitLogPage(c *fiber.Ctx) error {
	var pageQuery model.VisitLogPageQuery
	err := c.QueryParser(&pageQuery)

	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	visitLogs, count, err := h.visitService.QueryVisitLogPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	pageResult := model.PageResult[*ent.VisitLog]{
		Total:   int64(count),
		Records: visitLogs,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

func (h *VisitHandlerImpl) QueryVisitLog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	visitLog, err := h.visitService.QueryVisitLog(c, id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", visitLog))
}

func (h *VisitHandlerImpl) DeleteVisitLog(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	err = h.visitService.DeleteVisitLog(c, id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

func (h *VisitHandlerImpl) BatchDeleteVisitLog(c *fiber.Ctx) error {
	var req model.VisitLogBatchDeleteReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	err := h.visitService.BatchDeleteVisitLog(c, req.IDs)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
