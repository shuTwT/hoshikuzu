package middleware

import (
	"strings"

	"github.com/shuTwT/hoshikuzu/ent"
	"github.com/shuTwT/hoshikuzu/ent/personalaccesstoken"
	"github.com/shuTwT/hoshikuzu/pkg/config"
	"github.com/shuTwT/hoshikuzu/pkg/domain/model"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected 保护需要认证的路由
func Protected() fiber.Handler {
	// 从配置中获取密钥
	secret := config.GetString(config.AUTH_TOKEN_SECRET)
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: "HS256",
			Key:    []byte(secret)}, // 注意：在生产环境中应该使用环境变量存储密钥
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccess,
	})
}

// jwtError 处理JWT错误
func jwtError(c *fiber.Ctx, err error) error {
	c.Locals("authFailed", true)
	c.Locals("authType", "jwt")
	return nil
}

// jwtSuccess 处理JWT验证成功
func jwtSuccess(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	c.Locals("userId", claims["id"])
	c.Locals("userEmail", claims["email"])
	c.Locals("userName", claims["name"])
	c.Locals("authSuccess", true)
	c.Locals("authType", "jwt")

	return c.Next()
}

// FlexibleAuth 支持JWT或个人访问令牌任意一种认证方式
func FlexibleAuth(dbClient *ent.Client) fiber.Handler {
	var client = dbClient
	return func(c *fiber.Ctx) error {
		authToken := c.Get("Authorization")

		if authToken == "" {
			return c.JSON(model.NewError(fiber.StatusUnauthorized, "Authentication required"))
		}

		tokenString := strings.TrimPrefix(authToken, "Bearer ")
		if tokenString == "" {
			return c.JSON(model.NewError(fiber.StatusUnauthorized, "Invalid token format"))
		}

		secret := config.GetString(config.AUTH_TOKEN_SECRET)
		patSecret := config.GetString(config.AUTH_PAT_SECRET)

		hasAuth := false

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(model.NewError(fiber.StatusInternalServerError, "Invalid claims type"))
			}
			c.Locals("userId", claims["id"])
			c.Locals("userEmail", claims["email"])
			c.Locals("userName", claims["name"])
			c.Locals("authSuccess", true)
			c.Locals("authType", "jwt")
			hasAuth = true
		}

		if !hasAuth {
			token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(patSecret), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)

				pat, err := client.PersonalAccessToken.Query().
					Where(personalaccesstoken.Token(tokenString)).
					First(c.Context())

				if err == nil {
					if pat.Expires.IsZero() || pat.Expires.After(c.Context().Time()) {
						c.Locals("userId", claims["id"])
						c.Locals("userEmail", claims["email"])
						c.Locals("userName", claims["name"])
						c.Locals("patId", pat.ID)
						c.Locals("authSuccess", true)
						c.Locals("authType", "pat")
						hasAuth = true
					}
				}
			}
		}

		if hasAuth {
			return c.Next()
		}

		return c.JSON(model.NewError(fiber.StatusUnauthorized, "Authentication required"))
	}
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *fiber.Ctx) *model.LoginUser {
	if c.Locals("userId") == nil {
		return nil
	}
	return &model.LoginUser{
		ID:       int(c.Locals("userId").(float64)),
		Email:    c.Locals("userEmail").(string),
		Username: c.Locals("userName").(string),
	}
}

// PersonalAccessTokenProtected 使用个人访问令牌保护需要认证的路由
func PersonalAccessTokenProtected(dbClient *ent.Client) fiber.Handler {
	var client = dbClient
	return func(c *fiber.Ctx) error {
		authToken := c.Get("Authorization")
		if authToken == "" {
			return c.JSON(model.NewError(fiber.StatusUnauthorized, "Authentication required"))
		}
		tokenString := strings.TrimPrefix(authToken, "Bearer ")
		if tokenString == "" {
			return c.JSON(model.NewError(fiber.StatusUnauthorized, "Invalid token format"))
		}
		patSecret := config.GetString(config.AUTH_PAT_SECRET)

		hasAuth := false

		if !hasAuth {
			token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(patSecret), nil
			})

			if err == nil && token.Valid {
				claims := token.Claims.(jwt.MapClaims)

				pat, err := client.PersonalAccessToken.Query().
					Where(personalaccesstoken.Token(tokenString)).
					First(c.Context())

				if err == nil {
					if pat.Expires.IsZero() || pat.Expires.After(c.Context().Time()) {
						c.Locals("userId", claims["id"])
						c.Locals("userEmail", claims["email"])
						c.Locals("userName", claims["name"])
						c.Locals("patId", pat.ID)
						c.Locals("authSuccess", true)
						c.Locals("authType", "pat")
						hasAuth = true
					}
				}
			}
		}
		if hasAuth {
			return c.Next()
		}

		return c.JSON(model.NewError(fiber.StatusUnauthorized, "Authentication required"))
	}
}

// patError 处理个人访问令牌JWT错误
func patError(c *fiber.Ctx, err error) error {
	c.Locals("authFailed", true)
	c.Locals("authType", "pat")
	return nil
}
