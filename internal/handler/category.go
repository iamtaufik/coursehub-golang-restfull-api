package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/database"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/model"
)

func GetAllCategories(c *fiber.Ctx) error {
	db := database.StartDB()
	
	var categories []model.Category

	if err := db.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	type Result struct {
		ID uint `json:"id"`
		Name string `json:"name"`
	}

	var results []Result

	for _, category := range categories {
		results = append(results, Result{
			ID: category.ID,
			Name: category.CategoryName,
		})
	}

	return c.JSON(fiber.Map{
		"data": results,
	})
}

func GetCoursesByCategoryID(c *fiber.Ctx) error {
	type CourseResponse struct {
		ID          uint    `json:"id"`
		Title       string `json:"title"`
		Requirements string `json:"requirements"`
		Levels      model.Levels `json:"levels"`
		ImageURL    string `json:"imageURL"`
		Category    struct {
			ID           uint    `json:"id"`
			CategoryName string `json:"category_name"`
		} `json:"category"`
		Description string `json:"description"`
		Price       float64    `json:"price"`
		Author      string `json:"author"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}

	db := database.StartDB()

	categoryID := c.Params("id")

	var courses []model.Course

	if err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&courses).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var results []CourseResponse

	for _, course := range courses {
		var response CourseResponse
		response.ID = course.ID
		response.Title = course.Title
		response.Requirements = course.Requirements
		response.Levels = course.Levels
		response.ImageURL = course.ImageURL
		response.Category.ID = course.Category.ID
		response.Category.CategoryName = course.Category.CategoryName
		response.Description = course.Description
		response.Price = course.Price
		response.Author = course.Author
		response.CreatedAt = course.CreatedAt.Format("2006-01-02 15:04:05")
		response.UpdatedAt = course.UpdatedAt.Format("2006-01-02 15:04:05")
		results = append(results, response)
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"data": []string{},
		})
	}

	return c.JSON(fiber.Map{
		"data": results,
	})
}