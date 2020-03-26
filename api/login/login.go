package login

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var ActiveTokens map[string]bool = make(map[string]bool)

func LoginResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/login")
	username := c.FormValue("username")
	password := c.FormValue("password")
	if username == "admin" && password == "password" {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user"] = username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		claims["token"] = t
		ActiveTokens[t] = true
		fmt.Println("\t[OPERATION] Generated token")
		return c.JSON(http.StatusOK, map[string]string{

			"message": "Welcome " + username,
			"token":   t,
		})
	}
	return echo.ErrUnauthorized
}
