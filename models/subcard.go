package models

import "gorm.io/gorm"

type SubCard struct {
	gorm.Model
	CardID    uint   `gorm:"not null;index"`
	Title     string `gorm:"not null;size:200"`
	Completed bool   `gorm:"not null;default:false"`
	Position  int    `gorm:"not null;default:0"`
}
