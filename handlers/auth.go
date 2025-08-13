package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"tma/auth"
	"tma/models"
)

type AuthHandler struct {
	db         *sql.DB
	jwtManager *auth.JWTManager
	botToken   string
}

func NewAuthHandler(db *sql.DB, jwtManager *auth.JWTManager, botToken string) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
		botToken:   botToken,
	}
}

// Auth handles Telegram Mini App authentication
func (h *AuthHandler) Auth(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
			"received": c.Request.Body,
		})
		return
	}

	// Log the received init data for debugging
	if req.InitData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "init_data is required"})
		return
	}

	// Validate Telegram init data
	telegramUser, err := auth.ValidateTelegramInitData(req.InitData, h.botToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram data: " + err.Error()})
		return
	}

	// Get or create user
	user, err := h.getOrCreateUser(telegramUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user: " + err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	c.JSON(http.StatusOK, response)
}

// TestAuth handles test authentication without hash validation
func (h *AuthHandler) TestAuth(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.InitData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "init_data is required"})
		return
	}

	// Parse init data manually for testing
	values, err := url.ParseQuery(req.InitData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse init_data: " + err.Error()})
		return
	}

	// Extract user data
	userIDStr := values.Get("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in init_data"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id: " + err.Error()})
		return
	}

	telegramUser := &models.TelegramUser{
		ID:        userID,
		Username:  values.Get("username"),
		FirstName: values.Get("first_name"),
		LastName:  values.Get("last_name"),
	}

	// Get or create user
	user, err := h.getOrCreateUser(telegramUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user: " + err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test authentication successful",
		"data": response,
		"parsed_values": values,
	})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) getOrCreateUser(telegramUser *models.TelegramUser) (*models.User, error) {
	// Try to get existing user
	user, err := h.getUserByTelegramID(telegramUser.ID)
	if err == nil {
		// User exists, update if needed
		if user.Username != telegramUser.Username || 
		   user.FirstName != telegramUser.FirstName || 
		   user.LastName != telegramUser.LastName {
			return h.updateUser(user.ID, telegramUser)
		}
		return user, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new user
	return h.createUser(telegramUser)
}

func (h *AuthHandler) getUserByTelegramID(telegramID int64) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at
		FROM users WHERE telegram_id = $1
	`
	
	user := &models.User{}
	err := h.db.QueryRow(query, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, 
		&user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (h *AuthHandler) getUserByID(userID int) (*models.User, error) {
	query := `
		SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	user := &models.User{}
	err := h.db.QueryRow(query, userID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, 
		&user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (h *AuthHandler) createUser(telegramUser *models.TelegramUser) (*models.User, error) {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, telegram_id, username, first_name, last_name, created_at, updated_at
	`
	
	user := &models.User{}
	err := h.db.QueryRow(query, 
		telegramUser.ID, telegramUser.Username, telegramUser.FirstName, telegramUser.LastName,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, 
		&user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (h *AuthHandler) updateUser(userID int, telegramUser *models.TelegramUser) (*models.User, error) {
	query := `
		UPDATE users 
		SET username = $1, first_name = $2, last_name = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, telegram_id, username, first_name, last_name, created_at, updated_at
	`
	
	user := &models.User{}
	err := h.db.QueryRow(query, 
		telegramUser.Username, telegramUser.FirstName, telegramUser.LastName, time.Now(), userID,
	).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, 
		&user.LastName, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}
