package restful

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NotFoundHandler function ...
func NotFoundHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Message: "endpoint not found", Data: nil, Error: nil})
	}
}
