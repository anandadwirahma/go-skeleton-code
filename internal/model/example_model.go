package model

type CreateExampleRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"max=200,email"`
}

type ExampleResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
