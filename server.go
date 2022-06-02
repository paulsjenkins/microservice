package microservice

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"psj/microservice/models"
)

type Server struct {
	app *fiber.App
	log *logrus.Logger
	db  *gorm.DB
}

func NewServer(dsn string, l *logrus.Logger) (srv *Server, err error) {

	var db *gorm.DB
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		l.Error(err)
		return nil, err
	}

	if err = buildSchema(db); err != nil {
		l.Error(err)
		return nil, err
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	srv = &Server{
		app: app,
		log: l,
		db:  db,
	}

	srv.buildRoutes()
	l.Info("server initialised: ok")

	return srv, nil
}

func buildSchema(db *gorm.DB) (err error) {
	if err = db.AutoMigrate(models.Product{}); err != nil {
		return err
	}
	return nil
}

func (srv *Server) buildRoutes() {
	// see also: example_handler.go
	srv.app.Get("/shoes", srv.GetHandlerExample())
}

// Listen is a method on Server that starts the Fiber server
func (srv *Server) Listen(addr string) error {
	srv.log.Info(fmt.Sprintf("server is running at %s", addr))
	if err := srv.app.Listen(addr); err != nil {
		srv.log.Error(fmt.Sprintf("error at server.Listen: '%srv'", err.Error()))
		return err
	}
	return nil
}
