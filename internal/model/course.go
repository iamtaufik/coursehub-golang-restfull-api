package model

import (
	"gorm.io/gorm"
)

type Levels string

const (
	Beginner     Levels = "beginner"
	Intermediate Levels = "intermediate"
	Advanced     Levels = "advanced"
)

type Course struct {
	gorm.Model
	Title 			string
	Description 	string
	Requirements 	string  	
	Levels 			Levels 		`gorm:"type:levels;default:'beginner'"`
	Chapters 		[]Chapter  	`gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
	CategoryID 		uint
	Category 		Category
	ImageURL 		string
	TelegramLink 	*string
	Price 			float64
	Author 			string
	Users 			[]User 	`gorm:"many2many:user_courses;"`
}

type CourseResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Author      string `json:"author"`
	CategoryID  int    `json:"-"`
	Category    struct {
		ID           int    `json:"id"`
		CategoryName string `json:"category_name"`
	} `json:"category"`
}

