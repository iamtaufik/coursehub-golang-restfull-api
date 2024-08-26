package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    	string  	`gorm:"not null;unique"`
	Password 	string  	`gorm:"not null"`
	Profile  	Profile		`gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Courses 	[]Course	`gorm:"many2many:user_courses;"` 	
}