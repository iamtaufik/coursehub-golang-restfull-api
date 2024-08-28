package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	OrderID 	string
	UserID 		uint
	CourseID 	uint
	Amount 		float64
	Status 		string
	PaymentType string
}