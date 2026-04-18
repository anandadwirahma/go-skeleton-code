package usecase

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/yourusername/go-skeleton-code/internal/entity"
	"github.com/yourusername/go-skeleton-code/internal/model"
	"github.com/yourusername/go-skeleton-code/internal/model/converter"
	"github.com/yourusername/go-skeleton-code/internal/repository"
)

type ExampleUsecase interface {
	Create(ctx context.Context, request *model.CreateExampleRequest) (*model.ExampleResponse, error)
}

type exampleUsecase struct {
	db                *gorm.DB
	log               *logrus.Logger
	exampleRepository *repository.ExampleRepository
}

func NewExampleUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	exampleRepo *repository.ExampleRepository,
) ExampleUsecase {
	return &exampleUsecase{
		db:                db,
		log:               log,
		exampleRepository: exampleRepo,
	}
}

func (e *exampleUsecase) Create(ctx context.Context, request *model.CreateExampleRequest) (*model.ExampleResponse, error) {
	example := &entity.Example{
		Name:  request.Name,
		Email: request.Email,
	}

	err := e.exampleRepository.Create(e.db, example)
	if err != nil {
		e.log.WithError(err).Error("error creating example")
		return nil, fiber.ErrInternalServerError
	}

	return converter.ExampleToResponse(example), nil
}
