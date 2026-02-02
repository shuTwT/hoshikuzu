package tag

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/post"
	"github.com/shuTwT/hoshikuzu/ent/tag"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type TagService interface {
	QueryTag(c *fiber.Ctx, id int) (*ent.Tag, error)
	QueryTagList(c *fiber.Ctx) ([]model.TagResp, error)
	QueryTagPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Tag, error)
	CreateTag(c context.Context, createReq model.TagCreateReq) (*ent.Tag, error)
	UpdateTag(c *fiber.Ctx, id int, updateReq model.TagUpdateReq) (*ent.Tag, error)
	DeleteTag(c *fiber.Ctx, id int) error
}

type TagServiceImpl struct {
	client *ent.Client
}

func NewTagServiceImpl(client *ent.Client) *TagServiceImpl {
	return &TagServiceImpl{client: client}
}

func (s *TagServiceImpl) QueryTag(c *fiber.Ctx, id int) (*ent.Tag, error) {
	tag, err := s.client.Tag.Query().
		Where(tag.ID(id)).
		Only(c.Context())
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagServiceImpl) QueryTagList(c *fiber.Ctx) ([]model.TagResp, error) {
	tags, err := s.client.Tag.Query().
		Order(ent.Desc(tag.FieldID)).
		All(c.Context())
	if err != nil {
		return nil, err
	}

	var resp []model.TagResp
	for _, t := range tags {
		postCount, err := s.client.Post.Query().
			Where(post.HasTagsWith(tag.ID(t.ID))).
			Count(c.Context())
		if err != nil {
			return nil, err
		}

		resp = append(resp, model.TagResp{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Slug:        t.Slug,
			Color:       t.Color,
			SortOrder:   t.SortOrder,
			Active:      t.Active,
			PostCount:   postCount,
		})
	}

	return resp, nil
}

func (s *TagServiceImpl) QueryTagPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Tag, error) {
	count, err := s.client.Tag.Query().Count(c.UserContext())
	if err != nil {
		return 0, nil, err
	}

	tags, err := s.client.Tag.Query().
		Order(ent.Desc(tag.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		return 0, nil, err
	}

	return count, tags, nil
}

func (s *TagServiceImpl) CreateTag(c context.Context, createReq model.TagCreateReq) (*ent.Tag, error) {
	tag, err := s.client.Tag.Create().
		SetName(createReq.Name).
		SetNillableDescription(createReq.Description).
		SetNillableSlug(createReq.Slug).
		SetNillableColor(createReq.Color).
		SetNillableSortOrder(createReq.SortOrder).
		SetNillableActive(createReq.Active).
		Save(c)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagServiceImpl) UpdateTag(c *fiber.Ctx, id int, updateReq model.TagUpdateReq) (*ent.Tag, error) {
	update := s.client.Tag.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetNillableName(updateReq.Name)
	}

	if updateReq.Description != nil {
		update.SetNillableDescription(updateReq.Description)
	}

	if updateReq.Slug != nil {
		update.SetNillableSlug(updateReq.Slug)
	}

	if updateReq.Color != nil {
		update.SetNillableColor(updateReq.Color)
	}

	if updateReq.SortOrder != nil {
		update.SetNillableSortOrder(updateReq.SortOrder)
	}

	if updateReq.Active != nil {
		update.SetNillableActive(updateReq.Active)
	}

	updatedTag, err := update.Save(c.Context())
	if err != nil {
		return nil, err
	}

	return updatedTag, nil
}

func (s *TagServiceImpl) DeleteTag(c *fiber.Ctx, id int) error {
	err := s.client.Tag.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return err
	}
	return nil
}
