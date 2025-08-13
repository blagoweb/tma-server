package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"tma/models"
)

type ItemsHandler struct {
	db *sql.DB
}

func NewItemsHandler(db *sql.DB) *ItemsHandler {
	return &ItemsHandler{db: db}
}

// GetItems returns all items for the authenticated user
func (h *ItemsHandler) GetItems(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM items 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	
	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Description,
			&item.Status, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan item"})
			return
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

// GetItem returns a specific item by ID
func (h *ItemsHandler) GetItem(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM items 
		WHERE id = $1 AND user_id = $2
	`
	
	var item models.Item
	err = h.db.QueryRow(query, itemID, userID).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CreateItem creates a new item
func (h *ItemsHandler) CreateItem(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	query := `
		INSERT INTO items (user_id, title, description, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, user_id, title, description, status, created_at, updated_at
	`
	
	var item models.Item
	err := h.db.QueryRow(query, userID, req.Title, req.Description).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem updates an existing item
func (h *ItemsHandler) UpdateItem(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Build dynamic query based on provided fields
	query := `
		UPDATE items 
		SET updated_at = $1
	`
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.Title != "" {
		query += `, title = $` + strconv.Itoa(argIndex)
		args = append(args, req.Title)
		argIndex++
	}

	if req.Description != "" {
		query += `, description = $` + strconv.Itoa(argIndex)
		args = append(args, req.Description)
		argIndex++
	}

	if req.Status != "" {
		query += `, status = $` + strconv.Itoa(argIndex)
		args = append(args, req.Status)
		argIndex++
	}

	query += ` WHERE id = $` + strconv.Itoa(argIndex) + ` AND user_id = $` + strconv.Itoa(argIndex+1)
	args = append(args, itemID, userID)

	query += ` RETURNING id, user_id, title, description, status, created_at, updated_at`

	var item models.Item
	err = h.db.QueryRow(query, args...).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem deletes an item
func (h *ItemsHandler) DeleteItem(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	query := `DELETE FROM items WHERE id = $1 AND user_id = $2`
	
	result, err := h.db.Exec(query, itemID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
