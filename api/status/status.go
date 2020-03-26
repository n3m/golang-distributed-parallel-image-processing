package status

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func StatusResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/status")
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	token := user.Raw

	valid := helpers.IsTokenActive(token)

	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}
	return helpers.ReturnJSONMap(c, http.StatusOK, map[string]string{
		"message": "Hi " + claims["user"].(string) + ", the Distributed Parallel Image Processing System is up and running!",
		"time":    time.Now().String(),
	})
}
