package fiber

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

type HTTPServer struct {
	app     *fiber.App
	payload chan []byte
}

func NewHTTPServer() *HTTPServer {
	engine := html.New("./fiber", ".gohtml")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	res := &HTTPServer{
		app:     app,
		payload: make(chan []byte),
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("page_template", fiber.Map{
			"Link": "ws://localhost:7000/wsendpoint/1",
		})
	})

	app.Use("/wsendpoint", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/wsendpoint/:id", websocket.New(res.handlerWebSocket))

	return res
}

func (h *HTTPServer) Work() {
	defer h.cleanUp()

	log.Fatal(h.app.Listen(":7000"))
}

func (h *HTTPServer) Write(msg []byte) (int, error) {
	h.payload <- msg

	return 0, nil
}

func (h *HTTPServer) handlerWebSocket(c *websocket.Conn) {
	for {
		message := <-h.payload

		if err := c.WriteMessage(1, message); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (h *HTTPServer) cleanUp() {
	close(h.payload)
}
