package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"tma/models"

	"github.com/gin-gonic/gin"
)

type PagesHandler struct {
	db *sql.DB
}

func NewPagesHandler(db *sql.DB) *PagesHandler {
	return &PagesHandler{db: db}
}

// GetPages returns all pages for the authenticated user
func (h *PagesHandler) GetPages(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	query := `
		SELECT id, user_id, title, json_data, created_at, updated_at
		FROM pages
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pages"})
		return
	}
	defer rows.Close()

	pages := make([]models.Page, 0)
	for rows.Next() {
		var page models.Page
		err := rows.Scan(
			&page.ID, &page.UserID, &page.Title, &page.JSONData,
			&page.CreatedAt, &page.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan page"})
			return
		}
		pages = append(pages, page)
	}

	c.JSON(http.StatusOK, pages)
}

// GetPage returns a specific page by ID
func (h *PagesHandler) GetPage(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

		query := `
		SELECT id, user_id, title, json_data, created_at, updated_at
		FROM pages
		WHERE id = $1 AND user_id = $2
	`
	
	var page models.Page
	err = h.db.QueryRow(query, pageID, userID).Scan(
		&page.ID, &page.UserID, &page.Title, &page.JSONData,
		&page.CreatedAt, &page.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// CreatePage creates a new page
func (h *PagesHandler) CreatePage(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

		query := `
		INSERT INTO pages (user_id, title, json_data)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, title, json_data, created_at, updated_at
	`
	
	var page models.Page
	err := h.db.QueryRow(query, userID, req.Title, req.JSONData).Scan(
		&page.ID, &page.UserID, &page.Title, &page.JSONData,
		&page.CreatedAt, &page.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// UpdatePage updates an existing page
func (h *PagesHandler) UpdatePage(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	var req models.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Build dynamic query based on provided fields
	query := `
		UPDATE pages
		SET updated_at = $1
	`
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.Title != "" {
		query += `, title = $` + strconv.Itoa(argIndex)
		args = append(args, req.Title)
		argIndex++
	}

	if req.JSONData != nil {
		query += `, json_data = $` + strconv.Itoa(argIndex)
		args = append(args, req.JSONData)
		argIndex++
	}

	query += ` WHERE id = $` + strconv.Itoa(argIndex) + ` AND user_id = $` + strconv.Itoa(argIndex+1)
	args = append(args, pageID, userID)

	query += ` RETURNING id, user_id, title, json_data, created_at, updated_at`

	var page models.Page
	err = h.db.QueryRow(query, args...).Scan(
		&page.ID, &page.UserID, &page.Title, &page.JSONData,
		&page.CreatedAt, &page.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update page"})
		return
	}

	c.JSON(http.StatusOK, page)
}

// DeletePage deletes an page
func (h *PagesHandler) DeletePage(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page ID"})
		return
	}

	query := `DELETE FROM pages WHERE id = $1 AND user_id = $2`

	result, err := h.db.Exec(query, pageID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}
