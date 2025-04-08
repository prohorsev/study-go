package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var postLikes = map[string]int64{}

func main() {
	webApp := fiber.New(fiber.Config{
		Immutable:      true,
		ReadBufferSize: 16 * 1024})
	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Go to /likes/12345")
	})
	webApp.Get("/likes/:post_id", func(c *fiber.Ctx) error {
		postId := c.Params("post_id", "")
		if likes, ok := postLikes[postId]; ok {
			return c.SendString(strconv.Itoa(int(likes)))
		}
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})
	webApp.Post("/likes/:post_id", func(c *fiber.Ctx) error {
		postId := c.Params("post_id", "")
		if _, ok := postLikes[postId]; ok {
			postLikes[postId] = postLikes[postId] + 1
			return c.SendString(strconv.Itoa(int(postLikes[postId])))
		}
		postLikes[postId] = 1
		return c.Status(fiber.StatusCreated).SendString(strconv.Itoa(int(postLikes[postId])))
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
