package microservice

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"microservice/middleware/jwtAuth"
	"microservice/models"
	"strings"
)

type Server struct {
	app *fiber.App
	log *logrus.Logger
	db  *gorm.DB
}

func NewServer(dsn string, logger *logrus.Logger) (srv *Server, err error) {

	var db *gorm.DB
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err = buildSchema(db); err != nil {
		logger.Error(err)
		return nil, err
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// middleware configuration
	// any route that is not 'public' should be protected
	// refer to documentation in the middleware folder
	const publicPath = "/public/"

	app.Use(jwtAuth.New(jwtAuth.Config{
		Filter: func(c *fiber.Ctx) bool {
			path := strings.ToLower(c.Path())
			isPublic := strings.Contains(path, publicPath)
			return isPublic
		},
	}, logger))

	srv = &Server{
		app: app,
		log: logger,
		db:  db,
	}

	srv.buildRoutes()
	logger.Info("server initialised")
	return srv, nil
}

func buildSchema(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(models.Product{}); err != nil {
		return err
	}
	return nil
}

func (srv *Server) buildRoutes() {

	srv.app.Get("/shoes", srv.GetHandlerExample())
	srv.app.Get("/public/shoes", srv.GetHandlerExample())
	srv.app.Get("/public/token", srv.GenerateToken())
	srv.app.Get("/protected", srv.ProtectedRoute())
}

// Listen is a method on Server that starts the Fiber server
func (srv *Server) Listen(addr string) error {
	if err := srv.app.Listen(addr); err != nil {
		srv.log.Error(fmt.Sprintf("error at server.Listen: '%srv'", err.Error()))
		return err
	}
	return nil
}
