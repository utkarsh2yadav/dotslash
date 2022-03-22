package main

import (
	"dotslash/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:lang", websocket.New(func(c *websocket.Conn) {
		lang := c.Params("lang")
		l := service.GetLanguage(lang)
		if l == nil {
			c.WriteMessage(websocket.TextMessage, []byte("Request language is not yet supported"))
			return
		}
		l.Handler(c)
	}))

	log.Fatal(app.Listen(":3000"))
}
