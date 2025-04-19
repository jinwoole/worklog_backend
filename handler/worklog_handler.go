// handler/worklog_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinwoole/worklog-backend/service"
)

// WorkLogHandler handles work log HTTP requests
type WorkLogHandler struct {
	workLogService service.WorkLogService
}

// NewWorkLogHandler constructs a new WorkLogHandler
func NewWorkLogHandler(wls service.WorkLogService) *WorkLogHandler {
	return &WorkLogHandler{workLogService: wls}
}

// CreateWorkLog handles creation of today's work log
func (h *WorkLogHandler) CreateWorkLog(c *gin.Context) {
	userID := c.GetInt("userID")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logEntry, err := h.workLogService.CreateWorkLog(userID, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workLog": logEntry})
}

// UpdateWorkLog handles updating today's work log
func (h *WorkLogHandler) UpdateWorkLog(c *gin.Context) {
	userID := c.GetInt("userID")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.workLogService.UpdateWorkLog(userID, req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetAllWorkLogs retrieves all work logs for the authenticated user
func (h *WorkLogHandler) GetAllWorkLogs(c *gin.Context) {
	userID := c.GetInt("userID")
	logs, err := h.workLogService.GetAllWorkLogs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workLogs": logs})
}

// GetMe returns the authenticated user's ID (or full user data if desired)
func (h *WorkLogHandler) GetMe(c *gin.Context) {
	userID := c.GetInt("userID")
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
