package main

import "time"

// User definition
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
}

// Session definition
type Session struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UUID      string     `json:"UUID"`
	UserID    uint       `json:"userID"`
}

// UserSetting definition
type UserSetting struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
	UserID     uint       `json:"userID"`
	LastOption string     `json:"lastOption"`
}
