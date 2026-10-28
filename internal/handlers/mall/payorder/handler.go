package payorder_handler

import (
	"errors"
	"html"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/payorder"
	"github.com/shuTwT/hoshikuzu/internal/middleware"
	payorder_service "github.com/shuTwT/hoshikuzu/internal/services/mall/payorder"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type PayOrderHandler struct {
	client          *ent.Client
	payOrderService payorder_service.PayOrderService
}

func NewPayOrderHandler(client *ent.Client, service payorder_service.PayOrderService) *PayOrderHandler {
	return &PayOrderHandler{client: client, payOrderService: service}
}

// @Summary 获取支付订单列表
// @Description 获取所有支付订单的分页列表
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} model.HttpSuccess{data=model.PageResult[ent.PayOrder]}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/page [get]
func (h *PayOrderHandler) ListPayOrderPage(c *fiber.Ctx) error {
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
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Param payorder body model.PayOrderUpdateReq true "支付订单信息"
// @Success 200 {object} ent.PayOrder
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/update/{id} [put]
func (h *PayOrderHandler) UpdatePayOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest,
			"Invalid ID format"))
	}

	var order model.PayOrderUpdateReq
	if err = c.BodyParser(&order); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	updatedOrder, err := h.client.PayOrder.UpdateOneID(id).
		SetChannelType(order.ChannelType).
		SetOrderID(order.OrderID).
		SetOutTradeNo(order.OutTradeNo).
		SetOrderPrice(order.OrderPrice).
		SetPrice(order.Price).
		SetChannelFeePrice(order.ChannelFeePrice).
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
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Success 200 {object} ent.PayOrder
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/query/{id} [get]
func (h *PayOrderHandler) QueryPayOrder(c *fiber.Ctx) error {
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
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 404 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/delete/{id} [delete]
func (h *PayOrderHandler) DeletePayOrder(c *fiber.Ctx) error {
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
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param payorder body model.PayOrderSubmitReq true "支付订单提交请求"
// @Success 200 {object} model.HttpSuccess
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/submit [post]
func (h *PayOrderHandler) SubmitPayOrder(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(fiber.StatusUnauthorized, "请先登录"))
	}

	var req model.PayOrderSubmitReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if req.Money <= 0 {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "金额必须大于0"))
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
	default:
		return c.JSON(model.NewError(fiber.StatusBadRequest, "无效的订单类型"))
	}

	resp, err := h.payOrderService.SubmitPayOrder(c.Context(), loginUser.ID, &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 支付回调
// @Description 易支付异步通知回调，公开接口
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce plain
// @Success 200 {string} string
// @Router /api/v1/pay-order/notify [post]
func (h *PayOrderHandler) NotifyPayOrder(c *fiber.Ctx) error {
	params := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = string(v)
	})
	c.Context().PostArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = string(v)
	})

	if err := h.payOrderService.HandleNotify(c.Context(), params); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("fail")
	}
	return c.SendString("success")
}

// @Summary 查询订单状态
// @Description 主动同步易支付订单状态并返回
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param id path string true "支付订单ID"
// @Success 200 {object} model.HttpSuccess{data=model.PayOrderStatusResp}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/status/{id} [get]
func (h *PayOrderHandler) QueryOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	resp, err := h.payOrderService.SyncOrderStatus(c.Context(), id)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 余额充值
// @Description 提交余额充值订单，支付成功后自动入账钱包并发放会员积分
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param req body model.PayOrderRechargeReq true "充值请求"
// @Success 200 {object} model.HttpSuccess{data=model.PayOrderSubmitResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/recharge [post]
func (h *PayOrderHandler) RechargePayOrder(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(fiber.StatusUnauthorized, "请先登录"))
	}

	var req model.PayOrderRechargeReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}
	if req.Amount <= 0 {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "充值金额必须大于0"))
	}

	resp, err := h.payOrderService.RechargePayOrder(c.Context(), loginUser.ID, &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 订单退款
// @Description 对已支付订单发起全额退款（后台管理接口）
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Param id path int true "支付订单ID"
// @Param req body model.PayOrderRefundReq true "退款请求"
// @Success 200 {object} model.HttpSuccess{data=model.PayOrderRefundResp}
// @Failure 400 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/refund/{id} [post]
func (h *PayOrderHandler) RefundPayOrder(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(fiber.StatusUnauthorized, "请先登录"))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	var req model.PayOrderRefundReq
	if err := c.BodyParser(&req); err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, err.Error()))
	}

	resp, err := h.payOrderService.RefundOrder(c.Context(), id, &req)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("退款成功", resp))
}

// @Summary 模拟支付页
// @Description 模拟支付页面（仅测试环境使用），展示订单信息并提供模拟支付成功/失败按钮
// @Tags 公开接口/支付
// @Produce html
// @Param id path int true "支付订单ID"
// @Success 200 {string} string "模拟支付页面 HTML"
// @Router /api/v1/pay-order/mock-pay/{id} [get]
func (h *PayOrderHandler) MockPayPage(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}

	order, err := h.payOrderService.GetMockPayOrder(c.Context(), id)
	if err != nil {
		return mockPayError(c, err)
	}

	returnURL := "/orders.html"
	if order.ReturnURL != nil && *order.ReturnURL != "" {
		returnURL = *order.ReturnURL
	}

	subject := order.Subject
	if subject == "" {
		subject = "支付订单"
	}
	subjectHTML := html.EscapeString(subject)
	orderNoHTML := html.EscapeString(orderIDHtml(order))

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>模拟支付 - ` + subjectHTML + `</title>
<style>
  body { font-family: -apple-system, 'PingFang SC', 'Microsoft YaHei', sans-serif; background: #f5f6f8; margin: 0; }
  .box { max-width: 420px; margin: 80px auto; background: #fff; border-radius: 10px; padding: 32px; box-shadow: 0 2px 12px rgba(0,0,0,.06); }
  h2 { margin: 0 0 20px; font-size: 18px; }
  .row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px dashed #eee; font-size: 14px; }
  .row .k { color: #888; }
  .price { color: #f5222d; font-size: 24px; font-weight: 700; text-align: center; margin: 16px 0; }
  .btns { display: flex; gap: 12px; margin-top: 24px; }
  button { flex: 1; height: 40px; border: none; border-radius: 6px; font-size: 15px; cursor: pointer; }
  .ok { background: #52c41a; color: #fff; }
  .fail { background: #ff4d4f; color: #fff; }
  .tip { text-align: center; color: #999; font-size: 12px; margin-top: 16px; }
  .state { text-align: center; margin-top: 12px; font-size: 14px; color: #2f54eb; }
</style>
</head>
<body>
<div class="box">
  <h2>模拟支付（测试环境）</h2>
  <div class="row"><span class="k">订单号</span><span>` + orderNoHTML + `</span></div>
  <div class="row"><span class="k">标题</span><span>` + subjectHTML + `</span></div>
  <div class="row"><span class="k">当前状态</span><span>` + stateText(order.State) + `</span></div>
  <div class="price">¥` + fenToYuan(order.Price) + `</div>
  <div class="btns">
    <button class="ok" onclick="mockPay('success')">模拟支付成功</button>
    <button class="fail" onclick="mockPay('fail')">模拟支付失败</button>
  </div>
  <div class="state" id="msg"></div>
  <div class="tip">点击后订单按真实回调链路处理（履约/入账/积分），随后跳转回商城</div>
</div>
<script>
  async function mockPay(action) {
    var msg = document.getElementById('msg')
    var btn = document.querySelectorAll('button')
    btn.forEach(function (b) { b.disabled = true })
    msg.textContent = '处理中…'
    try {
      var resp = await fetch('/api/v1/pay-order/mock-pay/` + strconv.Itoa(order.ID) + `/' + action, { method: 'POST' })
      var body = await resp.json()
      if (body.code !== 200) {
        msg.textContent = '处理失败: ' + (body.msg || '未知错误')
        btn.forEach(function (b) { b.disabled = false })
        return
      }
      location.href = '` + returnURL + `'
    } catch (e) {
      msg.textContent = '请求失败: ' + e.message
      btn.forEach(function (b) { b.disabled = false })
    }
  }
</script>
</body>
</html>`
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}

// @Summary 模拟支付成功
// @Description 模拟支付成功回调（仅测试环境使用），走真实支付成功履约链路
// @Tags 公开接口/支付
// @Param id path int true "支付订单ID"
// @Success 200 {object} model.HttpSuccess
// @Router /api/v1/pay-order/mock-pay/{id}/success [post]
func (h *PayOrderHandler) MockPaySuccess(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	if err := h.payOrderService.MockPaySuccess(c.Context(), id); err != nil {
		return mockPayError(c, err)
	}
	return c.JSON(model.NewSuccess("模拟支付成功", nil))
}

// @Summary 模拟支付失败
// @Description 模拟支付失败回调（仅测试环境使用）
// @Tags 公开接口/支付
// @Param id path int true "支付订单ID"
// @Success 200 {object} model.HttpSuccess
// @Router /api/v1/pay-order/mock-pay/{id}/fail [post]
func (h *PayOrderHandler) MockPayFail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusBadRequest, "Invalid ID format"))
	}
	if err := h.payOrderService.MockPayFail(c.Context(), id); err != nil {
		return mockPayError(c, err)
	}
	return c.JSON(model.NewSuccess("模拟支付失败", nil))
}

// @Summary 可用支付方式
// @Description 查询当前用户可用的支付方式（支付宝/微信/余额）
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=model.PayMethodResp}
// @Failure 401 {object} model.HttpError
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/pay-methods [get]
func (h *PayOrderHandler) GetPayMethods(c *fiber.Ctx) error {
	loginUser := middleware.GetCurrentUser(c)
	if loginUser == nil {
		return c.JSON(model.NewError(fiber.StatusUnauthorized, "请先登录"))
	}

	resp, err := h.payOrderService.GetPayMethods(c.Context(), loginUser.ID)
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", resp))
}

// @Summary 获取今日统计
// @Description 获取今日支付订单统计信息
// @Tags 后台管理接口/支付订单
// @Accept json
// @Produce json
// @Success 200 {object} model.HttpSuccess{data=model.PayOrderTodayStats}
// @Failure 500 {object} model.HttpError
// @Router /api/v1/pay-order/today-stats [get]
func (h *PayOrderHandler) GetTodayStats(c *fiber.Ctx) error {
	stats, err := h.payOrderService.GetTodayStats(c.Context())
	if err != nil {
		return c.JSON(model.NewError(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(model.NewSuccess("success", stats))
}

// mockPayError 模拟支付接口的错误响应：未开启或订单不符时返回 403/404 业务码，不暴露内部错误细节。
func mockPayError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, model.ErrMockPayDisabled):
		return c.JSON(model.NewError(fiber.StatusForbidden, err.Error()))
	case errors.Is(err, model.ErrMockPayOrderInvalid):
		return c.JSON(model.NewError(fiber.StatusNotFound, err.Error()))
	default:
		slog.Error("模拟支付处理失败", "error", err)
		return c.JSON(model.NewError(fiber.StatusInternalServerError, "模拟支付处理失败"))
	}
}

// orderIDHtml mock 支付页展示用：优先商户订单号，缺省用订单ID。
func orderIDHtml(o *ent.PayOrder) string {
	if o.OutTradeNo != nil && *o.OutTradeNo != "" {
		return *o.OutTradeNo
	}
	return strconv.Itoa(o.ID)
}

// fenToYuan mock 支付页展示用：分转元字符串。
func fenToYuan(fen int) string {
	return strconv.FormatFloat(float64(fen)/100, 'f', 2, 64)
}

// stateText mock 支付页展示用：订单状态文案。
func stateText(state string) string {
	switch state {
	case "0":
		return "已取消"
	case "1":
		return "待支付"
	case "2":
		return "已支付"
	case "3":
		return "支付失败"
	case "4":
		return "已退款"
	}
	return state
}
