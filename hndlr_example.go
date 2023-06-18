package microservice

import (
	"github.com/gofiber/fiber/v2"
)

// GetHandlerExample is an example handler
// {get} http://example.com/shoes?order=desc&brand=nike
func (srv *Server) GetHandlerExample() fiber.Handler {

	srv.log.Info("report that route was called")

	return func(c *fiber.Ctx) (err error) {

		c.Query("order")         // "desc"
		c.Query("brand")         // "nike"
		c.Query("empty", "nike") // "nike"

		c.Status(fiber.StatusOK)
		_ = c.Send([]byte("OK"))
		return nil
	}

}
