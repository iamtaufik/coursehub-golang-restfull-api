package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	searchParams := c.Query("search")

	if searchParams != "" {
		var courses []model.Course
		db.Preload("Category").Where("title ILIKE ?", "%"+searchParams+"%").Find(&courses)

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
				"message": "No course found",
			})
		}

		return c.JSON(results)
	}


		// Select specific columns id, title, description, price, author
	
		
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

func GetDetailCourse(c *fiber.Ctx) error {
	db := database.StartDB()

	courseID := c.Params("id")

	var course model.Course
	if err := db.Preload("Chapters").Preload("Chapters.Modules").Preload("Category").Where("id = ?", courseID).First(&course).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Course not found",
		})
	}

	// Struct untuk hasil yang akan di-return
	type Result struct {
		ID          uint    `json:"id"`
		Title       string  `json:"title"`
		Requirements string `json:"requirements"`
		Levels      model.Levels  `json:"levels"`
		ImageURL    string  `json:"imageURL"`
		Category    struct {
			ID           uint   `json:"id"`
			CategoryName string `json:"category_name"`
		} `json:"category"`
		Chapters []struct {
			ID      uint   `json:"id"`
			Name    string `json:"name"`
			Modules []struct {
				ID        uint   `json:"id"`
				Title     string `json:"title"`
				VideoURL  string `json:"videoURL"`
				Duration  int    `json:"duration"`
				IsTrailer bool   `json:"isTrailer"`
			} `json:"modules"`
		} `json:"chapters"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Author      string  `json:"author"`
		CreatedAt   string  `json:"createdAt"`
		UpdatedAt   string  `json:"updatedAt"`
	}

	// Map data dari model ke struct Result
	var result Result
	result.ID = course.ID
	result.Title = course.Title
	result.Requirements = course.Requirements
	result.Levels = course.Levels
	result.ImageURL = course.ImageURL
	result.Description = course.Description
	result.Price = course.Price
	result.Author = course.Author
	result.CreatedAt = course.CreatedAt.Format("2006-01-02 15:04:05")
	result.UpdatedAt = course.UpdatedAt.Format("2006-01-02 15:04:05")

	// Map data Category
	result.Category.ID = course.Category.ID
	result.Category.CategoryName = course.Category.CategoryName

	// Map data Chapters dan Modules
	for _, chapter := range course.Chapters {
		var chapterResult struct {
			ID      uint   `json:"id"`
			Name    string `json:"name"`
			Modules []struct {
				ID        uint   `json:"id"`
				Title     string `json:"title"`
				VideoURL  string `json:"videoURL"`
				Duration  int    `json:"duration"`
				IsTrailer bool   `json:"isTrailer"`
			} `json:"modules"`
		}
		chapterResult.ID = chapter.ID
		chapterResult.Name = chapter.Name

		for _, module := range chapter.Modules {
			var moduleResult struct {
				ID        uint   `json:"id"`
				Title     string `json:"title"`
				VideoURL  string `json:"videoURL"`
				Duration  int    `json:"duration"`
				IsTrailer bool   `json:"isTrailer"`
			}
			moduleResult.ID = module.ID
			moduleResult.Title = module.Title
			moduleResult.VideoURL = module.VideoURL
			moduleResult.Duration = module.Duration
			moduleResult.IsTrailer = module.IsTrailer

			chapterResult.Modules = append(chapterResult.Modules, moduleResult)
		}

		result.Chapters = append(result.Chapters, chapterResult)
	}

	return c.JSON(result)
}

func JoinCourse(c *fiber.Ctx) error {
	db := database.StartDB()

	userClaim := c.Locals("user").(*jwt.Token)
	claims := userClaim.Claims.(jwt.MapClaims)

	email := claims["email"].(string)
	courseID := c.Params("id")

	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	var course model.Course
	if err := db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Course not found",
		})
	}

	var existingCourse model.Course
	if err := db.Model(&user).Where("id = ?", courseID).Association("Courses").Find(&existingCourse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check if user joined course",
		})
	}

	if existingCourse.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already joined the course",
		})
	}

	if err := db.Model(&user).Association("Courses").Append(&course); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to join course",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Successfully join course",
	})
}

func DeleteCourse(c *fiber.Ctx) error {
	db := database.StartDB()

	courseID := c.Params("id")

	var course model.Course
	if err := db.Where("id = ?", courseID).First(&course).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Course not found",
		})
	}

	if err := db.Delete(&course).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete course",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Course deleted successfully",
	})
}