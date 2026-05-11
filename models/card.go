package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	BoardID     uint       `gorm:"not null;index"`
	Column      string     `gorm:"not null"`
	Position    int        `gorm:"not null;default:0"`
	Title       string     `gorm:"not null;size:200"`
	Description string
	DueDate     *time.Time
	Priority    string
	Labels      []Label   `gorm:"many2many:card_labels;"`
	SubCards    []SubCard `gorm:"foreignKey:CardID"`
}
