package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	SendPushNotificationRequest struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}

	PushNotification struct {
		Message string `json:"message"`
		UserID  int64  `json:"user_id"`
	}
)

var pushNotificationsQueue []PushNotification

func main() {
	webApp := fiber.New(fiber.Config{
		ReadTimeout:  300 * time.Millisecond,
		WriteTimeout: 300 * time.Millisecond,
	})
	webApp.Use(recover.New())
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/push/send", func(c *fiber.Ctx) error {
		var req SendPushNotificationRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
		}

		pushNotificationsQueue = append(pushNotificationsQueue, PushNotification{
			Message: req.Message,
			UserID:  req.UserID,
		})
		if len(pushNotificationsQueue) > 3 {
			panic("Queue is full")
		}

		return c.SendString("OK")
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
