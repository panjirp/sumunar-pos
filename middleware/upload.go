package middleware

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

const maxFileSize = 2 * 1024 * 1024 // 2MB

func ValidateImageFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("logo")
		if err != nil {
			return echo.NewHTTPError(400, "logo file is required")
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			return echo.NewHTTPError(400, fmt.Sprintf("file extension %s is not allowed", ext))
		}

		// Size check
		if file.Size > maxFileSize {
			return echo.NewHTTPError(400, "file is too large (max 2MB)")
		}

		return next(c)
	}
}
