package initialize

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/role"
	"github.com/shuTwT/hoshikuzu/ent/setting"
	user_service "github.com/shuTwT/hoshikuzu/internal/services/system/user"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"
)

const systemInitializedKey = "system_initialized"

// 初始化时创建的内置角色，ID 固定，超级管理员账户关联 ID 为 1 的角色。
var builtinRoles = []struct {
	id   int
	name string
	code string
}{
	{1, "超级管理员", "superAdmin"},
	{2, "访客", "guest"},
	{3, "普通用户", "common"},
}

type InitializeService interface {
	Initialize(ctx context.Context, req *model.InitializeRequest) error
}

type InitializeServiceImpl struct {
	client      *ent.Client
	userService user_service.UserService
	// mu 串行化同一进程内的并发初始化请求，配合事务内的状态检查保证只初始化一次
	mu sync.Mutex
}

func NewInitializeServiceImpl(client *ent.Client, userService user_service.UserService) *InitializeServiceImpl {
	return &InitializeServiceImpl{client: client, userService: userService}
}

// Initialize 在一个事务中完成系统初始化：检查初始化标记、创建内置角色、创建超级管理员并写入初始化标记。
// 任一步失败都会整体回滚，不会留下"角色已建、管理员未建"的半初始化状态。
func (s *InitializeServiceImpl) Initialize(ctx context.Context, req *model.InitializeRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := s.initialize(ctx, tx, req); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: rollback failed: %v", err, rerr)
		}
		return err
	}
	return tx.Commit()
}

func (s *InitializeServiceImpl) initialize(ctx context.Context, tx *ent.Tx, req *model.InitializeRequest) error {
	initialized, err := tx.Setting.Query().Where(setting.KeyEQ(systemInitializedKey)).Exist(ctx)
	if err != nil {
		return fmt.Errorf("检查系统初始化状态失败: %w", err)
	}
	if initialized {
		return model.ErrSystemInitialized
	}

	for _, r := range builtinRoles {
		exists, err := tx.Role.Query().Where(role.ID(r.id)).Exist(ctx)
		if err != nil {
			return fmt.Errorf("检查内置角色失败: %w", err)
		}
		if exists {
			continue
		}
		if _, err := tx.Role.Create().
			SetID(r.id).
			SetName(r.name).
			SetCode(r.code).
			SetIsDefault(true).
			Save(ctx); err != nil {
			return fmt.Errorf("创建内置角色失败: %w", err)
		}
	}

	adminUser, err := s.userService.CreateUserWithTx(ctx, tx, model.UserCreateReq{
		Name:     req.AdminUsername,
		Password: req.AdminPassword,
		Email:    req.AdminEmail,
		RoleID:   builtinRoles[0].id,
	})
	if err != nil {
		return fmt.Errorf("创建超级管理员账户失败: %w", err)
	}

	if _, err := tx.Setting.Create().SetKey(systemInitializedKey).SetValue("true").Save(ctx); err != nil {
		return fmt.Errorf("标记系统初始化状态失败: %w", err)
	}

	slog.Info("系统初始化完成，超级管理员账户已创建", "user_id", adminUser.ID, "user_name", adminUser.Name)
	return nil
}
