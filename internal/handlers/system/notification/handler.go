package notification

import (
	"strconv"

	"github.com/shuTwT/hoshikuzu/ent"
	notification_service "github.com/shuTwT/hoshikuzu/internal/services/system/notification"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler interface {
	ListNotificationPage(c *fiber.Ctx) error
	QueryNotification(c *fiber.Ctx) error
	DeleteNotification(c *fiber.Ctx) error
	BatchMarkAsRead(c *fiber.Ctx) error
}

type NotificationHandlerImpl struct {
	notificationService notification_service.NotificationService
}

func NewNotificationHandlerImpl(notificationService notification_service.NotificationService) *NotificationHandlerImpl {
	return &NotificationHandlerImpl{
		notificationService: notificationService,
	}
}

// @Summary 查询通知分页列表
// @Description 查询通知的分页列表，支持按已读/未读状态过滤
// @Tags 后台管理接口/通知
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param is_read query bool false "是否已读" default(false)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.Notification]}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/notifications/page [get]
func (h *NotificationHandlerImpl) ListNotificationPage(c *fiber.Ctx) error {
	var pageQuery = model.NotificationPageQuery{}
	err := c.QueryParser(&pageQuery)

	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	notifications, count, err := h.notificationService.QueryNotificationPage(c, pageQuery)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	pageResult := model.PageResult[*ent.Notification]{
		Total:   int64(count),
		Records: notifications,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 查询通知详情
// @Description 查询指定通知的详细信息
// @Tags 后台管理接口/通知
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} model.HttpSuccess{data=ent.Notification}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/notifications/query/:id [get]
func (h *NotificationHandlerImpl) QueryNotification(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	notification, err := h.notificationService.QueryNotification(c, id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", notification))
}

// @Summary 删除通知
// @Description 删除指定通知
// @Tags 后台管理接口/通知
// @Accept json
// @Produce json
// @Param id path int true "通知ID"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/notifications/delete/:id [delete]
func (h *NotificationHandlerImpl) DeleteNotification(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}
	err = h.notificationService.DeleteNotification(c, id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", nil))
}

// @Summary 批量标记为已读
// @Description 批量将通知标记为已读
// @Tags 后台管理接口/通知
// @Accept json
// @Produce json
// @Param request body model.NotificationBatchReadReq true "批量已读请求"
// @Success 200 {object} model.HttpSuccess{data=nil}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/notifications/batch/read [post]
func (h *NotificationHandlerImpl) BatchMarkAsRead(c *fiber.Ctx) error {
	var req model.NotificationBatchReadReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			err.Error(),
		))
	}

	err := h.notificationService.BatchMarkAsRead(c, req.IDs)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
	}
	return c.JSON(model.NewSuccess("success", nil))
}
