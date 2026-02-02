package auth

import (
	"context"
	"errors"
	"time"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/user"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResp, error)
}

type AuthServiceImpl struct {
	client *ent.Client
}

func NewAuthServiceImpl(client *ent.Client) *AuthServiceImpl {
	return &AuthServiceImpl{client: client}
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResp, error) {
	u, err := s.client.User.Query().
		Where(user.EmailEQ(req.Email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("找不到该用户")
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	expires := time.Now().Add(time.Hour * 24).UnixMilli()
	claims := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
		"exp":   expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetString(config.AUTH_TOKEN_SECRET)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &model.LoginResp{
		AccessToken: t,
		Expires:     expires,
		Username:    u.Name,
		Roles:       []string{"admin"},
	}, nil
}
