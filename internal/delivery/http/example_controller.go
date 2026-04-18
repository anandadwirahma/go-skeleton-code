package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"

	"github.com/yourusername/go-skeleton-code/internal/model"
	"github.com/yourusername/go-skeleton-code/internal/usecase"
)

type ExampleController struct {
	useCase  usecase.ExampleUsecase
	validate *validator.Validate
	log      *logrus.Logger
}

func NewExampleController(
	useCase usecase.ExampleUsecase,
	validator *validator.Validate,
	log *logrus.Logger,
) *ExampleController {
	return &ExampleController{
		log:      log,
		validate: validator,
		useCase:  useCase,
	}
}

func (c *ExampleController) Create(ctx fiber.Ctx) error {
	request := new(model.CreateExampleRequest)

	err := ctx.Bind().Body(request)
	if err != nil {
		c.log.WithError(err).Error("failed to bind request body")
		return fiber.ErrBadRequest
	}

	err = c.validate.Struct(request)
	if err != nil {
		c.log.WithError(err).Errorf("invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.useCase.Create(ctx.Context(), request)
	if err != nil {
		c.log.WithError(err).Error("failed to create examples")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ExampleResponse]{Data: response})
}
