package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// TelegramUserData represents the user data structure from Telegram init_data
type TelegramUserData struct {
	ID                int64  `json:"id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	LanguageCode      string `json:"language_code"`
	AllowsWriteToPM   bool   `json:"allows_write_to_pm"`
	PhotoURL          string `json:"photo_url"`
}

// ValidateTelegramInitData validates Telegram Mini App init data
func ValidateTelegramInitData(initData, botToken string) (*models.TelegramUser, error) {
	if initData == "" {
		return nil, fmt.Errorf("init_data is empty")
	}

	// Log the received init data for debugging
	fmt.Printf("Received initData: %s\n", initData)

	// Parse the init data
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init_data: %w", err)
	}

	// Log parsed values for debugging
	fmt.Printf("Parsed values: %+v\n", values)

	// Extract user data
	userStr := values.Get("user")
	if userStr == "" {
		return nil, fmt.Errorf("user data not found in init_data")
	}

	// URL decode the user string
	decodedUserStr, err := url.QueryUnescape(userStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user data: %w", err)
	}

	fmt.Printf("Decoded user string: %s\n", decodedUserStr)

	// Parse user JSON
	var userData TelegramUserData
	if err := json.Unmarshal([]byte(decodedUserStr), &userData); err != nil {
		return nil, fmt.Errorf("failed to parse user JSON: %w", err)
	}

	fmt.Printf("Parsed user data: %+v\n", userData)

	// Validate hash if bot token is provided
	if botToken != "" {
		hash := values.Get("hash")
		if hash == "" {
			return nil, fmt.Errorf("hash not found in init_data")
		}

		fmt.Printf("Validating hash: %s\n", hash)
		if !validateHash(values, hash, botToken) {
			return nil, fmt.Errorf("invalid hash")
		}
		fmt.Printf("Hash validation successful\n")
	} else {
		fmt.Printf("Skipping hash validation (no bot token provided)\n")
	}

	user := &models.TelegramUser{
		ID:        userData.ID,
		Username:  userData.Username,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
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
