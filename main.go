package main

import (
	"errors"
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type (
	SignUpRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SignInResponse struct {
		JWTToken string `json:"jwt_token"`
	}

	ProfileResponse struct {
		Email string `json:"email"`
	}

	User struct {
		Email    string
		password string
	}

	AuthHandler struct{}

	UserHandler struct{}
)

var (
	webApiPort = ":8080"

	users = map[string]User{}

	secretKey = []byte("qwerty123456")

	contextKeyUser = "user"
)

func main() {
	webApp := fiber.New(fiber.Config{
		ReadBufferSize: 16 * 1024})

	authHandler := &AuthHandler{}
	userHandler := &UserHandler{}

	publicGroup := webApp.Group("")
	publicGroup.Post("/signup", authHandler.Signup)
	publicGroup.Post("/signin", authHandler.Signin)

	authorizedGroup := webApp.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: secretKey,
		},
		ContextKey: contextKeyUser,
	}))
	authorizedGroup.Get("/profile", userHandler.Profile)

	logrus.Fatal(webApp.Listen(webApiPort))
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	req := SignUpRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if _, exists := users[req.Email]; exists {
		return c.SendStatus(fiber.StatusConflict)
	}

	users[req.Email] = User{
		Email:    req.Email,
		password: req.Password,
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) Signin(c *fiber.Ctx) error {
	req := SignInRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	user, exists := users[req.Email]
	if !exists {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if user.password != req.Password {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(secretKey)
	if err != nil {
		logrus.WithError(err).Error("JWT token signing")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(SignInResponse{JWTToken: t})
}

func (h *UserHandler) Profile(c *fiber.Ctx) error {
	jwtPayload, ok := jwtPayloadFromRequest(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userInfo, ok := users[jwtPayload["sub"].(string)]
	if !ok {
		return errors.New("user not found")
	}

	return c.JSON(ProfileResponse{
		Email: userInfo.Email,
	})
}

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, bool) {
	jwtToken, ok := c.Context().Value(contextKeyUser).(*jwt.Token)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_context_value": c.Context().Value(contextKeyUser),
		}).Error("wrong type of JWT token in context")
		return nil, false
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"jwt_token_claims": jwtToken.Claims,
		}).Error("wrong type of JWT token claims")
		return nil, false
	}

	return payload, true
}
