package domain

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepository định nghĩa những gì tầng usecase cần từ DB
type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	GetByID(id uint) (*User, error)
}
