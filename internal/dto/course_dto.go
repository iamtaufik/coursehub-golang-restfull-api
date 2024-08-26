package dto

import "github.com/iamtaufik/coursehub-golang-restfull-api/internal/model"

type CreateCourseDTO struct {
	Title        string             `json:"title" binding:"required"`
	Description  string             `json:"description" binding:"required"`
	Requirements string             `json:"requirements" binding:"required"`
	Levels       model.Levels             `json:"levels" binding:"required"`
	ImageURL     string             `json:"imageURL" binding:"required"`
	CategoryID   uint               `json:"categoryID" binding:"required"`
	Price        float64            `json:"price" binding:"required"`
	Author       string             `json:"author" binding:"required"`
	Chapters     []CreateChapterDTO `json:"chapters" binding:"required"`
}

type CreateChapterDTO struct {
	Name    string            `json:"name" binding:"required"`
	Modules []CreateModuleDTO `json:"modules" binding:"required"`
}

type CreateModuleDTO struct {
	Title     string `json:"title" binding:"required"`
	Duration  int    `json:"duration" binding:"required"`
	VideoURL  string `json:"videoURL" binding:"required"`
	IsTrailer bool   `json:"isTrailer" binding:"required"`
}