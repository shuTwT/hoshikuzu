package payorder_handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/payorder"
	payorder_service "github.com/shuTwT/hoshikuzu/internal/services/mall/payorder"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type PayOrderHandler interface {
	ListPayOrderPage(c *fiber.Ctx) error
	UpdatePayOrder(c *fiber.Ctx) error
	QueryPayOrder(c *fiber.Ctx) error
	DeletePayOrder(c *fiber.Ctx) error
	SubmitPayOrder(c *fiber.Ctx) error
}

type PayOrderHandlerImpl struct {
	client          *ent.Client
	payOrderService payorder_service.PayOrderService
}

func NewPayOrderHandlerImpl(client *ent.Client, service payorder_service.PayOrderService) *PayOrderHandlerImpl {
	return &PayOrderHandlerImpl{client: client, payOrderService: service}
}

// @Summary 获取支付订单列表
// @Description 获取所有支付订单的列表
// @Tags 支付订单
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=[]ent.PayOrder}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/list [get]
func (h *PayOrderHandlerImpl) ListPayOrderPage(c *fiber.Ctx) error {
	var req model.PageQuery
	if err := c.QueryParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	orders, count, err := h.payOrderService.ListPayOrderPage(c.Context(), &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	pageResult := model.PageResult[*ent.PayOrder]{
		Total:   int64(count),
		Records: orders,
	}
	return c.JSON(model.NewSuccess("success", pageResult))
}

// @Summary 更新支付订单
// @Description 更新指定支付订单的信息
// @Tags 支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Param payorder body ent.PayOrder true "支付订单信息"
// @Success 200 {object} ent.PayOrder
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/update/{id} [put]
func (h *PayOrderHandlerImpl) UpdatePayOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	var order *model.PayOrderUpdateReq
	if err = c.BodyParser(&order); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedOrder, err := h.client.PayOrder.UpdateOneID(id).
		SetChannelType(order.ChannelType).
		SetOrderID(order.OrderID).
		SetOutTradeNo(order.OutTradeNo).
		SetTotalFee(order.TotalFee).
		SetSubject(order.Subject).
		SetBody(order.Body).
		SetNotifyURL(order.NotifyURL).
		SetReturnURL(order.ReturnURL).
		SetExtra(order.Extra).
		SetPayURL(order.PayURL).
		SetState(order.State).
		SetErrorMsg(order.ErrorMsg).
		SetRaw(order.Raw).
		Save(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", updatedOrder))
}

// @Summary 查询支付订单
// @Description 查询指定支付订单的详细信息
// @Tags 支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Success 200 {object} ent.PayOrder
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/query/{id} [get]
func (h *PayOrderHandlerImpl) QueryPayOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	order, err := h.client.PayOrder.Query().
		Where(payorder.ID(id)).
		Only(c.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusNotFound,
				"PayOrder not found"))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", order))
}

// @Summary 删除支付订单
// @Description 删除指定支付订单
// @Tags 支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/delete/{id} [delete]
func (h *PayOrderHandlerImpl) DeletePayOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format",
		))
	}

	err = h.client.PayOrder.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.JSON(model.NewError(fiber.StatusBadRequest,
				"PayOrder not found",
			))
		}
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success",
		nil,
	))
}

// @Summary 提交支付订单
// @Description 提交一个新的支付订单
// @Tags 支付订单
// @Accept json
// @Produce json
// @Param payorder body model.PayOrderSubmitReq true "支付订单提交请求"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/submit [post]
func (h *PayOrderHandlerImpl) SubmitPayOrder(c *fiber.Ctx) error {
	var req model.PayOrderSubmitReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	switch req.OrderType {
	case model.PayOrderTypePost:
		if req.PostId <= 0 {
			return c.JSON(model.NewError(fiber.StatusBadRequest,
				"PostId is required",
			))
		}
	case model.PayOrderTypeProduct:
		if req.ProductId <= 0 {
			return c.JSON(model.NewError(fiber.StatusBadRequest,
				"ProductId is required",
			))
		}
	}
	err := h.payOrderService.SubmitPayOrder(c.Context(), &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return nil
}
