package helpers

import (
	"github.com/labstack/echo"
)

var ActiveTokens map[string]bool = make(map[string]bool)

func IsTokenActive(token string) bool {
	if _, ok := ActiveTokens[token]; ok {
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
