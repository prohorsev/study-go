package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"slices"
)

type (
	BinarySearchRequest struct {
		Numbers []int `json:"numbers"`
		Target  int   `json:"target"`
	}

	BinarySearchResponse struct {
		TargetIndex int    `json:"target_index"`
		Error       string `json:"error,omitempty"`
	}
)

const targetNotFound = -1

func main() {
	webApp := fiber.New(fiber.Config{
		ReadBufferSize: 16 * 1024})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	webApp.Post("/search", func(c *fiber.Ctx) error {
		var request BinarySearchRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(BinarySearchResponse{
				TargetIndex: -1,
				Error:       "Invalid JSON",
			})
		}
		targetIndex := slices.Index(request.Numbers, request.Target)
		if targetIndex == targetNotFound {
			return c.Status(fiber.StatusNotFound).JSON(BinarySearchResponse{
				TargetIndex: -1,
				Error:       "Target was not found",
			})
		}

		return c.JSON(BinarySearchResponse{
			TargetIndex: targetIndex,
			Error:       "",
		})
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
