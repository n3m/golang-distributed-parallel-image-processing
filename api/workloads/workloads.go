package workloads

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func WorkloadsResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/workloads/test")
	user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	token := user.Raw

	valid := helpers.IsTokenActive(token)

	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}

	return helpers.ReturnJSON(c, http.StatusOK, "A test is now running!")
}
