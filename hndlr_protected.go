package microservice

import "github.com/gofiber/fiber/v2"

func (srv *Server) ProtectedRoute() fiber.Handler {

	srv.log.Info("report that route was called")
	return func(c *fiber.Ctx) (err error) {
		claimData := c.Locals("jwtClaims")
		if claimData == nil {
			return c.SendString("Jwt was bypassed")
		}
		return c.JSON(claimData)
	}
}
