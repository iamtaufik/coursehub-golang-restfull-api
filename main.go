package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/database"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/routes"
)

func init () {
	database.StartDB()
}


func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}