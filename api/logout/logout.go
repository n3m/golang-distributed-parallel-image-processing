package logout

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/login"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func LogoutResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/logout")
	user := c.Get("user").(*jwt.Token)

	if _, ok := login.ActiveTokens[user.Raw]; ok {
		fmt.Println("\t[OPERATION] Revoked token")
		delete(login.ActiveTokens, user.Raw)
		return c.JSON(http.StatusOK, &map[string]string{"message": "the provided token has been revoked"})
	}
	return c.JSON(http.StatusOK, &map[string]string{"message": "Your token is invalid"})
}
