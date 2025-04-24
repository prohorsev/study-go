package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {

	file, err := os.OpenFile(".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	webApp := fiber.New(fiber.Config{
		ReadBufferSize: 16 * 1024})

	webApp.Use(requestid.New())
	webApp.Use(logger.New(logger.Config{
		Format: "${locals:requestid}: ${method} ${path} - ${status}\n",
		Output: file,
	}))

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	fooGroup := webApp.Group("/foo")
	fooGroup.Use(limiter.New(limiter.Config{
		Max:        1,
		Expiration: 2 * time.Second,
	}))
	fooGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	barGroup := webApp.Group("/bar")
	barGroup.Use(limiter.New(limiter.Config{
		Max:        1,
		Expiration: 2 * time.Second,
	}))
	barGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
