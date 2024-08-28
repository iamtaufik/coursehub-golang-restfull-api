package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/config"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config("JWT_SECRET"))},
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Missing or malformed JWT",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired JWT",
			})
		},
	})
}


func IsAdmin(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized, invalid token",
		})
	}

	// Ambil klaim dari token
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized, invalid claims",
		})
	}

	// Cek apakah email sesuai dengan email admin
	email := claims["email"].(string)
	if email != config.Config("ADMIN_EMAIL") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized, not an admin",
		})
	}

	// Jika semuanya sesuai, lanjutkan ke handler berikutnya
	return c.Next()
}