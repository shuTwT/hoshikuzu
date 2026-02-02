package couponusage

import (
	"context"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/couponusage"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type CouponUsageService interface {
	CreateCouponUsage(ctx context.Context, req *model.CouponUsageCreateReq) (*ent.CouponUsage, error)
	UpdateCouponUsage(ctx context.Context, id int, req *model.CouponUsageUpdateReq) (*ent.CouponUsage, error)
	DeleteCouponUsage(ctx context.Context, id int) error
	GetCouponUsage(ctx context.Context, id int) (*ent.CouponUsage, error)
	ListCouponUsages(ctx context.Context, page, pageSize int) ([]*ent.CouponUsage, int, error)
	ListAllCouponUsages(ctx context.Context) ([]*ent.CouponUsage, error)
	BatchUpdateCouponUsages(ctx context.Context, ids []int, req *model.CouponUsageBatchUpdateReq) error
	BatchDeleteCouponUsages(ctx context.Context, ids []int) error
	SearchCouponUsages(ctx context.Context, req model.CouponUsageSearchReq) ([]*ent.CouponUsage, int, error)
}

type CouponUsageServiceImpl struct {
	client *ent.Client
}

func NewCouponUsageServiceImpl(client *ent.Client) *CouponUsageServiceImpl {
	return &CouponUsageServiceImpl{client: client}
}

func (s *CouponUsageServiceImpl) CreateCouponUsage(ctx context.Context, req *model.CouponUsageCreateReq) (*ent.CouponUsage, error) {
	expireAt, err := time.Parse(time.RFC3339, req.ExpireAt)
	if err != nil {
		return nil, err
	}

	createBuilder := s.client.CouponUsage.Create().
		SetCouponCode(req.CouponCode).
		SetUserID(req.UserID).
		SetDiscountAmount(req.DiscountAmount).
		SetExpireAt(expireAt)

	if req.OrderID != nil {
		createBuilder.SetOrderID(*req.OrderID)
	}
	if req.Remark != nil {
		createBuilder.SetRemark(*req.Remark)
	}

	return createBuilder.Save(ctx)
}

func (s *CouponUsageServiceImpl) UpdateCouponUsage(ctx context.Context, id int, req *model.CouponUsageUpdateReq) (*ent.CouponUsage, error) {
	updateBuilder := s.client.CouponUsage.UpdateOneID(id)

	if req.OrderID != nil {
		updateBuilder.SetOrderID(*req.OrderID)
	}
	if req.Status != nil {
		updateBuilder.SetStatus(*req.Status)
	}
	if req.UsedAt != nil {
		usedAt, err := time.Parse(time.RFC3339, *req.UsedAt)
		if err == nil {
			updateBuilder.SetUsedAt(usedAt)
		}
	}
	if req.DiscountAmount != nil {
		updateBuilder.SetDiscountAmount(*req.DiscountAmount)
	}
	if req.ExpireAt != nil {
		expireAt, err := time.Parse(time.RFC3339, *req.ExpireAt)
		if err == nil {
			updateBuilder.SetExpireAt(expireAt)
		}
	}
	if req.Remark != nil {
		updateBuilder.SetRemark(*req.Remark)
	}

	return updateBuilder.Save(ctx)
}

func (s *CouponUsageServiceImpl) DeleteCouponUsage(ctx context.Context, id int) error {
	return s.client.CouponUsage.DeleteOneID(id).Exec(ctx)
}

func (s *CouponUsageServiceImpl) GetCouponUsage(ctx context.Context, id int) (*ent.CouponUsage, error) {
	return s.client.CouponUsage.Query().
		Where(couponusage.ID(id)).
		First(ctx)
}

func (s *CouponUsageServiceImpl) ListCouponUsages(ctx context.Context, page, pageSize int) ([]*ent.CouponUsage, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total, err := s.client.CouponUsage.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	usages, err := s.client.CouponUsage.Query().
		Order(ent.Desc(couponusage.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)

	return usages, total, err
}

func (s *CouponUsageServiceImpl) ListAllCouponUsages(ctx context.Context) ([]*ent.CouponUsage, error) {
	return s.client.CouponUsage.Query().
		Order(ent.Desc(couponusage.FieldCreatedAt)).
		All(ctx)
}

func (s *CouponUsageServiceImpl) BatchUpdateCouponUsages(ctx context.Context, ids []int, req *model.CouponUsageBatchUpdateReq) error {
	if len(ids) == 0 {
		return nil
	}

	updateBuilder := s.client.CouponUsage.Update().
		Where(couponusage.IDIn(ids...))

	if req.Status != nil {
		updateBuilder.SetStatus(*req.Status)
	}

	_, err := updateBuilder.Save(ctx)
	return err
}

func (s *CouponUsageServiceImpl) BatchDeleteCouponUsages(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.client.CouponUsage.Delete().
		Where(couponusage.IDIn(ids...)).
		Exec(ctx)
	return err
}

func (s *CouponUsageServiceImpl) SearchCouponUsages(ctx context.Context, req model.CouponUsageSearchReq) ([]*ent.CouponUsage, int, error) {
	query := s.client.CouponUsage.Query()

	if req.CouponCode != nil {
		query.Where(couponusage.CouponCode(*req.CouponCode))
	}
	if req.UserID != nil {
		query.Where(couponusage.UserID(*req.UserID))
	}
	if req.Status != nil {
		query.Where(couponusage.Status(*req.Status))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	usages, err := query.
		Order(ent.Desc(couponusage.FieldCreatedAt)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		All(ctx)

	return usages, total, err
}
