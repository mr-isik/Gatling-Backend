package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type contextKey string

const claimsKey contextKey = "claims"

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": domain.ErrUnauthorized.Error()})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": domain.ErrUnauthorized.Error()})
		}

		tokenStr := parts[1]
		claims, err := service.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": domain.ErrUnauthorized.Error()})
		}

		// Store claims in user context
		ctx := context.WithValue(c.UserContext(), claimsKey, claims)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func GetClaimsFromContext(ctx context.Context) *service.Claims {
	claims, ok := ctx.Value(claimsKey).(*service.Claims)
	if !ok {
		return nil
	}
	return claims
}

func GetUserIDFromContext(ctx context.Context) string {
	claims := GetClaimsFromContext(ctx)
	if claims == nil {
		return ""
	}
	return claims.UserID
}
