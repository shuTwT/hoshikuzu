package memberlevel

import (
	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/memberlevel"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberLevelService interface {
	QueryMemberLevel(c *fiber.Ctx, id int) (*ent.MemberLevel, error)
	QueryMemberLevelList(c *fiber.Ctx) ([]*ent.MemberLevel, error)
	QueryMemberLevelPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.MemberLevel, error)
	CreateMemberLevel(c *fiber.Ctx, createReq *model.MemberLevelCreateReq) (*ent.MemberLevel, error)
	UpdateMemberLevel(c *fiber.Ctx, id int, updateReq *model.MemberLevelUpdateReq) (*ent.MemberLevel, error)
	DeleteMemberLevel(c *fiber.Ctx, id int) error
}

type MemberLevelServiceImpl struct {
	client *ent.Client
}

func NewMemberLevelServiceImpl(client *ent.Client) *MemberLevelServiceImpl {
	return &MemberLevelServiceImpl{client: client}
}

func (s *MemberLevelServiceImpl) QueryMemberLevel(c *fiber.Ctx, id int) (*ent.MemberLevel, error) {
	ml, err := s.client.MemberLevel.Query().
		Where(memberlevel.ID(id)).
		Only(c.Context())
	if err != nil {
		return nil, err
	}
	return ml, nil
}

func (s *MemberLevelServiceImpl) QueryMemberLevelList(c *fiber.Ctx) ([]*ent.MemberLevel, error) {
	memberLevels, err := s.client.MemberLevel.Query().
		Order(ent.Desc(memberlevel.FieldID)).
		All(c.Context())
	if err != nil {
		return nil, err
	}
	return memberLevels, nil
}

func (s *MemberLevelServiceImpl) QueryMemberLevelPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.MemberLevel, error) {
	count, err := s.client.MemberLevel.Query().Count(c.UserContext())
	if err != nil {
		return 0, nil, err
	}

	memberLevels, err := s.client.MemberLevel.Query().
		Order(ent.Desc(memberlevel.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		return 0, nil, err
	}

	return count, memberLevels, nil
}

func (s *MemberLevelServiceImpl) CreateMemberLevel(c *fiber.Ctx, createReq *model.MemberLevelCreateReq) (*ent.MemberLevel, error) {

	ml, err := s.client.MemberLevel.Create().
		SetName(createReq.Name).
		SetDescription(createReq.Description).
		SetLevel(createReq.Level).
		SetMinPoints(createReq.MinPoints).
		SetDiscountRate(createReq.DiscountRate).
		SetPrivileges(createReq.Privileges).
		SetIcon(createReq.Icon).
		SetSortOrder(createReq.SortOrder).
		Save(c.Context())
	if err != nil {
		return nil, err
	}
	return ml, nil
}

func (s *MemberLevelServiceImpl) UpdateMemberLevel(c *fiber.Ctx, id int, updateReq *model.MemberLevelUpdateReq) (*ent.MemberLevel, error) {
	update := s.client.MemberLevel.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetName(*updateReq.Name)
	}

	if updateReq.Description != nil {
		update.SetDescription(*updateReq.Description)
	}

	if updateReq.Level != nil {
		update.SetLevel(*updateReq.Level)
	}

	if updateReq.MinPoints != nil {
		update.SetMinPoints(*updateReq.MinPoints)
	}

	if updateReq.DiscountRate != nil {
		update.SetDiscountRate(*updateReq.DiscountRate)
	}

	if updateReq.Privileges != nil {
		update.SetPrivileges(*updateReq.Privileges)
	}

	if updateReq.Icon != nil {
		update.SetIcon(*updateReq.Icon)
	}

	if updateReq.Active != nil {
		update.SetActive(*updateReq.Active)
	}

	if updateReq.SortOrder != nil {
		update.SetSortOrder(*updateReq.SortOrder)
	}

	updatedMemberLevel, err := update.Save(c.Context())
	if err != nil {
		return nil, err
	}

	return updatedMemberLevel, nil
}

func (s *MemberLevelServiceImpl) DeleteMemberLevel(c *fiber.Ctx, id int) error {
	err := s.client.MemberLevel.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return err
	}
	return nil
}
