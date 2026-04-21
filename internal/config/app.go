package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/yourusername/go-skeleton-code/internal/delivery/http"
	"github.com/yourusername/go-skeleton-code/internal/delivery/http/router/route"

	// httpGateway "github.com/yourusername/go-skeleton-code/internal/gateway/http"
	"github.com/yourusername/go-skeleton-code/internal/repository"
	"github.com/yourusername/go-skeleton-code/internal/usecase"
)

type AppConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func App(config *AppConfig) {
	// setup gateway httpClient
	// gtw := httpGateway.New(
	// 	httpGateway.WithBaseURL(config.Config.GetString("api.example.base_url")),
	// 	httpGateway.WithTimeout(config.Config.GetDuration("api.example.duration")),
	// )

	// setup repositories
	exampleRepository := repository.NewExampleRepository(config.Log)

	// setup use cases
	exampleUseCase := usecase.NewExampleUsecase(config.DB, config.Log, exampleRepository)

	// setup controller
	exampleController := http.NewExampleController(exampleUseCase, config.Validate, config.Log)

	routeConfig := route.RouteConfig{
		App:               config.App,
		ExampleController: exampleController,
	}
	routeConfig.Setup()
}
