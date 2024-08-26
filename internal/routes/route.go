package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/handler"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	course := api.Group("/courses")
	course.Get("/", middleware.Protected(), handler.GetAllCourses)
	course.Post("/", handler.CreateCourse)

	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)
}