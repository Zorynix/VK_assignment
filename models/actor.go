package models

import "time"

type Actor struct {
	ID          int       `gorm:"primary_key"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Gender      string    `gorm:"type:char(1)"`
	DateOfBirth time.Time `gorm:"type:date"`
}
