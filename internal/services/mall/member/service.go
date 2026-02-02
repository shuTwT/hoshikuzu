package member

import (
	"context"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/member"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MemberService interface {
	QueryMember(c *fiber.Ctx, userId int) (*ent.Member, error)
	QueryMemberPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Member, error)
	CreateMember(c context.Context, userId int, memberLevel int, memberNo string) (*ent.Member, error)
	UpdateMember(c *fiber.Ctx, id int, updateReq *model.MemberUpdateReq) (*ent.Member, error)
	DeleteMember(c *fiber.Ctx, id int) error
}

type MemberServiceImpl struct {
	client *ent.Client
}

func NewMemberServiceImpl(client *ent.Client) *MemberServiceImpl {
	return &MemberServiceImpl{client: client}
}

func (s *MemberServiceImpl) QueryMember(c *fiber.Ctx, userId int) (*ent.Member, error) {
	m, err := s.client.Member.Query().
		Where(member.UserIDEQ(userId)).
		Only(c.Context())
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MemberServiceImpl) QueryMemberPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Member, error) {
	count, err := s.client.Member.Query().Count(c.UserContext())
	if err != nil {
		return 0, nil, err
	}

	members, err := s.client.Member.Query().
		Order(ent.Desc(member.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		return 0, nil, err
	}

	return count, members, nil
}

func (s *MemberServiceImpl) CreateMember(c context.Context, userId int, memberLevel int, memberNo string) (*ent.Member, error) {
	m, err := s.client.Member.Create().
		SetUserID(userId).
		SetMemberLevel(memberLevel).
		SetMemberNo(memberNo).
		Save(c)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MemberServiceImpl) UpdateMember(c *fiber.Ctx, id int, updateReq *model.MemberUpdateReq) (*ent.Member, error) {
	update := s.client.Member.UpdateOneID(id)

	if updateReq.MemberLevel != nil {
		update.SetMemberLevel(*updateReq.MemberLevel)
	}

	if updateReq.ExpireTime != nil {
		expireTime, err := time.Parse(time.RFC3339, *updateReq.ExpireTime)
		if err == nil {
			update.SetExpireTime(model.ParseTime(expireTime).Time())
		}
	}

	if updateReq.Points != nil {
		update.SetPoints(*updateReq.Points)
	}

	if updateReq.TotalSpent != nil {
		update.SetTotalSpent(*updateReq.TotalSpent)
	}

	if updateReq.OrderCount != nil {
		update.SetOrderCount(*updateReq.OrderCount)
	}

	if updateReq.Active != nil {
		update.SetActive(*updateReq.Active)
	}

	if updateReq.Remark != nil {
		update.SetRemark(*updateReq.Remark)
	}

	updatedMember, err := update.Save(c.Context())
	if err != nil {
		return nil, err
	}

	return updatedMember, nil
}

func (s *MemberServiceImpl) DeleteMember(c *fiber.Ctx, id int) error {
	err := s.client.Member.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return err
	}
	return nil
}
