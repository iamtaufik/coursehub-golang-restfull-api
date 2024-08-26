package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	UserID    	uint   	
	FirstName 	string 	`gorm:"size:256;not null"`
	LastName  	string 	`gorm:"size:256;not null"`
	PhoneNumber *string 
	Address    	*string	
}