package login

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func LoginResponse(c echo.Context) error {

	fmt.Println("[ACCESS] New connection to:\t/login")
	username, password, isOk := c.(*helpers.CustomContext).Context.Request().BasicAuth()
	if !isOk {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERROR] Problem parsing Basic Auth")
	}
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
		helpers.ActiveTokens[t] = true
		fmt.Println("\t[OPERATION] Generated token")
		return c.JSON(http.StatusOK, map[string]string{

			"message": "Welcome " + username,
			"token":   t,
		})
	}
	return helpers.ReturnJSON(c, http.StatusConflict, "[ERROR] Credentials are wrong")
}
