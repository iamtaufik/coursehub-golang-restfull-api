package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/database"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/dto"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/model"
)

func CreateCourse(c *fiber.Ctx) error {

	db := database.StartDB()

	var createCourseDTO dto.CreateCourseDTO
	if err := c.BodyParser(&createCourseDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Mapping Chapters and Modules from DTO
	var chapters []model.Chapter
	for _, chapterDTO := range createCourseDTO.Chapters {
		var modules []model.Module
		for _, moduleDTO := range chapterDTO.Modules {
			module := model.Module{
				Title:    moduleDTO.Title,
				Duration: moduleDTO.Duration,
				VideoURL: moduleDTO.VideoURL,
				IsTrailer: moduleDTO.IsTrailer,
			}
			modules = append(modules, module)
		}

		chapter := model.Chapter{
			Name:    chapterDTO.Name,
			Modules: modules,
		}
		chapters = append(chapters, chapter)
	}

	// Create Course model with mapped data
	course := model.Course{
		Title:        createCourseDTO.Title,
		Description:  createCourseDTO.Description,
		Requirements: createCourseDTO.Requirements,
		Levels:       model.Levels(createCourseDTO.Levels),
		ImageURL:     createCourseDTO.ImageURL,
		CategoryID:   createCourseDTO.CategoryID,
		Price:        createCourseDTO.Price,
		Author:       createCourseDTO.Author,
		Chapters:     chapters,
	}

	if err := db.Omit("Category").Create(&course).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create course",
			"error":   err.Error(),
		})
	}

	response := dto.CreateCourseDTO{
		Title:        course.Title,
		Description:  course.Description,
		Requirements: course.Requirements,
		Levels:       model.Levels(course.Levels),
		ImageURL:     course.ImageURL,
		CategoryID:   course.CategoryID,
		Price:        course.Price,
		Author: 	 course.Author,
		Chapters:     createCourseDTO.Chapters,
	}


	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Course created successfully",
		"course":  response,
	})
}

func GetAllCourses(c *fiber.Ctx) error {
	db := database.StartDB()
		// Select specific columns id, title, description, price, author
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
		
		var courses []model.Course
		db.Preload("Category").Find(&courses)

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

		return c.JSON(results)
}