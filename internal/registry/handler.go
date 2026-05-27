package registry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Handler handles HTTP requests for the registry
type Handler struct {
	repo   Repository
	logger *zap.Logger
}

// NewHandler creates a new registry handler
func NewHandler(repo Repository, logger *zap.Logger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

// RegisterRoutes registers the registry routes to the Gin engine
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/services", h.CreateService)
		api.GET("/services", h.ListServices)
		api.GET("/services/:id", h.GetService)
		
		api.POST("/services/:id/contracts", h.CreateContract)
		api.GET("/contracts", h.ListContracts)
		api.GET("/incidents", h.ListIncidents)
	}
}

// CreateService creates a new microservice registration
func (h *Handler) CreateService(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Owner       string `json:"owner" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := &Service{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Owner:       req.Owner,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.CreateService(c.Request.Context(), service); err != nil {
		h.logger.Error("Failed to create service", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// ListServices lists all registered services
func (h *Handler) ListServices(c *gin.Context) {
	services, err := h.repo.ListServices(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list services", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetService gets a single service by ID
func (h *Handler) GetService(c *gin.Context) {
	id := c.Param("id")
	service, err := h.repo.GetService(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get service", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// CreateContract creates a new SLA contract for a service
func (h *Handler) CreateContract(c *gin.Context) {
	serviceID := c.Param("id")

	var req struct {
		Policy string `json:"policy" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract := &SLAContract{
		ID:        uuid.New().String(),
		ServiceID: serviceID,
		Policy:    req.Policy,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.CreateSLAContract(c.Request.Context(), contract); err != nil {
		h.logger.Error("Failed to create contract", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SLA contract"})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// ListContracts lists all SLA contracts globally
func (h *Handler) ListContracts(c *gin.Context) {
	contracts, err := h.repo.ListSLAContracts(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list contracts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list contracts"})
		return
	}
	c.JSON(http.StatusOK, contracts)
}

// ListIncidents lists all recorded incidents
func (h *Handler) ListIncidents(c *gin.Context) {
	incidents, err := h.repo.ListIncidents(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to list incidents", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list incidents"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}
