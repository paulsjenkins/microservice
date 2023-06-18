package jwtAuth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Config is needed by all middleware
type Config struct {
	Filter       func(c *fiber.Ctx) bool                    // Required
	Unauthorized fiber.Handler                              // middleware specfic
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error) // middleware specfic
	Secret       string                                     // middleware specfic
	Expiry       int64                                      // middleware specfic
}

// ConfigDefault gives default values if not passed
var ConfigDefault = Config{

	// Filter: when returned true, our middleware is skipped
	Filter: nil,

	// Decode decodes the token
	Decode: nil,

	// Unauthorized is a function that is run when there is an error decoding the token
	Unauthorized: nil,

	// Secret is used to sign the token
	Secret: "a_very_weak_secret",

	// Expiry is the token expiry in seconds
	Expiry: 60,
}
