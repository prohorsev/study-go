package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var exchangeRate = map[string]float64{
	"USD/EUR": 0.8,
	"EUR/USD": 1.25,
	"USD/GBP": 0.7,
	"GBP/USD": 1.43,
	"USD/JPY": 110,
	"JPY/USD": 0.0091,
}

func main() {
	webApp := fiber.New()
	webApp.Get("/convert", func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		key := from + "/" + to
		if rate, ok := exchangeRate[key]; ok {
			return c.SendString(fmt.Sprintf("%.2f", rate))
		}

		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	logrus.Fatal(webApp.Listen(":8080"))
}
