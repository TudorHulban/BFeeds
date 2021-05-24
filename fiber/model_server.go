package fiber

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

type HTTPServer struct {
	app *fiber.App
}

func NewHTTPServer() *HTTPServer {
	engine := html.New("./fiber", ".gohtml")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("page_template", fiber.Map{
			"Link": "ws://" + c.Hostname() + "/ws",
		})
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(handlerWebSocket))

	return &HTTPServer{
		app: app,
	}
}

func (h *HTTPServer) Work() {
	log.Fatal(h.app.Listen(":7000"))
}

func handlerWebSocket(c *websocket.Conn) {
	log.Println(c.Locals("allowed"))  // true
	log.Println(c.Params("id"))       // 123
	log.Println(c.Query("v"))         // 1.0
	log.Println(c.Cookies("session")) // ""

	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	var (
		messageType int
		msg         []byte
		err         error
	)

	for {
		if messageType, msg, err = c.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", msg)

		if err = c.WriteMessage(messageType, msg); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
