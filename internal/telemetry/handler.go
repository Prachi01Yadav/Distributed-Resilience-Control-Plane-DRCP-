package telemetry

import (
	"net/http"
	"time"

	"github.com/arche/sentinelmesh/pkg/kafka"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TelemetryData struct {
	ServiceID  string  `json:"service_id"`
	StatusCode int     `json:"status_code"`
	LatencyMs  float64 `json:"latency_ms"`
	Timestamp  int64   `json:"timestamp"`
}

type Handler struct {
	producer *kafka.Producer
	logger   *zap.Logger
}

func NewHandler(producer *kafka.Producer, logger *zap.Logger) *Handler {
	return &Handler{
		producer: producer,
		logger:   logger,
	}
}

func (h *Handler) ReceiveTelemetry(c *gin.Context) {
	var req TelemetryData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Timestamp = time.Now().UnixMilli()

	if err := h.producer.ProduceTelemetry(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process telemetry"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}
