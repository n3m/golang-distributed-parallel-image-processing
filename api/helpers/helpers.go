package helpers

import (
	"golang-distributed-parallel-image-processing/api/login"

	"github.com/labstack/echo"
)

func IsTokenActive(token string) bool {
	if _, ok := login.ActiveTokens[token]; ok {
		return true
	}
	return false
}

func ReturnJSON(c echo.Context, status int, data string) error {
	return c.JSON(status, &map[string]string{"message": data})
}

func ReturnJSONMap(c echo.Context, status int, data map[string]string) error {
	return c.JSON(status, &data)
}
