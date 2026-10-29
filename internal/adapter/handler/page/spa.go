package page

import (
	"github.com/labstack/echo/v5"
)

func ServeSPA(c *echo.Context) error {
	c.Response().Header().Set("Cache-Control", "no-cache")
	return c.File("static/app/index.html")
}
