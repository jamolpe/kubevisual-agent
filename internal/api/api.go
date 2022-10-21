package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jamolpe/kubevisual-agent/internal/podDescriber"
)

type API struct {
	app          *fiber.App
	podDescriber *podDescriber.PodDescriber
}

func New() *API {
	app := fiber.New()
	return &API{app: app, podDescriber: podDescriber.New()}
}

func (api *API) Configure() {
	api.app.Use(recover.New())
	api.app.Use(cors.New())
	api.DefineRoutes()
}

func (api *API) DefineRoutes() {
	api.AgentRoutes(api.app.Group("/kubevisual-agent", func(c *fiber.Ctx) error {
		return c.Next()
	}))
}

func (api *API) Listen(port string) {
	if api.app != nil {
		api.app.Listen(":" + port)
	}
}
