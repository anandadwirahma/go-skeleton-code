package converter

import (
	"time"

	"github.com/yourusername/go-skeleton-code/internal/entity"
	"github.com/yourusername/go-skeleton-code/internal/model"
)

func ExampleToResponse(example *entity.Example) *model.ExampleResponse {
	return &model.ExampleResponse{
		ID:        example.ID,
		Name:      example.Name,
		Email:     example.Email,
		CreatedAt: example.CreatedAt.Format(time.RFC3339),
		UpdatedAt: example.UpdatedAt.Format(time.RFC3339),
	}
}
