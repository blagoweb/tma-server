package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSONData - кастомный тип для работы с JSON в PostgreSQL
type JSONData json.RawMessage

// Value - реализация интерфейса driver.Valuer для записи в БД
func (j JSONData) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}

	// Логируем, что отправляется в БД
	value := string(j)
	if len(value) > 100 {
		value = value[:100] + "..."
	}
	fmt.Printf("JSONData.Value(): Sending to DB: %s\n", value)

	return string(j), nil
}

// Scan - реализация интерфейса sql.Scanner для чтения из БД
func (j *JSONData) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSONData(v)
	case string:
		*j = JSONData(v)
	default:
		return json.Unmarshal([]byte(v.(string)), j)
	}
	return nil
}

// MarshalJSON - для сериализации в JSON
func (j JSONData) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON - для десериализации из JSON
func (j *JSONData) UnmarshalJSON(data []byte) error {
	// Проверяем, что это валидный JSON
	if len(data) == 0 {
		*j = nil
		return nil
	}

	// Проверяем, что это не строка "null"
	if string(data) == "null" {
		*j = nil
		return nil
	}

	// Валидируем JSON перед сохранением
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*j = JSONData(data)
	return nil
}

type User struct {
	ID         int       `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   string    `json:"username" db:"username"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthRequest struct {
	InitData string `json:"init_data"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Page struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	JSONData  JSONData  `json:"json_data" db:"json_data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePageRequest struct {
	Title    string   `json:"title" binding:"required"`
	JSONData JSONData `json:"json_data"`
}

type UpdatePageRequest struct {
	Title    string   `json:"title"`
	JSONData JSONData `json:"json_data"`
}
