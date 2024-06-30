package middleware

import (
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

func InfoLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestTime := time.Now()
		method := c.Request().Method
		path := c.Request().URL.Path
		err := next(c)
		responseTime := time.Since(requestTime)
		log.Printf("- %s request to %s took %v\n", method, path, responseTime)
		return err
	}
}
