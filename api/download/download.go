package download

import (
	"golang-distributed-parallel-image-processing/api/helpers"
	"net/http"

	"github.com/labstack/echo"
)

func DownloadResponse(c echo.Context) error {
	workload_id := c.FormValue("workload_id")
	image_id := c.FormValue("image_id")
	token := c.FormValue("worker-token")

	valid := helpers.IsBotTokenActive(token)
	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}

	cc := c.(*helpers.CustomContext)

	if len(cc.DB) == 0 {
		return helpers.ReturnJSON(c, http.StatusConflict, "There are no registered workers")
	}

	/* Get File ID */
	if _, ok := cc.WorkloadsFileID[workload_id]; ok == false {
		return helpers.ReturnJSON(c, http.StatusConflict, "There's no workload registered under: "+workload_id)
	}

	return c.File("public/download/" + workload_id + "/" + image_id)
}
