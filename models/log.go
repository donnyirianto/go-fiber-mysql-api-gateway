package models

import "time"

type Log struct {
	ID        uint   `gorm:"primaryKey"`
	Level     string `gorm:"not null"`
	Message   string `gorm:"not null"`
	CreatedAt time.Time
}
