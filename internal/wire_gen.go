// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	logger2 "github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"go-clean-architecture-example/config"
	"go-clean-architecture-example/docs"
	"go-clean-architecture-example/internal/api"
	"go-clean-architecture-example/internal/app"
	"go-clean-architecture-example/internal/commom/exception"
	"go-clean-architecture-example/internal/infrastructure/notification"
	"go-clean-architecture-example/internal/infrastructure/persistence"
	"go-clean-architecture-example/internal/probes"
	"go-clean-architecture-example/internal/router"
	"go-clean-architecture-example/pkg/logger"
	"os"
	"time"
)

// Injectors from server.go:

func New() (*Server, error) {
	configuration, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	repository := persistence.NewCragMemRepository()
	service := notification.NewNotificationService()
	application := app.NewApplication(repository, service)
	cragHttpApi := api.NewCragHttpApi(application)
	cragRouter := router.NewCragRouter(cragHttpApi)
	healthCheckApplication := probes.NewHealthChecker(configuration)
	server := NewServer(configuration, cragRouter, healthCheckApplication)
	return server, nil
}

// server.go:

// Server struct
type Server struct {
	app    *fiber.App
	cfg    *config.Configuration
	logger logger.Logger
}

// @title  My SERVER
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email minkj1992@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000
// @BasePath /
func NewServer(
	cfg *config.Configuration,
	cragRouter router.CragRouter,
	healthCheckApp probes.HealthCheckApplication) *Server {
	logger3 := logger.NewApiLogger(cfg)
	app2 := fiber.New(fiber.Config{
		ErrorHandler: exception.CustomErrorHandler,
		ReadTimeout:  time.Second * cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * cfg.Server.WriteTimeout,
	})
	app2.
		Use(logger2.New(logger2.Config{
			Next:         nil,
			Done:         nil,
			Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stdout,
		}))
	app2.
		Use(cors.New())
	app2.
		Use(etag.New())
	app2.
		Use(recover2.New())

	setSwagger(cfg.Server.BaseURI)
	app2.
		Get("/swagger/*", swagger.HandlerDefault)
	app2.
		Get("/liveness", func(c *fiber.Ctx) error {
			result := healthCheckApp.LiveEndpoint()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})
	app2.
		Get("/readiness", func(c *fiber.Ctx) error {
			result := healthCheckApp.ReadyEndpoint()
			if result.Status {
				return c.Status(fiber.StatusOK).JSON(result)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(result)
		})
	api2 := app2.Group("/api")
	v1 := api2.Group("/v1")
	cragRouter.Init(&v1)

	return &Server{
		cfg:    cfg,
		logger: logger3,
		app:    app2,
	}
}

func (serv Server) App() *fiber.App {
	return serv.app
}

func (serv Server) Config() *config.Configuration {
	return serv.cfg
}

func (serv Server) Logger() logger.Logger {
	return serv.logger
}

func setSwagger(baseURI string) {
	docs.SwaggerInfo.
		Title = "Go Clean Architecture Example ✈️"
	docs.SwaggerInfo.
		Description = "This is a go clean architecture example."
	docs.SwaggerInfo.
		Version = "1.0"
	docs.SwaggerInfo.
		Host = baseURI
	docs.SwaggerInfo.
		BasePath = "/api/v1"
}
