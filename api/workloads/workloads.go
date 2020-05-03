package workloads

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"golang-distributed-parallel-image-processing/scheduler"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func WorkloadsResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/workloads/test")
	user := c.Get("user").(*jwt.Token)
	token := user.Raw

	valid := helpers.IsTokenActive(token)

	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}

	cc := c.(*helpers.CustomContext)

	if len(cc.DB) == 0 {
		return helpers.ReturnJSON(c, http.StatusConflict, "There are no registered workers")
	}

	/*TEST*/
	cc.JOBS <- scheduler.Job{RPCName: "test"}

	return helpers.ReturnJSON(c, http.StatusOK, "A test is now running!")
}
