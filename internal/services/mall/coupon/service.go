package coupon

import (
	"context"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/coupon"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type CouponService interface {
	CreateCoupon(ctx context.Context, req *model.CouponCreateReq) (*ent.Coupon, error)
	UpdateCoupon(ctx context.Context, id int, req *model.CouponUpdateReq) (*ent.Coupon, error)
	DeleteCoupon(ctx context.Context, id int) error
	GetCoupon(ctx context.Context, id int) (*ent.Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*ent.Coupon, error)
	ListCoupons(ctx context.Context, page, pageSize int) ([]*ent.Coupon, int, error)
	ListAllCoupons(ctx context.Context) ([]*ent.Coupon, error)
	BatchUpdateCoupons(ctx context.Context, ids []int, req *model.CouponBatchUpdateReq) error
	BatchDeleteCoupons(ctx context.Context, ids []int) error
	SearchCoupons(ctx context.Context, req model.CouponSearchReq) ([]*ent.Coupon, int, error)
}

type CouponServiceImpl struct {
	client *ent.Client
}

func NewCouponServiceImpl(client *ent.Client) *CouponServiceImpl {
	return &CouponServiceImpl{client: client}
}

func (s *CouponServiceImpl) CreateCoupon(ctx context.Context, req *model.CouponCreateReq) (*ent.Coupon, error) {
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	createBuilder := s.client.Coupon.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetType(req.Type).
		SetValue(req.Value).
		SetMinAmount(req.MinAmount).
		SetMaxDiscount(req.MaxDiscount).
		SetTotalCount(req.TotalCount).
		SetPerUserLimit(req.PerUserLimit).
		SetStartTime(startTime).
		SetEndTime(endTime)

	if req.Description != nil {
		createBuilder.SetDescription(*req.Description)
	}
	if req.Active != nil {
		createBuilder.SetActive(*req.Active)
	}
	if req.Image != nil {
		createBuilder.SetImage(*req.Image)
	}
	if req.ProductIds != nil {
		createBuilder.SetProductIds(req.ProductIds)
	}
	if req.CategoryIds != nil {
		createBuilder.SetCategoryIds(req.CategoryIds)
	}

	return createBuilder.Save(ctx)
}

func (s *CouponServiceImpl) UpdateCoupon(ctx context.Context, id int, req *model.CouponUpdateReq) (*ent.Coupon, error) {
	updateBuilder := s.client.Coupon.UpdateOneID(id)

	if req.Name != nil {
		updateBuilder.SetName(*req.Name)
	}
	if req.Description != nil {
		updateBuilder.SetDescription(*req.Description)
	}
	if req.Type != nil {
		updateBuilder.SetType(*req.Type)
	}
	if req.Value != nil {
		updateBuilder.SetValue(*req.Value)
	}
	if req.MinAmount != nil {
		updateBuilder.SetMinAmount(*req.MinAmount)
	}
	if req.MaxDiscount != nil {
		updateBuilder.SetMaxDiscount(*req.MaxDiscount)
	}
	if req.TotalCount != nil {
		updateBuilder.SetTotalCount(*req.TotalCount)
	}
	if req.PerUserLimit != nil {
		updateBuilder.SetPerUserLimit(*req.PerUserLimit)
	}
	if req.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *req.StartTime)
		if err == nil {
			updateBuilder.SetStartTime(startTime)
		}
	}
	if req.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *req.EndTime)
		if err == nil {
			updateBuilder.SetEndTime(endTime)
		}
	}
	if req.Active != nil {
		updateBuilder.SetActive(*req.Active)
	}
	if req.Image != nil {
		updateBuilder.SetImage(*req.Image)
	}
	if req.ProductIds != nil {
		updateBuilder.SetProductIds(req.ProductIds)
	}
	if req.CategoryIds != nil {
		updateBuilder.SetCategoryIds(req.CategoryIds)
	}

	return updateBuilder.Save(ctx)
}

func (s *CouponServiceImpl) DeleteCoupon(ctx context.Context, id int) error {
	return s.client.Coupon.DeleteOneID(id).Exec(ctx)
}

func (s *CouponServiceImpl) GetCoupon(ctx context.Context, id int) (*ent.Coupon, error) {
	return s.client.Coupon.Query().
		Where(coupon.ID(id)).
		First(ctx)
}

func (s *CouponServiceImpl) GetCouponByCode(ctx context.Context, code string) (*ent.Coupon, error) {
	return s.client.Coupon.Query().
		Where(coupon.Code(code)).
		First(ctx)
}

func (s *CouponServiceImpl) ListCoupons(ctx context.Context, page, pageSize int) ([]*ent.Coupon, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total, err := s.client.Coupon.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	coupons, err := s.client.Coupon.Query().
		Order(ent.Desc(coupon.FieldCreatedAt)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)

	return coupons, total, err
}

func (s *CouponServiceImpl) ListAllCoupons(ctx context.Context) ([]*ent.Coupon, error) {
	return s.client.Coupon.Query().
		Order(ent.Desc(coupon.FieldCreatedAt)).
		All(ctx)
}

func (s *CouponServiceImpl) BatchUpdateCoupons(ctx context.Context, ids []int, req *model.CouponBatchUpdateReq) error {
	if len(ids) == 0 {
		return nil
	}

	updateBuilder := s.client.Coupon.Update().
		Where(coupon.IDIn(ids...))

	if req.Active != nil {
		updateBuilder.SetActive(*req.Active)
	}

	_, err := updateBuilder.Save(ctx)
	return err
}

func (s *CouponServiceImpl) BatchDeleteCoupons(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := s.client.Coupon.Delete().
		Where(coupon.IDIn(ids...)).
		Exec(ctx)
	return err
}

func (s *CouponServiceImpl) SearchCoupons(ctx context.Context, req model.CouponSearchReq) ([]*ent.Coupon, int, error) {
	query := s.client.Coupon.Query()

	if req.Keyword != "" {
		query.Where(coupon.Or(
			coupon.NameContains(req.Keyword),
			coupon.CodeContains(req.Keyword),
			coupon.DescriptionContains(req.Keyword),
		))
	}
	if req.Type != nil {
		query.Where(coupon.Type(*req.Type))
	}
	if req.Active != nil {
		query.Where(coupon.Active(*req.Active))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	coupons, err := query.
		Order(ent.Desc(coupon.FieldCreatedAt)).
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		All(ctx)

	return coupons, total, err
}
