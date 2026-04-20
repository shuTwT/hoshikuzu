package menu

import (
	"context"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/menu"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
)

type MenuService interface {
	QueryMenu(c *fiber.Ctx, id int) (*ent.Menu, error)
	QueryMenuList(c *fiber.Ctx) ([]model.MenuResp, error)
	QueryMenuPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Menu, error)
	CreateMenu(c context.Context, createReq model.MenuCreateReq) (*ent.Menu, error)
	UpdateMenu(c *fiber.Ctx, id int, updateReq model.MenuUpdateReq) (*ent.Menu, error)
	DeleteMenu(c *fiber.Ctx, id int) error
}

type MenuServiceImpl struct {
	client *ent.Client
}

func NewMenuServiceImpl(client *ent.Client) *MenuServiceImpl {
	return &MenuServiceImpl{client: client}
}

func (s *MenuServiceImpl) QueryMenu(c *fiber.Ctx, id int) (*ent.Menu, error) {
	m, err := s.client.Menu.Query().
		Where(menu.ID(id)).
		Only(c.Context())
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MenuServiceImpl) QueryMenuList(c *fiber.Ctx) ([]model.MenuResp, error) {
	menus, err := s.client.Menu.Query().
		Order(ent.Asc(menu.FieldSortOrder), ent.Desc(menu.FieldID)).
		All(c.Context())
	if err != nil {
		return nil, err
	}

	var resp []model.MenuResp
	for _, m := range menus {
		resp = append(resp, model.MenuResp{
			ID:        m.ID,
			Name:      m.Name,
			Title:     m.Title,
			Path:      m.Path,
			Icon:      m.Icon,
			ParentID:  m.ParentID,
			SortOrder: m.SortOrder,
			Visible:   m.Visible,
			Target:    m.Target,
		})
	}

	return resp, nil
}

func (s *MenuServiceImpl) QueryMenuPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.Menu, error) {
	count, err := s.client.Menu.Query().Count(c.UserContext())
	if err != nil {
		return 0, nil, err
	}

	menus, err := s.client.Menu.Query().
		Order(ent.Asc(menu.FieldSortOrder), ent.Desc(menu.FieldID)).
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		return 0, nil, err
	}

	return count, menus, nil
}

func (s *MenuServiceImpl) CreateMenu(c context.Context, createReq model.MenuCreateReq) (*ent.Menu, error) {
	m, err := s.client.Menu.Create().
		SetName(createReq.Name).
		SetNillableTitle(createReq.Title).
		SetNillablePath(createReq.Path).
		SetNillableIcon(createReq.Icon).
		SetNillableParentID(createReq.ParentID).
		SetNillableSortOrder(createReq.SortOrder).
		SetNillableVisible(createReq.Visible).
		SetNillableTarget(createReq.Target).
		Save(c)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MenuServiceImpl) UpdateMenu(c *fiber.Ctx, id int, updateReq model.MenuUpdateReq) (*ent.Menu, error) {
	update := s.client.Menu.UpdateOneID(id)

	if updateReq.Name != nil {
		update.SetNillableName(updateReq.Name)
	}

	if updateReq.Title != nil {
		update.SetNillableTitle(updateReq.Title)
	}

	if updateReq.Path != nil {
		update.SetNillablePath(updateReq.Path)
	}

	if updateReq.Icon != nil {
		update.SetNillableIcon(updateReq.Icon)
	}

	if updateReq.ParentID != nil {
		update.SetNillableParentID(updateReq.ParentID)
	}

	if updateReq.SortOrder != nil {
		update.SetNillableSortOrder(updateReq.SortOrder)
	}

	if updateReq.Visible != nil {
		update.SetNillableVisible(updateReq.Visible)
	}

	if updateReq.Target != nil {
		update.SetNillableTarget(updateReq.Target)
	}

	updatedMenu, err := update.Save(c.Context())
	if err != nil {
		return nil, err
	}

	return updatedMenu, nil
}

func (s *MenuServiceImpl) DeleteMenu(c *fiber.Ctx, id int) error {
	err := s.client.Menu.DeleteOneID(id).Exec(c.Context())
	if err != nil {
		return err
	}
	return nil
}
