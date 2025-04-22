package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID      int64
	Email   string
	Age     int
	Country string
}

var users = map[int64]User{}

type (
	CreateUserRequest struct {
		ID      int64  `json:"id" validate:"required,min=1"`
		Email   string `json:"email" validate:"required,email"`
		Age     int    `json:"age" validate:"required,min=18,max=130"`
		Country string `json:"country" validate:"required,allowable_country"`
	}
)

var allowableCountries = []string{
	"USA",
	"Germany",
	"France",
}

func main() {
	webApp := fiber.New(fiber.Config{ReadBufferSize: 16 * 1024})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("1334414")
	})

	validate := validator.New()
	vErr := validate.RegisterValidation("allowable_country", func(fl validator.FieldLevel) bool {
		text := fl.Field().String()
		for _, country := range allowableCountries {
			if text == country {
				return true
			}
		}

		return false
	})
	if vErr != nil {
		logrus.Fatal("register validation ", vErr)
	}

	webApp.Post("/users", func(ctx *fiber.Ctx) error {
		var req CreateUserRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}
		err := validate.Struct(req)
		if err != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
		}

		return ctx.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
