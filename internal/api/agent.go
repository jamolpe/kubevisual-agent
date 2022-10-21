package api

import "github.com/gofiber/fiber/v2"

func (api *API) AgentRoutes(router fiber.Router) {
	api.NodeRoutes(router.Group("/node"))
	api.PodRoutes(router.Group("/pod"))
}

func (api *API) NodeRoutes(router fiber.Router) {
	router.Get("", func(c *fiber.Ctx) error {
		return c.Status(200).JSON("ok")
	})
}

func (api *API) PodRoutes(router fiber.Router) {
	router.Get("", func(c *fiber.Ctx) error {
		return c.Status(200).JSON("ok")
	})
}
