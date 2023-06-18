package microservice

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"microservice/middleware/jwtAuth"
)

// GenerateToken generates a new jwt token
func (srv *Server) GenerateToken() fiber.Handler {

	srv.log.Info("generating a new token")
	return func(c *fiber.Ctx) (err error) {
		token, err := jwtAuth.Encode(&jwt.MapClaims{
			"email": "test@gmail.com",
			"role":  "admin",
			"id":    7,
		}, 1000)

		if err != nil {
			return c.SendStatus(500)
		}
		return c.SendString(token)
	}
}
