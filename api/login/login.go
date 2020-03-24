package login

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func LoginResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/login")
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "admin" && password == "password" {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{

			"message": "Welcome " + username,
			"token":   t,
		})
	}
	return echo.ErrUnauthorized
}
