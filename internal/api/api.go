package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type API struct {
	app *fiber.App
}

func New() *API {
	app := fiber.New()
	return &API{app: app}
}

func (api *API) Configure() {
	api.app.Use(recover.New())
	api.app.Use(cors.New())
	api.DefineRoutes()
}

func (api *API) DefineRoutes() {
	api.AgentRoutes(api.app.Group("/kubevisual-agent", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("V1")
	}))
}

func (api *API) Listen(port string) {
	if api.app != nil {
		api.app.Listen(":" + port)
	}
}
