package role

import (
	"context"
	"fmt"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/role"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

type RoleService interface {
	QueryRoleList(c context.Context) ([]*ent.Role, error)
	QueryRolePage(c context.Context, req model.PageQuery) ([]*ent.Role, int, error)
	QueryRole(c context.Context, id int) (*ent.Role, error)
	CreateRole(c context.Context, req model.RoleCreateReq) (*ent.Role, error)
	UpdateRole(c context.Context, id int, req model.RoleUpdateReq) (*ent.Role, error)
	DeleteRole(c context.Context, id int) error
}

type RoleServiceImpl struct {
	client *ent.Client
}

func NewRoleServiceImpl(client *ent.Client) *RoleServiceImpl {
	return &RoleServiceImpl{client: client}
}

func (s *RoleServiceImpl) QueryRoleList(c context.Context) ([]*ent.Role, error) {
	return s.client.Role.Query().
		Order(ent.Desc(role.FieldID)).
		All(c)
}

func (s *RoleServiceImpl) QueryRolePage(c context.Context, req model.PageQuery) ([]*ent.Role, int, error) {
	count, err := s.client.Role.Query().Count(c)
	if err != nil {
		return nil, 0, err
	}
	roles, err := s.client.Role.Query().
		Order(ent.Desc(role.FieldID)).
		Limit(req.Size).
		Offset((req.Page - 1) * req.Size).
		All(c)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func (s *RoleServiceImpl) QueryRole(c context.Context, id int) (*ent.Role, error) {
	return s.client.Role.Query().Where(role.IDEQ(id)).First(c)

}

func (s *RoleServiceImpl) CreateRole(c context.Context, req model.RoleCreateReq) (*ent.Role, error) {
	// 检查角色代码是否已存在
	exists, err := s.client.Role.Query().
		Where(role.CodeEQ(req.Code)).
		Exist(c)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("Role code already exists")
	}

	role, err := s.client.Role.Create().
		SetName(req.Name).
		SetCode(req.Code).
		Save(c)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleServiceImpl) UpdateRole(c context.Context, id int, req model.RoleUpdateReq) (*ent.Role, error) {
	// 检查角色代码是否已存在
	exists, err := s.client.Role.Query().
		Where(role.CodeEQ(req.Code)).
		Exist(c)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf(
			"Role code already exists",
		)
	}
	// 开始构建更新
	update := s.client.Role.UpdateOneID(id)

	// 如果提供了新名称
	if req.Name != "" {
		update.SetName(req.Name)
	}

	// 如果提供了新代码
	if req.Code != "" {
		update.SetCode(req.Code)
	}

	update.SetDescription(req.Description)

	newRole, err := update.Save(c)
	if err != nil {
		return nil, err
	}

	return newRole, nil
}

func (s *RoleServiceImpl) DeleteRole(c context.Context, id int) error {
	return s.client.Role.DeleteOneID(id).Exec(c)
}
