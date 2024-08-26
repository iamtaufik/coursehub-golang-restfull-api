package model

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	Name    	string
	CourseID 	uint
	Modules 	[]Module `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}