package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/personalaccesstoken"
	"github.com/shuTwT/hoshikuzu/ent/user"
	"github.com/shuTwT/hoshikuzu/pkg/cache"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUser(c *fiber.Ctx) ([]*ent.User, error)
	ListUserPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.User, error)
	QueryUserById(ctx context.Context, id int) (*ent.User, error)
	CreateUser(ctx context.Context, req model.UserCreateReq) (*ent.User, error)
	UpdateUser(ctx context.Context, id int, req model.UserUpdateReq) (*ent.User, error)
	GetUserCount(ctx context.Context) (int, error)
	DeleteUser(ctx context.Context, id int) error
	GetPersonalAccessTokenList(ctx context.Context, userId int) ([]*ent.PersonalAccessToken, error)
	GetPersonalAccessToken(ctx context.Context, userId int, id int) (*ent.PersonalAccessToken, error)
	CreatePersonalAccessToken(ctx context.Context, id int, req model.PersonalAccessTokenCreateReq) (*ent.PersonalAccessToken, error)
	SearchUsers(ctx context.Context, req model.UserSearchReq) ([]*model.UserSearchResp, int, error)
}

type UserServiceImpl struct {
	client *ent.Client
}

func NewUserServiceImpl(client *ent.Client) *UserServiceImpl {
	return &UserServiceImpl{client: client}
}

func (s *UserServiceImpl) ListUser(c *fiber.Ctx) ([]*ent.User, error) {
	client := s.client
	users, err := client.User.Query().All(c.Context())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserServiceImpl) ListUserPage(c *fiber.Ctx, pageQuery model.PageQuery) (int, []*ent.User, error) {
	count, err := s.client.User.Query().Count(c.UserContext())

	if err != nil {
		c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
		return 0, nil, err
	}

	users, err := s.client.User.Query().
		Offset((pageQuery.Page - 1) * pageQuery.Size).
		Limit(pageQuery.Size).
		All(c.Context())
	if err != nil {
		c.JSON(model.NewError(fiber.StatusInternalServerError,
			err.Error(),
		))
		return 0, nil, err
	}

	return count, users, err
}

func (s *UserServiceImpl) QueryUserById(ctx context.Context, id int) (*ent.User, error) {
	user, err := s.client.User.Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("User not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req model.UserCreateReq) (*ent.User, error) {

	// 检查邮箱是否已存在
	exists, err := s.client.User.Query().
		Where(user.EmailEQ(req.Email)).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("email %s already exists", req.Email)
	}

	// 如果提供了手机号，检查手机号是否已存在
	if req.PhoneNumber != "" {
		exists, err = s.client.User.Query().
			Where(user.PhoneNumberEQ(req.PhoneNumber)).
			Exist(ctx)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("phone number %s already exists", req.PhoneNumber)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 使用事务创建用户和钱包
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 创建用户
	createdUser, err := tx.User.Create().
		SetName(req.Name).
		SetPassword(string(hashedPassword)).
		SetPhoneNumber(req.PhoneNumber).
		SetEmail(req.Email).
		SetRoleID(req.RoleID).
		Save(ctx)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 为用户创建钱包
	_, err = tx.Wallet.Create().
		SetUserID(createdUser.ID).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 为用户创建会员
	memberNo := fmt.Sprintf("M%06d", createdUser.ID)
	_, err = tx.Member.Create().
		SetUserID(createdUser.ID).
		SetMemberLevel(1).
		SetMemberNo(memberNo).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id int, req model.UserUpdateReq) (*ent.User, error) {
	// 开始构建更新
	update := s.client.User.UpdateOneID(id)

	// 如果提供了新名称
	if req.Name != "" {
		update.SetName(req.Name)
	}

	// 如果提供了新手机号
	if req.PhoneNumber != "" {
		// 检查手机号是否已被其他用户使用
		var exists bool
		exists, err := s.client.User.Query().
			Where(
				user.And(
					user.PhoneNumberEQ(req.PhoneNumber),
					user.IDNEQ(id),
				),
			).
			Exist(ctx)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("Phone number already exists")
		}
		update.SetPhoneNumber(req.PhoneNumber)
	}

	// 如果提供了新密码
	if req.Password != "" {
		var hashedPassword []byte
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("Failed to hash password")
		}
		update.SetPassword(string(hashedPassword))
	}

	// 执行更新
	updatedUser, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserServiceImpl) GetUserCount(ctx context.Context) (int, error) {
	count, err := s.client.User.Query().Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id int) error {
	err := s.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("User not found")
		}
		return err
	}

	return nil
}

func (s *UserServiceImpl) GetPersonalAccessTokenList(ctx context.Context, userId int) ([]*ent.PersonalAccessToken, error) {
	tokens, err := s.client.PersonalAccessToken.Query().Where(personalaccesstoken.UserID(userId)).All(ctx)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *UserServiceImpl) GetPersonalAccessToken(ctx context.Context, userId int, id int) (*ent.PersonalAccessToken, error) {
	token, err := s.client.PersonalAccessToken.Query().Where(personalaccesstoken.UserIDEQ(userId), personalaccesstoken.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserServiceImpl) CreatePersonalAccessToken(ctx context.Context, userId int, req model.PersonalAccessTokenCreateReq) (*ent.PersonalAccessToken, error) {
	u, _ := s.client.User.Query().
		Where(user.IDEQ(userId)).
		Only(ctx)

	claims := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"exp":   req.Expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetString(config.AUTH_PAT_SECRET)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	pat, err := s.client.PersonalAccessToken.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetExpires(req.Expires.Time()).
		SetToken(t).
		SetUserID(userId).Save(ctx)
	if err != nil {
		return nil, err
	}
	return pat, nil
}

func (s *UserServiceImpl) SearchUsers(ctx context.Context, req model.UserSearchReq) ([]*model.UserSearchResp, int, error) {
	cacheKey := fmt.Sprintf("user:search:%s:%d:%d", req.Keyword, req.Page, req.Size)

	if cached, found := cache.GetCache().Get(cacheKey); found {
		if result, ok := cached.([]*model.UserSearchResp); ok {
			return result, len(result), nil
		}
	}

	keyword := strings.ToLower(req.Keyword)

	users, err := s.client.User.Query().
		Where(
			user.Or(
				user.NameContains(keyword),
				user.EmailContains(keyword),
			),
		).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	var results []*model.UserSearchResp

	for _, u := range users {
		relevance := s.calculateUserRelevance(u, keyword)
		if relevance > 0 {
			results = append(results, &model.UserSearchResp{
				ID:     u.ID,
				Name:   u.Name,
				Email:  u.Email,
				RoleID: u.RoleID,
			})
		}
	}

	total := len(results)

	start := (req.Page - 1) * req.Size
	end := start + req.Size

	if start >= total {
		return []*model.UserSearchResp{}, total, nil
	}
	if end > total {
		end = total
	}

	pagedResults := results[start:end]

	cache.GetCache().Set(cacheKey, pagedResults, 5*time.Minute)

	return pagedResults, total, nil
}

func (s *UserServiceImpl) calculateUserRelevance(u *ent.User, keyword string) float64 {
	var relevance float64 = 0

	name := strings.ToLower(u.Name)
	email := strings.ToLower(u.Email)

	if strings.Contains(name, keyword) {
		if name == keyword {
			relevance += 10.0
		} else if strings.HasPrefix(name, keyword) {
			relevance += 8.0
		} else {
			relevance += 5.0
		}
	}

	if strings.Contains(email, keyword) {
		relevance += 3.0
	}

	return relevance
}
