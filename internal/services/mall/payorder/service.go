package payorder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/member"
	"github.com/shuTwT/hoshikuzu/ent/payorder"
	"github.com/shuTwT/hoshikuzu/ent/postpurchase"
	"github.com/shuTwT/hoshikuzu/ent/wallet"
	"github.com/shuTwT/hoshikuzu/internal/infra/pay/epay"
	setting_service "github.com/shuTwT/hoshikuzu/internal/services/system/setting"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type PayOrderService interface {
	ListPayOrderPage(ctx context.Context, req *model.PageQuery) ([]*ent.PayOrder, int, error)
	SubmitPayOrder(ctx context.Context, userID int, req *model.PayOrderSubmitReq) (*model.PayOrderSubmitResp, error)
	HandleNotify(ctx context.Context, params map[string]string) error
	SyncOrderStatus(ctx context.Context, orderID int) (*model.PayOrderStatusResp, error)
	GetTodayStats(ctx context.Context) (*model.PayOrderTodayStats, error)
	CloseTimeoutOrders(ctx context.Context) error
	RefundOrder(ctx context.Context, orderID int, req *model.PayOrderRefundReq) (*model.PayOrderRefundResp, error)
	RechargePayOrder(ctx context.Context, userID int, req *model.PayOrderRechargeReq) (*model.PayOrderSubmitResp, error)
	GetMockPayOrder(ctx context.Context, orderID int) (*ent.PayOrder, error)
	MockPaySuccess(ctx context.Context, orderID int) error
	MockPayFail(ctx context.Context, orderID int) error
	GetPayMethods(ctx context.Context, userID int) (*model.PayMethodResp, error)
}

type PayOrderServiceImpl struct {
	db             *ent.Client
	settingService setting_service.SettingService
}

func NewPayOrderServiceImpl(db *ent.Client, settingService setting_service.SettingService) PayOrderService {
	return &PayOrderServiceImpl{db: db, settingService: settingService}
}

// paymentSettings 对应系统设置中 key='payment' 的 JSON 配置，仅声明后端消费的字段。
type paymentSettings struct {
	EnableEpay      bool   `json:"enableEpay"`
	EpayApiUrl      string `json:"epayApiUrl"`
	EpayMerchantId  string `json:"epayMerchantId"`
	EpayMerchantKey string `json:"epayMerchantKey"`
	EpayNotifyUrl   string `json:"epayNotifyUrl"`
	EpayReturnUrl   string `json:"epayReturnUrl"`
	// 订单超时分钟数，超时未支付自动关单
	OrderTimeout int `json:"orderTimeout"`
	// 充值积分比例：每支付 1 分钱发放的积分数量
	RechargePointsRate int `json:"rechargePointsRate"`
	// 模拟支付开关：开启后不调用真实支付网关，返回本地模拟支付页（仅测试环境使用）
	EnableMockPay bool `json:"enableMockPay"`
	// 直连渠道开关（用于支付方式可用性判断，直连支付暂未实现）
	EnableAlipay    bool `json:"enableAlipay"`
	EnableWechatPay bool `json:"enableWechatPay"`
}

// loadPaymentSettings 读取系统设置 key='payment' 的 JSON 配置。
func (s *PayOrderServiceImpl) loadPaymentSettings(ctx context.Context) (*paymentSettings, error) {
	setting, err := s.settingService.GetSettingByKey(ctx, "payment")
	if err != nil {
		return nil, err
	}

	var ps paymentSettings
	if err := json.Unmarshal([]byte(setting.Value), &ps); err != nil {
		return nil, fmt.Errorf("解析支付配置失败: %w", err)
	}
	return &ps, nil
}

// loadEpayConfig 从系统设置读取易支付配置，每次下单按需读取以保证设置修改即时生效。
func (s *PayOrderServiceImpl) loadEpayConfig(ctx context.Context) (epay.Config, bool, error) {
	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return epay.Config{}, false, fmt.Errorf("支付配置未初始化，请先在系统设置-支付设置中配置")
		}
		return epay.Config{}, false, err
	}

	cfg := epay.Config{
		MchID:     ps.EpayMerchantId,
		Key:       ps.EpayMerchantKey,
		APIURL:    ps.EpayApiUrl,
		NotifyURL: ps.EpayNotifyUrl,
		ReturnURL: ps.EpayReturnUrl,
	}
	return cfg, ps.EnableEpay, nil
}

// 查询支付订单列表
func (s *PayOrderServiceImpl) ListPayOrderPage(ctx context.Context, req *model.PageQuery) ([]*ent.PayOrder, int, error) {
	orders, err := s.db.PayOrder.Query().
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.db.PayOrder.Query().
		Order(ent.Desc(payorder.FieldID)).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return orders, count, nil
}

// SubmitPayOrder 提交文章付费/商品购买订单并跳转易支付下单。
func (s *PayOrderServiceImpl) SubmitPayOrder(ctx context.Context, userID int, req *model.PayOrderSubmitReq) (*model.PayOrderSubmitResp, error) {
	// 余额支付：直接扣减钱包余额完成支付，不经过第三方网关
	if req.ChannelType == model.PayChannelBalance {
		return s.submitBalancePay(ctx, userID, req)
	}

	order, err := s.createPayOrder(ctx, userID, req.OrderType, req.ChannelType, req.Name, req.Name, req.Money, req.PostId, req.ProductId)
	if err != nil {
		return nil, err
	}

	payURL, tradeNO, err := s.submitToEpay(ctx, order, req.ChannelType, req.ReturnUrl, req.Name)
	if err != nil {
		return nil, err
	}

	return &model.PayOrderSubmitResp{
		OrderID:    order.ID,
		OutTradeNo: *order.OutTradeNo,
		PayURL:     payURL,
		TradeNO:    tradeNO,
	}, nil
}

// submitBalancePay 余额支付：扣减钱包余额并将订单置为已支付（含履约），全部在一个事务内完成。
func (s *PayOrderServiceImpl) submitBalancePay(ctx context.Context, userID int, req *model.PayOrderSubmitReq) (*model.PayOrderSubmitResp, error) {
	order, err := s.createPayOrder(ctx, userID, req.OrderType, model.PayChannelBalance, req.Name, req.Name, req.Money, req.PostId, req.ProductId)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	// 校验并扣减钱包余额
	walletEnt, err := tx.Wallet.Query().Where(wallet.UserIDEQ(userID)).Only(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	if walletEnt.Balance < req.Money {
		_ = tx.Rollback()
		return nil, fmt.Errorf("余额不足，当前余额 ¥%s", strconv.FormatFloat(float64(walletEnt.Balance)/100, 'f', 2, 64))
	}
	if _, err := tx.Wallet.UpdateOneID(walletEnt.ID).
		SetBalance(walletEnt.Balance - req.Money).
		SetTotalExpense(walletEnt.TotalExpense + req.Money).
		Save(ctx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 订单置为已支付并落库交易号
	order, err = tx.PayOrder.UpdateOneID(order.ID).
		SetState("2").
		SetOrderID(fmt.Sprintf("BAL%d", order.ID)).
		SetPayURL("").
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 履约（文章付费写购买记录）
	if err := writePostPurchase(ctx, tx, order); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	slog.Info("余额支付成功", "user_id", userID, "order_id", order.ID, "amount", req.Money)
	return &model.PayOrderSubmitResp{
		OrderID:    order.ID,
		OutTradeNo: *order.OutTradeNo,
		PayURL:     "",
		TradeNO:    fmt.Sprintf("BAL%d", order.ID),
	}, nil
}

// RechargePayOrder 提交余额充值订单并跳转易支付下单，支付成功后回调入账钱包并发放会员积分。
func (s *PayOrderServiceImpl) RechargePayOrder(ctx context.Context, userID int, req *model.PayOrderRechargeReq) (*model.PayOrderSubmitResp, error) {
	// 充值不能使用余额支付
	if req.ChannelType == model.PayChannelBalance {
		return nil, fmt.Errorf("充值不支持余额支付")
	}

	subject := "余额充值"
	order, err := s.createPayOrder(ctx, userID, model.PayOrderTypeRecharge, req.ChannelType, subject, subject, req.Amount, 0, 0)
	if err != nil {
		return nil, err
	}

	payURL, tradeNO, err := s.submitToEpay(ctx, order, req.ChannelType, req.ReturnUrl, subject)
	if err != nil {
		return nil, err
	}

	return &model.PayOrderSubmitResp{
		OrderID:    order.ID,
		OutTradeNo: *order.OutTradeNo,
		PayURL:     payURL,
		TradeNO:    tradeNO,
	}, nil
}

// createPayOrder 创建支付订单并生成商户订单号，金额单位为分。
func (s *PayOrderServiceImpl) createPayOrder(ctx context.Context, userID int, orderType, channelType, subject, body string, amount, postID, productID int) (*ent.PayOrder, error) {
	builder := s.db.PayOrder.Create().
		SetUserID(userID).
		SetOrderType(orderType).
		SetChannelType(channelType).
		SetSubject(subject).
		SetBody(body).
		SetOrderPrice(amount).
		SetPrice(amount).
		SetState("1")
	if postID > 0 {
		builder = builder.SetPostID(postID)
	}
	if productID > 0 {
		builder = builder.SetProductID(productID)
	}
	order, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	// 生成商户订单号
	orderNo := time.Now().Format("20060102150405") + fmt.Sprintf("%09d", order.ID)
	return s.db.PayOrder.UpdateOneID(order.ID).SetOutTradeNo(orderNo).Save(ctx)
}

// submitToEpay 下单并落库支付链接与易支付订单号。开启模拟支付时跳过真实网关，返回本地模拟支付页，金额单位为分。
func (s *PayOrderServiceImpl) submitToEpay(ctx context.Context, order *ent.PayOrder, channelType, returnURL, name string) (string, string, error) {
	moneyYuan := strconv.FormatFloat(float64(order.Price)/100, 'f', 2, 64)

	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return "", "", fmt.Errorf("支付配置未初始化，请先在系统设置-支付设置中配置")
		}
		return "", "", err
	}

	// 模拟支付：不调用真实网关，返回本地 mock 支付页
	if ps.EnableMockPay {
		mockURL := fmt.Sprintf("/api/v1/pay-order/mock-pay/%d", order.ID)
		if _, err := s.db.PayOrder.UpdateOneID(order.ID).
			SetPayURL(mockURL).
			SetOrderID(fmt.Sprintf("MOCK%d", order.ID)).
			SetReturnURL(returnURL).
			Save(ctx); err != nil {
			return "", "", err
		}
		return mockURL, fmt.Sprintf("MOCK%d", order.ID), nil
	}

	if !ps.EnableEpay {
		return "", "", fmt.Errorf("易支付未启用")
	}

	cfg := epay.Config{
		MchID:     ps.EpayMerchantId,
		Key:       ps.EpayMerchantKey,
		APIURL:    ps.EpayApiUrl,
		NotifyURL: ps.EpayNotifyUrl,
		ReturnURL: ps.EpayReturnUrl,
	}
	params := epay.V1PayRequestParams{
		PID:        cfg.MchID,
		Type:       channelType,
		OutTradeNo: *order.OutTradeNo,
		Name:       name,
		Money:      moneyYuan,
		ReturnURL:  &returnURL,
	}
	resp, err := epay.NewV1Client(cfg).CreateOrder(params)
	if err != nil {
		// 下单失败，记录错误信息
		_ = s.db.PayOrder.UpdateOneID(order.ID).SetState("3").SetErrorMsg(err.Error()).Exec(ctx)
		return "", "", err
	}

	// 存储支付链接与易支付订单号
	payURL := ""
	if resp.Payurl != nil {
		payURL = *resp.Payurl
	}
	if _, err := s.db.PayOrder.UpdateOneID(order.ID).
		SetPayURL(payURL).
		SetOrderID(resp.TradeNO).
		SetReturnURL(returnURL).
		Save(ctx); err != nil {
		return "", "", err
	}

	return payURL, resp.TradeNO, nil
}

// GetMockPayOrder 查询模拟支付订单：要求支付设置开启了模拟支付，且订单是通过模拟支付下单的。
// 模拟支付接口是公开的，所有入口都必须先经过这里的校验。
func (s *PayOrderServiceImpl) GetMockPayOrder(ctx context.Context, orderID int) (*ent.PayOrder, error) {
	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrMockPayDisabled
		}
		return nil, err
	}
	if !ps.EnableMockPay {
		return nil, model.ErrMockPayDisabled
	}

	order, err := s.db.PayOrder.Get(ctx, orderID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, model.ErrMockPayOrderInvalid
		}
		return nil, err
	}
	if order.OrderID == nil || *order.OrderID != fmt.Sprintf("MOCK%d", order.ID) {
		return nil, model.ErrMockPayOrderInvalid
	}
	return order, nil
}

// MockPaySuccess 模拟支付成功：复用真实支付成功的履约逻辑（幂等，含文章购买记录/充值入账积分）。
func (s *PayOrderServiceImpl) MockPaySuccess(ctx context.Context, orderID int) error {
	order, err := s.GetMockPayOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order.State == "2" {
		return nil
	}
	if order.State != "1" {
		return model.ErrMockPayOrderInvalid
	}
	return s.handlePaySuccess(ctx, order, fmt.Sprintf("MOCK%d", order.ID))
}

// MockPayFail 模拟支付失败：待支付订单置为失败状态。
func (s *PayOrderServiceImpl) MockPayFail(ctx context.Context, orderID int) error {
	order, err := s.GetMockPayOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order.State != "1" {
		return model.ErrMockPayOrderInvalid
	}
	_, err = s.db.PayOrder.UpdateOneID(orderID).
		Where(payorder.State("1")).
		SetState("3").
		SetErrorMsg("模拟支付失败").
		Save(ctx)
	return err
}

// GetPayMethods 查询当前用户可用的支付方式：支付宝/微信取决于易支付或直连渠道启用，余额取决于钱包余额。
func (s *PayOrderServiceImpl) GetPayMethods(ctx context.Context, userID int) (*model.PayMethodResp, error) {
	resp := &model.PayMethodResp{}

	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// 未配置支付设置：仅余额可用
			ps = &paymentSettings{}
		} else {
			return nil, err
		}
	}

	// 模拟支付模式下所有渠道均可测试
	epayAvailable := ps.EnableEpay || ps.EnableMockPay
	resp.Alipay = epayAvailable || ps.EnableAlipay
	resp.Wechat = epayAvailable || ps.EnableWechatPay

	// 余额支付：钱包存在即可选（余额不足时由下单接口提示，前端始终展示该方式）
	walletEnt, err := s.db.Wallet.Query().Where(wallet.UserIDEQ(userID)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
	} else {
		resp.Balance = true
		resp.BalanceAmount = walletEnt.Balance
	}

	return resp, nil
}

// handlePaySuccess 标记订单为已支付并执行业务履约（文章付费写购买记录、充值入账钱包与积分），幂等。
func (s *PayOrderServiceImpl) handlePaySuccess(ctx context.Context, order *ent.PayOrder, tradeNo string) error {
	if order.State == "2" {
		return nil
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return err
	}

	updater := tx.PayOrder.UpdateOneID(order.ID).SetState("2")
	if tradeNo != "" {
		// 落库易支付订单号，退款时需要使用
		updater.SetOrderID(tradeNo)
	}
	order, err = updater.Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// 文章付费：写购买记录（唯一索引保证不重复）
	if err := writePostPurchase(ctx, tx, order); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 余额充值：钱包入账并发放会员积分
	if order.OrderType == model.PayOrderTypeRecharge {
		if err := s.creditWalletAndPoints(ctx, tx, order); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// writePostPurchase 文章付费履约：写入购买记录（唯一索引保证不重复），在事务内调用。
func writePostPurchase(ctx context.Context, tx *ent.Tx, order *ent.PayOrder) error {
	if order.OrderType != model.PayOrderTypePost || order.PostID == 0 {
		return nil
	}

	exist, err := tx.PostPurchase.Query().
		Where(postpurchase.UserIDEQ(order.UserID), postpurchase.PostIDEQ(order.PostID)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exist {
		if _, err := tx.PostPurchase.Create().
			SetUserID(order.UserID).
			SetPostID(order.PostID).
			SetOrderID(order.ID).
			Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

// creditWalletAndPoints 充值履约：余额入账并发放会员积分，在支付成功事务内调用。
func (s *PayOrderServiceImpl) creditWalletAndPoints(ctx context.Context, tx *ent.Tx, order *ent.PayOrder) error {
	walletEnt, err := tx.Wallet.Query().Where(wallet.UserIDEQ(order.UserID)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return err
		}
	} else {
		if _, err := tx.Wallet.UpdateOneID(walletEnt.ID).
			SetBalance(walletEnt.Balance + order.Price).
			SetTotalIncome(walletEnt.TotalIncome + order.Price).
			Save(ctx); err != nil {
			return err
		}
	}

	memberEnt, err := tx.Member.Query().Where(member.UserIDEQ(order.UserID)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return err
		}
	} else {
		points := s.calcRechargePoints(ctx, order.Price)
		if points > 0 {
			if _, err := tx.Member.UpdateOneID(memberEnt.ID).
				SetPoints(memberEnt.Points + points).
				Save(ctx); err != nil {
				return err
			}
			// 落库实际发放的积分，退款时按此精确扣回
			if _, err := tx.PayOrder.UpdateOneID(order.ID).
				SetPointsGranted(points).
				Save(ctx); err != nil {
				return err
			}
			slog.Info("充值发放积分", "user_id", order.UserID, "amount", order.Price, "points", points)
		}
	}

	return nil
}

// revokeWalletAndPoints 充值退款履约：扣回钱包余额与会员积分，在退款事务内调用。
func (s *PayOrderServiceImpl) revokeWalletAndPoints(ctx context.Context, tx *ent.Tx, order *ent.PayOrder, refundAmount int) error {
	walletEnt, err := tx.Wallet.Query().Where(wallet.UserIDEQ(order.UserID)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return err
		}
	} else {
		if _, err := tx.Wallet.UpdateOneID(walletEnt.ID).
			SetBalance(walletEnt.Balance - refundAmount).
			SetTotalExpense(walletEnt.TotalExpense + refundAmount).
			Save(ctx); err != nil {
			return err
		}
	}

	memberEnt, err := tx.Member.Query().Where(member.UserIDEQ(order.UserID)).Only(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			return err
		}
	} else {
		// 优先按订单记录的实际发放积分扣回，兼容旧订单（无记录时按当前比例计算）
		points := order.PointsGranted
		if points <= 0 {
			points = s.calcRechargePoints(ctx, refundAmount)
		}
		if points > 0 {
			if _, err := tx.Member.UpdateOneID(memberEnt.ID).
				SetPoints(memberEnt.Points - points).
				Save(ctx); err != nil {
				return err
			}
			slog.Info("充值退款扣回积分", "user_id", order.UserID, "order_id", order.ID, "points", points)
		}
	}

	return nil
}

// calcRechargePoints 按充值金额（分）与系统设置的兑换比例计算应发积分。
func (s *PayOrderServiceImpl) calcRechargePoints(ctx context.Context, amountFen int) int {
	return amountFen * s.getRechargePointsRate(ctx)
}

// getRechargePointsRate 读取充值积分比例（每 1 分钱兑换的积分数量），未配置或非法时默认 1。
func (s *PayOrderServiceImpl) getRechargePointsRate(ctx context.Context) int {
	const defaultRate = 1

	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		return defaultRate
	}
	if ps.RechargePointsRate <= 0 {
		return defaultRate
	}
	return ps.RechargePointsRate
}

// HandleNotify 处理易支付异步回调，支付通知与退款通知共用回调地址，按参数区分。
func (s *PayOrderServiceImpl) HandleNotify(ctx context.Context, params map[string]string) error {
	// 退款通知带 refund_no 参数，支付通知带 trade_status
	if params["refund_no"] != "" {
		return s.handleRefundNotify(ctx, params)
	}

	cfg, _, err := s.loadEpayConfig(ctx)
	if err != nil {
		return err
	}

	sign, ok := params["sign"]
	if !ok {
		return fmt.Errorf("缺少签名参数")
	}
	if !epay.VerifySign(params, cfg.Key, sign) {
		return fmt.Errorf("签名校验失败")
	}

	// 仅处理支付成功的通知
	if params["trade_status"] != "TRADE_SUCCESS" {
		return nil
	}

	outTradeNo := params["out_trade_no"]
	order, err := s.db.PayOrder.Query().
		Where(payorder.OutTradeNoEQ(outTradeNo)).
		Only(ctx)
	if err != nil {
		return err
	}

	return s.handlePaySuccess(ctx, order, params["trade_no"])
}

// SyncOrderStatus 主动查询易支付订单状态并同步本地（用于补单/后台手动刷新）。
func (s *PayOrderServiceImpl) SyncOrderStatus(ctx context.Context, orderID int) (*model.PayOrderStatusResp, error) {
	order, err := s.db.PayOrder.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 已支付的订单无需同步
	if order.State == "2" {
		return toPayOrderStatusResp(order), nil
	}

	// 模拟支付模式下不查外部网关，直接返回当前状态
	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		return nil, err
	}
	if ps.EnableMockPay {
		return toPayOrderStatusResp(order), nil
	}
	if !ps.EnableEpay {
		return nil, fmt.Errorf("易支付未启用")
	}

	cfg := epay.Config{
		MchID:     ps.EpayMerchantId,
		Key:       ps.EpayMerchantKey,
		APIURL:    ps.EpayApiUrl,
		NotifyURL: ps.EpayNotifyUrl,
		ReturnURL: ps.EpayReturnUrl,
	}

	pid, err := strconv.Atoi(cfg.MchID)
	if err != nil {
		return nil, fmt.Errorf("商户ID无效: %w", err)
	}
	tradeNo := ""
	if order.OrderID != nil {
		tradeNo = *order.OrderID
	}
	outTradeNo := ""
	if order.OutTradeNo != nil {
		outTradeNo = *order.OutTradeNo
	}

	queryResp, err := epay.NewV1Client(cfg).QueryOrder(pid, cfg.Key, tradeNo, outTradeNo)
	if err != nil {
		return nil, err
	}

	// 易支付 status=2 表示已支付
	if queryResp.Status == 2 {
		if err := s.handlePaySuccess(ctx, order, queryResp.TradeNO); err != nil {
			return nil, err
		}
		order, err = s.db.PayOrder.Get(ctx, orderID)
		if err != nil {
			return nil, err
		}
	}

	return toPayOrderStatusResp(order), nil
}

func toPayOrderStatusResp(o *ent.PayOrder) *model.PayOrderStatusResp {
	outTradeNo := ""
	if o.OutTradeNo != nil {
		outTradeNo = *o.OutTradeNo
	}
	return &model.PayOrderStatusResp{
		ID:         o.ID,
		OutTradeNo: outTradeNo,
		State:      o.State,
		Subject:    o.Subject,
		Price:      o.Price,
	}
}

func (s *PayOrderServiceImpl) GetTodayStats(ctx context.Context) (*model.PayOrderTodayStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	orders, err := s.db.PayOrder.Query().
		Where(payorder.CreatedAtGTE(todayStart)).
		Where(payorder.CreatedAtLT(todayEnd)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	total := len(orders)
	var amount float64
	pending := 0
	successCount := 0

	for _, order := range orders {
		if order.State == "2" {
			successCount++
			amount += float64(order.Price) / 100
		} else if order.State == "1" {
			pending++
		}
	}

	var successRate float64
	if total > 0 {
		successRate = float64(successCount) / float64(total) * 100
	}

	return &model.PayOrderTodayStats{
		Total:       total,
		Amount:      amount,
		SuccessRate: successRate,
		Pending:     pending,
	}, nil
}

// getOrderTimeoutMinutes 读取订单超时分钟数，未配置或配置非法时使用默认值 30 分钟。
func (s *PayOrderServiceImpl) getOrderTimeoutMinutes(ctx context.Context) int {
	const defaultTimeout = 30

	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		return defaultTimeout
	}
	if ps.OrderTimeout <= 0 {
		return defaultTimeout
	}
	return ps.OrderTimeout
}

// CloseTimeoutOrders 关闭超时未支付订单（state 1 -> 0），由定时任务周期调用。
func (s *PayOrderServiceImpl) CloseTimeoutOrders(ctx context.Context) error {
	cutoff := time.Now().Add(-time.Duration(s.getOrderTimeoutMinutes(ctx)) * time.Minute)

	updated, err := s.db.PayOrder.Update().
		Where(
			payorder.StateEQ("1"),
			payorder.CreatedAtLT(cutoff),
		).
		SetState("0").
		SetErrorMsg("超时未支付，系统自动关单").
		Save(ctx)
	if err != nil {
		return err
	}

	if updated > 0 {
		slog.Info("超时关单完成", "count", updated)
	}
	return nil
}

// RefundOrder 发起退款（当前仅支持全额退款），state 2 -> 4，并回滚文章购买记录。
func (s *PayOrderServiceImpl) RefundOrder(ctx context.Context, orderID int, req *model.PayOrderRefundReq) (*model.PayOrderRefundResp, error) {
	order, err := s.db.PayOrder.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 幂等：已退款直接返回现有信息
	if order.State == "4" {
		return toPayOrderRefundResp(order), nil
	}
	if order.State != "2" {
		return nil, fmt.Errorf("订单状态 %s 不可退款", order.State)
	}

	// 无退款记录表，部分退款无法追踪累计金额与权益回滚，当前仅支持全额退款
	amount := req.Amount
	if amount == 0 {
		amount = order.Price
	}
	if amount < order.Price {
		return nil, fmt.Errorf("当前仅支持全额退款")
	}

	refundNo := ""
	if order.OutTradeNo != nil && *order.OutTradeNo != "" {
		refundNo = fmt.Sprintf("R%s%s", *order.OutTradeNo, time.Now().Format("20060102150405"))
	} else {
		refundNo = fmt.Sprintf("R%d%s", order.ID, time.Now().Format("20060102150405"))
	}

	// 真实退款：调用易支付网关；余额支付订单与模拟支付模式下跳过网关调用
	ps, err := s.loadPaymentSettings(ctx)
	if err != nil {
		return nil, err
	}
	isBalanceOrder := order.ChannelType != nil && *order.ChannelType == model.PayChannelBalance
	if !isBalanceOrder && !ps.EnableMockPay {
		if !ps.EnableEpay {
			return nil, fmt.Errorf("易支付未启用")
		}

		cfg := epay.Config{
			MchID:     ps.EpayMerchantId,
			Key:       ps.EpayMerchantKey,
			APIURL:    ps.EpayApiUrl,
			NotifyURL: ps.EpayNotifyUrl,
			ReturnURL: ps.EpayReturnUrl,
		}
		// 易支付订单号本地缺失时先查单补齐（部分订单经回调支付成功但未落库 trade_no）
		tradeNo, err := s.getEpayTradeNo(ctx, cfg, order)
		if err != nil {
			return nil, err
		}

		moneyYuan := strconv.FormatFloat(float64(amount)/100, 'f', 2, 64)
		pid, err := strconv.Atoi(cfg.MchID)
		if err != nil {
			return nil, fmt.Errorf("商户ID无效: %w", err)
		}
		if _, err := epay.NewV1Client(cfg).RefundOrder(pid, cfg.Key, tradeNo, refundNo, moneyYuan); err != nil {
			return nil, fmt.Errorf("易支付退款失败: %w", err)
		}
	}

	// 更新订单并回滚权益与资金，事务保证一致性
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := tx.PayOrder.UpdateOneID(order.ID).
		SetState("4").
		SetRefundNo(refundNo).
		SetRefundAmount(amount).
		SetRefundAt(time.Now()).
		Save(ctx); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	// 余额支付订单：退款退回钱包余额
	if isBalanceOrder {
		if err := s.refundToWallet(ctx, tx, order, amount); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}
	// 文章付费权益回滚
	if err := deletePostPurchase(ctx, tx, order); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	// 余额充值：扣回钱包余额与会员积分
	if order.OrderType == model.PayOrderTypeRecharge {
		if err := s.revokeWalletAndPoints(ctx, tx, order, amount); err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	order, err = s.db.PayOrder.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	slog.Info("订单退款成功", "order_id", order.ID, "refund_no", refundNo, "amount", amount)
	return toPayOrderRefundResp(order), nil
}

// handleRefundNotify 处理易支付退款异步通知：验签 -> 找订单 -> 标记退款并回滚权益，幂等。
func (s *PayOrderServiceImpl) handleRefundNotify(ctx context.Context, params map[string]string) error {
	cfg, _, err := s.loadEpayConfig(ctx)
	if err != nil {
		return err
	}

	sign, ok := params["sign"]
	if !ok {
		return fmt.Errorf("缺少签名参数")
	}
	if !epay.VerifySign(params, cfg.Key, sign) {
		return fmt.Errorf("签名校验失败")
	}

	refundNo := params["refund_no"]
	outTradeNo := params["out_trade_no"]
	order, err := s.db.PayOrder.Query().
		Where(payorder.OutTradeNoEQ(outTradeNo)).
		Only(ctx)
	if err != nil {
		return err
	}

	// 幂等：已退款订单直接跳过
	if order.State == "4" {
		return nil
	}

	amount := 0
	if money := params["money"]; money != "" {
		if moneyF, parseErr := strconv.ParseFloat(money, 64); parseErr == nil {
			amount = int(moneyF * 100)
		}
	}

	err = func() error {
		tx, err := s.db.Tx(ctx)
		if err != nil {
			return err
		}

		if _, err := tx.PayOrder.UpdateOneID(order.ID).
			SetState("4").
			SetRefundNo(refundNo).
			SetRefundAmount(amount).
			SetRefundAt(time.Now()).
			Save(ctx); err != nil {
			_ = tx.Rollback()
			return err
		}
		// 文章付费权益回滚
		if err := deletePostPurchase(ctx, tx, order); err != nil {
			_ = tx.Rollback()
			return err
		}
		// 余额充值：扣回钱包余额与会员积分
		if order.OrderType == model.PayOrderTypeRecharge {
			if err := s.revokeWalletAndPoints(ctx, tx, order, amount); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	}()
	if err != nil {
		return err
	}

	slog.Info("退款通知处理完成", "order_id", order.ID, "refund_no", refundNo, "amount", amount)
	return nil
}

// deletePostPurchase 文章付费权益回滚：删除对应购买记录，在事务内调用。
func deletePostPurchase(ctx context.Context, tx *ent.Tx, order *ent.PayOrder) error {
	if order.OrderType != model.PayOrderTypePost || order.PostID == 0 {
		return nil
	}
	_, err := tx.PostPurchase.Delete().
		Where(postpurchase.OrderIDEQ(order.ID)).
		Exec(ctx)
	return err
}

// refundToWallet 余额支付订单退款：退款金额退回钱包余额，在事务内调用。
func (s *PayOrderServiceImpl) refundToWallet(ctx context.Context, tx *ent.Tx, order *ent.PayOrder, refundAmount int) error {
	walletEnt, err := tx.Wallet.Query().Where(wallet.UserIDEQ(order.UserID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return err
	}
	if _, err := tx.Wallet.UpdateOneID(walletEnt.ID).
		SetBalance(walletEnt.Balance + refundAmount).
		SetTotalExpense(walletEnt.TotalExpense - refundAmount).
		Save(ctx); err != nil {
		return err
	}
	return nil
}

// getEpayTradeNo 获取易支付订单号，本地未落库时通过查询接口补齐。
func (s *PayOrderServiceImpl) getEpayTradeNo(ctx context.Context, cfg epay.Config, order *ent.PayOrder) (string, error) {
	if order.OrderID != nil && *order.OrderID != "" {
		return *order.OrderID, nil
	}
	if order.OutTradeNo == nil || *order.OutTradeNo == "" {
		return "", fmt.Errorf("订单缺少交易号")
	}

	pid, err := strconv.Atoi(cfg.MchID)
	if err != nil {
		return "", fmt.Errorf("商户ID无效: %w", err)
	}
	resp, err := epay.NewV1Client(cfg).QueryOrder(pid, cfg.Key, "", *order.OutTradeNo)
	if err != nil {
		return "", err
	}
	if resp.TradeNO == "" {
		return "", fmt.Errorf("易支付未找到对应订单")
	}
	return resp.TradeNO, nil
}

func toPayOrderRefundResp(o *ent.PayOrder) *model.PayOrderRefundResp {
	refundNo := ""
	if o.RefundNo != nil {
		refundNo = *o.RefundNo
	}
	return &model.PayOrderRefundResp{
		OrderID:      o.ID,
		RefundNo:     refundNo,
		RefundAmount: o.RefundAmount,
		State:        o.State,
	}
}
