package models

import "gorm.io/gorm"

type Label struct {
	gorm.Model
	BoardID uint   `gorm:"not null;index"`
	Name    string `gorm:"not null"`
	Color   string `gorm:"not null"`
}
