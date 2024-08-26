package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/config"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/database"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/dto"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	db := database.StartDB()

	var registerDTO dto.RegisterDTO

	if err := c.BodyParser(&registerDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	user := model.User{
		Email:    registerDTO.Email,
		Password: string(hashedPassword),
		Profile: model.Profile{
			FirstName: registerDTO.FirstName,
			LastName:  registerDTO.LastName,
		},
	}

	if err := db.Preload("Profile").Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	type Response struct {
		ID 			uint 	`json:"id"`
		FirstName	string	`json:"firstName"`
		LastName	string	`json:"lastName"`
		Email		string	`json:"email"`
	}

	response := Response{
		ID: user.ID,
		FirstName: user.Profile.FirstName,
		LastName: user.Profile.LastName,
		Email: user.Email,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    response,
	})
}

func GetUserByEmail(email string) (*model.User, error) {
	db := database.StartDB()
	
	var user model.User
	if err := db.Where(&model.User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil

}

func Login(c *fiber.Ctx) error {
	type LoginDTO struct {
		Email		string		`json:"email"`
		Password	string		`json:"password"`
	}

	type UserCredential struct {
		ID 			uint 	`json:"id"`
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}

	var loginDTO LoginDTO
	if err := c.BodyParser(&loginDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := new(model.User), *new(error)

	user, err = GetUserByEmail(loginDTO.Email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to query user",
			"error":   err.Error(),
		})
	} else if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	userCredential := UserCredential{
		ID: user.ID,
		Email: user.Email,
		Password: user.Password,
	}

	isMatch := bcrypt.CompareHashAndPassword([]byte(userCredential.Password), []byte(loginDTO.Password))
	if isMatch != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = userCredential.Email
	claims["id"] = userCredential.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // 3 days

	t, err := token.SignedString([]byte(config.Config("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":    t,
	})
}

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	email := claims["email"].(string)

	db := database.StartDB()

	var userDb model.User

	if err := db.Preload("Profile").Preload("Courses").Preload("Courses.Category").Where("email = ?", email).First(&userDb).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	type Result struct {
		ID        uint   `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Courses   []struct {
			ID          uint   `json:"id"`
			Title       string `json:"title"`
			Requirements string `json:"requirements"`
			Levels      model.Levels `json:"levels"` // Mengubah menjadi string jika model.Levels adalah enum atau string
			ImageURL    string `json:"imageURL"`
			Category    struct {
				ID           uint   `json:"id"`
				CategoryName string `json:"category_name"`
			} `json:"category"`
			Description string `json:"description"`
			Author      string `json:"author"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		} `json:"courses"`
	}

	var result Result

	// Memetakan data userDb ke dalam result
	result.ID = userDb.ID
	result.FirstName = userDb.Profile.FirstName
	result.LastName = userDb.Profile.LastName
	result.Email = userDb.Email

	// Memetakan data kursus ke dalam result
	for _, course := range userDb.Courses {
		var courseResult struct {
			ID          uint   `json:"id"`
			Title       string `json:"title"`
			Requirements string `json:"requirements"`
			Levels      model.Levels `json:"levels"`
			ImageURL    string `json:"imageURL"`
			Category    struct {
				ID           uint   `json:"id"`
				CategoryName string `json:"category_name"`
			} `json:"category"`
			Description string `json:"description"`
			Author      string `json:"author"`
			CreatedAt   string `json:"createdAt"`
			UpdatedAt   string `json:"updatedAt"`
		}

		courseResult.ID = course.ID
		courseResult.Title = course.Title
		courseResult.Requirements = course.Requirements
		courseResult.Levels = course.Levels
		courseResult.ImageURL = course.ImageURL
		courseResult.Description = course.Description
		courseResult.Author = course.Author
		courseResult.CreatedAt = course.CreatedAt.Format("2006-01-02 15:04:05")
		courseResult.UpdatedAt = course.UpdatedAt.Format("2006-01-02 15:04:05")

		courseResult.Category.ID = course.Category.ID
		courseResult.Category.CategoryName = course.Category.CategoryName

		result.Courses = append(result.Courses, courseResult)
	}

	// Return data pengguna yang sedang login beserta kursus yang dimiliki
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}