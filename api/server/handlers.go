package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/joel-thompson/my-go-service/constants"
	"github.com/joel-thompson/my-go-service/storage"
)

// handleHealth returns a simple health check response
func (a *API) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": constants.StatusHealthy,
	})
}

// handleHello returns a simple hello world response
func (a *API) handleHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": constants.MessageHello,
	})
}

// handleCreateItem creates a new item
func (a *API) handleCreateItem(c *gin.Context) {
	var req storage.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	item, err := a.store.CreateItem(c.Request.Context(), req)
	if err != nil {
		a.logger.Error("Failed to create item", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create item",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// handleListItems retrieves a paginated list of items
func (a *API) handleListItems(c *gin.Context) {
	var req storage.ListItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		a.logger.Error("Failed to bind query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	response, err := a.store.ListItems(c.Request.Context(), req)
	if err != nil {
		a.logger.Error("Failed to list items", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve items",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
