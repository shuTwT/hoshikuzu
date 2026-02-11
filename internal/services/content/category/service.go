package category

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/category"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type CategoryService interface {
	QueryCategory(c context.Context, id int) (*ent.Category, error)
	QueryCategoryList(c context.Context) ([]*ent.Category, error)
	QueryCategoryPage(c context.Context, pageQuery model.PageQuery) (int, []*ent.Category, error)
	CreateCategory(c context.Context, createReq model.CategoryCreateReq) (*ent.Category, error)
	UpdateCategory(c context.Context, id int, updateReq *model.CategoryUpdateReq) (*ent.Category, error)
	DeleteCategory(c context.Context, id int) error
}

type CategoryServiceImpl struct {
	client *ent.Client
}

func NewCategoryServiceImpl(client *ent.Client) *CategoryServiceImpl {
	return &CategoryServiceImpl{client: client}
}

func (s *CategoryServiceImpl) QueryCategory(c context.Context, id int) (*ent.Category, error) {
	category, err := s.client.Category.Query().
		Where(category.ID(id)).
		Only(c)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) QueryCategoryList(c context.Context) ([]*ent.Category, error) {

	categories, err := s.client.Category.Query().
		All(c)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryServiceImpl) QueryCategoryPage(c context.Context, pageQuery model.PageQuery) (int, []*ent.Category, error) {
	count, err := s.client.Category.Query().Count(c)
	if err != nil {
		return 0, nil, err
	}

	categories, err := s.client.Category.Query().
		Order(ent.Desc(category.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c)
	if err != nil {
		return 0, nil, err
	}

	return count, categories, nil
}

func (s *CategoryServiceImpl) CreateCategory(c context.Context, createReq model.CategoryCreateReq) (*ent.Category, error) {
	category, err := s.client.Category.Create().
		SetName(createReq.Name).
		SetNillableDescription(createReq.Description).
		SetNillableSlug(createReq.Slug).
		SetSortOrder(createReq.SortOrder).
		SetNillableActive(createReq.Active).
		Save(c)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryServiceImpl) UpdateCategory(c context.Context, id int, updateReq *model.CategoryUpdateReq) (*ent.Category, error) {

	update := s.client.Category.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetName(*updateReq.Name)
	}

	if updateReq.Description != nil {
		update.SetNillableDescription(updateReq.Description)
	}

	if updateReq.Slug != nil {
		update.SetNillableSlug(updateReq.Slug)
	}

	if updateReq.SortOrder != nil {
		update.SetNillableSortOrder(updateReq.SortOrder)
	}

	if updateReq.Active != nil {
		update.SetNillableActive(updateReq.Active)
	}

	updatedCategory, err := update.Save(c)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

func (s *CategoryServiceImpl) DeleteCategory(c context.Context, id int) error {
	err := s.client.Category.DeleteOneID(id).Exec(c)
	if err != nil {
		return err
	}
	return nil
}
