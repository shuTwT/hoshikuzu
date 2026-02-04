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

// @Summary 处理访客访问
// @Description 处理访客访问请求，记录访问日志
// @Tags 公开接口/访客访问
// @Accept json
// @Produce json
// @Param req body model.VisitLogReq true "访客访问请求"
// @Success 200 {object} model.HttpSuccess{data=ent.VisitLog}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/visit/handle [post]
func (h *VisitHandlerImpl) HandleVisitor(c *fiber.Ctx) error {
	var req model.VisitLogReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	return h.visitService.CreateVisitLog(c.Context(), req)
}

// @Summary 查询访客访问日志分页
// @Description 查询访客访问日志的分页列表
// @Tags 后台管理接口/访客访问
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.VisitLog]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/visit/list [get]
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

// @Summary 查询访客访问日志
// @Description 查询指定ID的访客访问日志
// @Tags 后台管理接口/访客访问
// @Accept json
// @Produce json
// @Param id path int true "Visit Log ID"
// @Success 200 {object} model.HttpSuccess{data=ent.VisitLog}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/visit/query/{id} [get]
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

// @Summary 删除访客访问日志
// @Description 删除指定ID的访客访问日志
// @Tags 后台管理接口/访客访问
// @Accept json
// @Produce json
// @Param id path int true "Visit Log ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/visit/delete/{id} [delete]
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

// @Summary 批量删除访客访问日志
// @Description 批量删除指定ID的访客访问日志
// @Tags 后台管理接口/访客访问
// @Accept json
// @Produce json
// @Param req body model.VisitLogBatchDeleteReq true "批量删除访客访问日志请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/visit/batch/delete [post]
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
