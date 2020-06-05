package workloads

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"golang-distributed-parallel-image-processing/scheduler"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NoOfTests int = 0
var NoOfJobs int = 0

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

func WorkloadsFilterResponse(c echo.Context) error {
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

	lastJobID := NoOfJobs
	NoOfJobs++

	/* Params */

	workloadID := c.FormValue("workload-id")
	filter := c.FormValue("filter")

	/* Folder Creation */

	_ = os.MkdirAll("public/download/"+workloadID+"/", 0755)
	_ = os.MkdirAll("public/results/"+workloadID+"/", 0755)

	/* File receiving */

	image, err := c.FormFile("data")
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] There was no file sent into the input. ("+err.Error()+")")
	}

	src, err := image.Open()
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error opening image. ("+err.Error()+")")
	}
	defer src.Close()

	objID := primitive.NewObjectID()

	fileURLOnServer := "public/download/" + workloadID + "/" + objID.Hex() + "_" + image.Filename
	downloadURLOnServer := "download/" + workloadID + "/" + objID.Hex() + "_" + image.Filename
	dst, err := os.Create(fileURLOnServer)
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error creating Image on Server. ("+err.Error()+")")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error copying information on new image. ("+err.Error()+")")
	}

	fileID := int64(0)

	/* Get File ID */
	if _, ok := cc.WorkloadsFileID[workloadID]; ok {
		cc.WorkloadsFileID[workloadID] = cc.WorkloadsFileID[workloadID] + int64(1)
		fileID = cc.WorkloadsFileID[workloadID]
	} else {
		cc.WorkloadsFileID[workloadID] = int64(1)
		fileID = int64(1)
	}

	fileExt := strings.Split(image.Filename, ".")[1]

	preJobString := workloadID + "|" + filter + "|" + strconv.Itoa(int(fileID)) + "|" + downloadURLOnServer + "|" + fileExt

	/*RPC Job*/
	cc.JOBS <- scheduler.Job{RPCName: "filter", Data: preJobString}

	return helpers.ReturnJSONMap(c, http.StatusOK, map[string]interface{}{
		"Workload ID": workloadID,
		"Job ID":      lastJobID,
		"Status":      "Scheduling",
		"Results":     "http://localhost:8080/results/" + workloadID + "/",
	})
}
