package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jamolpe/kubevisual-agent/internal/config"
	"github.com/jamolpe/kubevisual-agent/internal/nodeworker"
	"github.com/jamolpe/kubevisual-agent/internal/podworker"
)

type API struct {
	app        *fiber.App
	podWorker  *podworker.PodWorker
	nodeWorker *nodeworker.NodeWorker
}

func New() *API {
	app := fiber.New()
	cmdConfig := config.GetcmdConfig()
	kubeclient := config.GetKubernetesClient(cmdConfig)
	metricsClient := config.GetKubernetesMetrics(cmdConfig)
	return &API{app: app, podWorker: podworker.New(kubeclient, metricsClient), nodeWorker: nodeworker.New(kubeclient, metricsClient)}
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
