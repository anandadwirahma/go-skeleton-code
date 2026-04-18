package repository

import (
	"github.com/sirupsen/logrus"

	"github.com/yourusername/go-skeleton-code/internal/entity"
)

type ExampleRepository struct {
	Repository[entity.Example]
	log *logrus.Logger
}

func NewExampleRepository(
	log *logrus.Logger,
) *ExampleRepository {
	return &ExampleRepository{
		log: log,
	}
}
