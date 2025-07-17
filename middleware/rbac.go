package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			if userRole == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Role not found in context")
			}

			if userRole.(string) != requiredRole {
				return echo.NewHTTPError(http.StatusForbidden, "You do not have access to this resource")
			}

			return next(c)
		}
	}
}

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "No role found")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden: insufficient permissions")
		}
	}
}
