package workloads

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"golang-distributed-parallel-image-processing/scheduler"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var NoOfTests int = 0

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

	NoOfTests++
	for e := 0; e < 20; e++ {
		cc.JOBS <- scheduler.Job{RPCName: "test"}
	}

	return helpers.ReturnJSONMap(c, http.StatusOK, map[string]interface{}{
		"Workload": "test",
		"Job ID":   NoOfTests,
		"Status":   "Completed 20 tasks",
		"Result":   "Done!",
	})
}
