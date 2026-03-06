// Package handler provides Gin HTTP handlers for the Contact resource.
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yourusername/go-skeleton-code/internal/delivery/http/dto"
	"github.com/yourusername/go-skeleton-code/internal/domain/contact"
)

// ContactHandler holds the usecase dependency and logger for contact endpoints.
type ContactHandler struct {
	usecase contact.Usecase
	log     *zap.Logger
}

// NewContactHandler creates a new ContactHandler.
func NewContactHandler(uc contact.Usecase, log *zap.Logger) *ContactHandler {
	return &ContactHandler{usecase: uc, log: log}
}

// Create handles POST /contacts
//
//	@Summary      Create a new contact
//	@Description  Submit a new contact form entry
//	@Tags         contacts
//	@Accept       json
//	@Produce      json
//	@Param        body  body      dto.CreateContactRequest  true  "Contact payload"
//	@Success      201   {object}  dto.SuccessResponse
//	@Failure      400   {object}  dto.ErrorResponse
//	@Failure      422   {object}  dto.ErrorResponse
//	@Failure      500   {object}  dto.ErrorResponse
//	@Router       /contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	result, err := h.usecase.Create(c.Request.Context(), contact.CreateInput{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Data:    dto.ToResponse(result),
	})
}

// GetAll handles GET /contacts
func (h *ContactHandler) GetAll(c *gin.Context) {
	contacts, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    dto.ToResponseList(contacts),
	})
}

// GetByID handles GET /contacts/:id
func (h *ContactHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: "invalid id"})
		return
	}

	result, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    dto.ToResponse(result),
	})
}

// Update handles PUT /contacts/:id
func (h *ContactHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: "invalid id"})
		return
	}

	var req dto.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	result, err := h.usecase.Update(c.Request.Context(), id, contact.UpdateInput{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    dto.ToResponse(result),
	})
}

// Delete handles DELETE /contacts/:id
func (h *ContactHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: "invalid id"})
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    "contact deleted successfully",
	})
}

// --- helpers ---

// parseID extracts and validates the :id route parameter.
func parseID(c *gin.Context) (uint, error) {
	raw := c.Param("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// handleError maps domain errors to appropriate HTTP status codes.
func (h *ContactHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, contact.ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: "contact not found"})
	case errors.Is(err, contact.ErrEmailTaken):
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Success: false, Error: "email already registered"})
	default:
		h.log.Error("unhandled error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: "internal server error"})
	}
}
