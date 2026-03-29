package domain

import "time"

type User struct {
	UserID      string // Firebase Auth UID
	Email       string
	DisplayName string
	CreatedAt   time.Time
	LastLoginAt time.Time
}
