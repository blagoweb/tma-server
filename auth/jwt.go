package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"tma/models"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey: secretKey}
}

type Claims struct {
	UserID     int   `json:"user_id"`
	TelegramID int64 `json:"telegram_id"`
	jwt.RegisteredClaims
}

func (j *JWTManager) GenerateToken(user models.User) (string, error) {
	claims := &Claims{
		UserID:     user.ID,
		TelegramID: user.TelegramID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateTelegramInitData validates Telegram Mini App init data
func ValidateTelegramInitData(initData, botToken string) (*models.TelegramUser, error) {
	if initData == "" {
		return nil, fmt.Errorf("init_data is empty")
	}

	// Parse the init data
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init_data: %w", err)
	}

	// Extract user data
	userStr := values.Get("user")
	if userStr == "" {
		return nil, fmt.Errorf("user data not found in init_data")
	}

	// Parse user JSON (simplified - in production you'd use proper JSON parsing)
	// For now, we'll extract basic info from the query string
	userIDStr := values.Get("user_id")
	if userIDStr == "" {
		return nil, fmt.Errorf("user_id not found in init_data")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// Validate hash if bot token is provided
	if botToken != "" {
		hash := values.Get("hash")
		if hash == "" {
			return nil, fmt.Errorf("hash not found in init_data")
		}

		if !validateHash(values, hash, botToken) {
			return nil, fmt.Errorf("invalid hash")
		}
	}

	user := &models.TelegramUser{
		ID:        userID,
		Username:  values.Get("username"),
		FirstName: values.Get("first_name"),
		LastName:  values.Get("last_name"),
	}

	return user, nil
}

func validateHash(values url.Values, hash, botToken string) bool {
	// Remove hash from values
	dataCheckString := make([]string, 0)
	for key, values := range values {
		if key == "hash" {
			continue
		}
		for _, value := range values {
			dataCheckString = append(dataCheckString, key+"="+value)
		}
	}

	// Sort alphabetically
	sort.Strings(dataCheckString)

	// Create data check string
	dataCheckStr := strings.Join(dataCheckString, "\n")

	// Create secret key
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))

	// Calculate hash
	h := hmac.New(sha256.New, secretKey.Sum(nil))
	h.Write([]byte(dataCheckStr))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	return calculatedHash == hash
}
