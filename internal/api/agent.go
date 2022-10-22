package api

import (
	"github.com/gofiber/fiber/v2"
)

func (api *API) AgentRoutes(router fiber.Router) {
	api.PodRoutes(router)
}

func (api *API) PodRoutes(router fiber.Router) {
	router.Get("/pod", func(c *fiber.Ctx) error {
		pods, err := api.podDescriber.GetAllPodsInformation()
		if err != nil {
			return c.Status(500).JSON(&Error{message: "error geting pod info"})
		}
		return c.Status(200).JSON(pods)
	})
}
