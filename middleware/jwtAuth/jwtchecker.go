package jwtAuth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func New(config Config, logger *logrus.Logger) fiber.Handler {

	// For setting default config
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		requestPath := strings.ToLower(c.Path())
		// Don't execute middleware if Filter returns true
		if cfg.Filter != nil && cfg.Filter(c) {
			logger.Debugf("middleware skipped: '%s'", requestPath)
			return c.Next()
		}
		logger.Debugf("    middleware run: '%s'", requestPath)

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("jwtClaims", *claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}

func Encode(claims *jwt.MapClaims, expiryAfter int64) (string, error) {

	// setting default expiryAfter
	if expiryAfter == 0 {
		expiryAfter = ConfigDefault.Expiry
	}

	// or you can use time.Now().Add(time.Second * time.Duration(expiryAfter)).UTC().Unix()
	(*claims)["exp"] = time.Now().UTC().Unix() + expiryAfter

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// our signed jwt token string
	signedToken, err := token.SignedString([]byte(ConfigDefault.Secret))

	if err != nil {
		err := errors.New("error creating a token")
		return "", err
	}

	return signedToken, nil
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values if not passed
	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}

	// Set default secret if not passed
	if cfg.Secret == "" {
		cfg.Secret = ConfigDefault.Secret
	}

	// Set default expiry if not passed
	if cfg.Expiry == 0 {
		cfg.Expiry = ConfigDefault.Expiry
	}

	// this is the main jwt decode function of our middleware
	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {

			authHeader := c.Get("Authorization")

			if authHeader == "" {
				return nil, errors.New("Authorization header is required")
			}

			// we parse our jwt token and check for validity against our secret
			token, err := jwt.Parse(
				authHeader[7:],
				func(token *jwt.Token) (interface{}, error) {
					// verifying our algo
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(cfg.Secret), nil
				},
			)

			if err != nil {
				return nil, errors.New("Error parsing token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !(ok && token.Valid) {
				return nil, errors.New("Invalid token")
			}

			if expiresAt, ok := claims["exp"]; ok && int64(expiresAt.(float64)) < time.Now().UTC().Unix() {
				return nil, errors.New("jwt is expired")
			}

			return &claims, nil
		}
	}

	// Set default Unauthorized if not passed
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}
