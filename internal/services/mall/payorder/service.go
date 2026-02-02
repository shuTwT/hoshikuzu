package payorder

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/payorder"
	"github.com/shuTwT/hoshikuzu/internal/infra/pay/epay"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type PayOrderService interface {
	ListPayOrderPage(ctx context.Context, req *model.PageQuery) ([]*ent.PayOrder, int, error)
	SubmitPayOrder(ctx context.Context, req *model.PayOrderSubmitReq) error
}

type PayOrderServiceImpl struct {
	db         *ent.Client
	epayClient *epay.V1Client
}

func NewPayOrderServiceImpl(db *ent.Client) PayOrderService {
	epayConfig := epay.Config{
		MchID:     config.GetString(config.PAY_EPAY_MERCHANT_ID),
		Key:       config.GetString(config.PAY_EPAY_MERCHANT_KEY),
		APIURL:    config.GetString(config.PAY_EPAY_API_URL),
		NotifyURL: config.GetString(config.PAY_EPAY_NOTIFY_URL),
		ReturnURL: config.GetString(config.PAY_EPAY_RETURN_URL),
	}
	epayClient := epay.NewV1Client(epayConfig)
	return &PayOrderServiceImpl{db: db, epayClient: epayClient}
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

func (s *PayOrderServiceImpl) CreatePayOrder(ctx context.Context) (*ent.PayOrder, error) {
	order, err := s.db.PayOrder.Create().
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// 更新支付订单号
func (s *PayOrderServiceImpl) UpdatePayOrderNO(ctx context.Context, order *ent.PayOrder) (*ent.PayOrder, error) {
	id := order.ID

	now := time.Now().Format("20060102150405")
	idStr := fmt.Sprintf("%09d", id)

	orderNo := now + idStr
	return s.db.PayOrder.UpdateOne(order).
		SetOutTradeNo(orderNo).
		Save(ctx)
}

func (s *PayOrderServiceImpl) SubmitPayOrder(ctx context.Context, req *model.PayOrderSubmitReq) error {
	// 判断是否启用了易支付
	// 先创建一个支付订单
	order, err := s.CreatePayOrder(ctx)
	if err != nil {
		return err
	}
	// 更新支付订单号
	order, err = s.UpdatePayOrderNO(ctx, order)
	if err != nil {
		return err
	}

	if config.GetBool(config.PAY_EPAY_ENABLE) {
		// 调用易支付接口
		params := epay.V1PayRequestParams{
			PID:        config.GetString(config.PAY_EPAY_MERCHANT_ID),
			Type:       req.ChannelType,
			OutTradeNo: *order.OutTradeNo,
			Name:       req.Name,
			Money:      strconv.Itoa(req.Money),
		}
		resp, err := s.epayClient.CreateOrder(params)
		if err != nil {
			return err
		}
		if resp.Code != 1 {
			return fmt.Errorf("创建支付订单失败: %s", resp.Msg)
		}
	}
	return nil
}
