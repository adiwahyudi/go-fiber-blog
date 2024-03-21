package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go-blog/internal/model"
	"strings"
	"time"
)

type Middleware struct {
	Config *viper.Viper
	Log    *logrus.Logger
}

func NewMiddleware(v *viper.Viper, l *logrus.Logger) *Middleware {
	return &Middleware{
		Config: v,
		Log:    l,
	}
}

func (m *Middleware) ValidateJWT(ctx *fiber.Ctx) error {
	var token string
	authorization := ctx.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		token = strings.TrimPrefix(authorization, "Bearer ")
	}

	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Token Empty")
	}

	tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(m.Config.GetString("JWT_SECRET")), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "A "+err.Error())

	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "B "+err.Error())
	}

	subClaims := claims["auth"].(string)
	auth := new(model.Auth)

	if err := json.Unmarshal([]byte(subClaims), auth); err != nil {
		return fiber.ErrInternalServerError
	}

	ctx.Locals("auth", auth)
	return ctx.Next()
}

func (m *Middleware) BasicAuth(c *fiber.Ctx) error {
	config := basicauth.Config{
		Users: map[string]string{
			m.Config.GetString("BASIC_AUTH_USERNAME"): m.Config.GetString("BASIC_AUTH_PASSWORD"),
		},
		Unauthorized: func(ctx *fiber.Ctx) error {
			return fiber.ErrForbidden
		},
	}
	return basicauth.New(config)(c)
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}

func (m *Middleware) GenerateToken(auth *model.Auth) (string, error) {
	jwtSecret := m.Config.GetString("JWT_SECRET")
	authJSON, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"auth": string(authJSON),
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
