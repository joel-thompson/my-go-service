package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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

// handleGetItem retrieves a single item by ID
func (a *API) handleGetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		a.logger.Error("Invalid item ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID format",
		})
		return
	}

	item, err := a.store.GetItem(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found",
			})
			return
		}
		a.logger.Error("Failed to get item", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// handleUpdateItem updates an existing item
func (a *API) handleUpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		a.logger.Error("Invalid item ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID format",
		})
		return
	}

	var req storage.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		a.logger.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Ensure at least one field is provided
	if req.Name == nil && req.Description == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field (name or description) must be provided",
		})
		return
	}

	item, err := a.store.UpdateItem(c.Request.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found",
			})
			return
		}
		a.logger.Error("Failed to update item", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

// handleDeleteItem deletes an item by ID
func (a *API) handleDeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		a.logger.Error("Invalid item ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID format",
		})
		return
	}

	item, err := a.store.DeleteItem(c.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Item not found",
			})
			return
		}
		a.logger.Error("Failed to delete item", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
		"item":    item,
	})
}
