package model

import "gorm.io/gorm"

type Module struct {
	gorm.Model
	ChapterID 	uint
	Title    	string
	Duration 	int
	VideoURL 	string
	IsTrailer 	bool	`gorm:"default:false"`
}