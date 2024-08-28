package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/handler"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	category := api.Group("/categories")
	category.Get("/", handler.GetAllCategories)
	category.Get("/:id", handler.GetCoursesByCategoryID)

	course := api.Group("/courses")
	course.Get("/", handler.GetAllCourses)
	course.Post("/", middleware.Protected(), middleware.IsAdmin, handler.CreateCourse)
	course.Delete("/:id", middleware.Protected(), middleware.IsAdmin, handler.DeleteCourse)
	course.Get("/:id", handler.GetDetailCourse)
	course.Get("/:id/join",middleware.Protected(), handler.JoinCourse)


	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	user := api.Group("/users")
	user.Get("/me",middleware.Protected(), handler.GetMe)
}