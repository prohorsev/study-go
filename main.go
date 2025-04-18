package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/url"
)

type (
	CreateLinkRequest struct {
		External string `json:"external"`
		Internal string `json:"internal"`
	}

	GetLinkResponse struct {
		Internal string `json:"internal"`
	}
)

var links = make(map[string]string)

func main() {
	webApp := fiber.New(fiber.Config{
		ReadBufferSize: 16 * 1024})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	linkHandler := &LinkHandler{}
	webApp.Post("/links", linkHandler.CreateLink)
	webApp.Get("/links/:external", linkHandler.GetLink)

	logrus.Fatal(webApp.Listen(":8080"))
}

type LinkHandler struct{}

func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}
	links[request.External] = request.Internal

	return c.SendStatus(fiber.StatusOK)
}

func (h *LinkHandler) GetLink(c *fiber.Ctx) error {
	external, _ := url.QueryUnescape(c.Params("external"))
	internal, ok := links[external]
	if !ok {
		return c.Status(fiber.StatusNotFound).SendString("Link not found")
	}
	return c.JSON(GetLinkResponse{internal})
}
