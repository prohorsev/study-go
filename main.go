package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

type (
	CreateItemRequest struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

	Item struct {
		Name  string `json:"name"`
		Price uint   `json:"price"`
	}

	ItemsHandler struct{}
)

var (
	items []Item
)

func main() {
	viewsEngine := html.New("./templates", ".tmpl")
	webApp := fiber.New(fiber.Config{
		Views:          viewsEngine,
		ReadBufferSize: 16 * 1024,
	})
	itemsHandler := ItemsHandler{}

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	webApp.Post("/items", itemsHandler.Create)
	webApp.Get("/items/view", itemsHandler.View)

	logrus.Fatal(webApp.Listen(":8080"))
}

func (h *ItemsHandler) Create(c *fiber.Ctx) error {
	req := CreateItemRequest{}
	if err := c.BodyParser(&req); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	items = append(items, Item{
		Name:  req.Name,
		Price: req.Price,
	})
	logrus.Info(items)
	return c.SendString("OK")
}

func (h *ItemsHandler) View(c *fiber.Ctx) error {
	return c.Render("items", items)
}
