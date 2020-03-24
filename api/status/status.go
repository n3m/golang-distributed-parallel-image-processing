package status

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func StatusResponse(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.JSON(http.StatusOK, &Message{
			Message:"Welcome" + name + "!",
		}
	)
}
