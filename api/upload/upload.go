package upload

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

func UploadResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/upload")
	token := c.FormValue("worker-token")
	workloadName := c.FormValue("workload_id")

	valid := helpers.IsBotTokenActive(token)
	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}

	image, err := c.FormFile("image")
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] There was no file sent into the input. ("+err.Error()+")")
	}

	src, err := image.Open()
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error opening image. ("+err.Error()+")")
	}
	defer src.Close()

	fileURLOnServer := "public/results/" + workloadName + "/" + image.Filename
	dst, err := os.Create(fileURLOnServer)
	if err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error creating Image on Server. ("+err.Error()+")")
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return helpers.ReturnJSON(c, http.StatusConflict, "[ERR] Error copying information on new image. ("+err.Error()+")")
	}

	fmt.Println("\t[OPERATION] Uploaded Image")
	return helpers.ReturnJSONMap(c, http.StatusOK, map[string]interface{}{
		"message":  "An image has been successfully uploaded",
		"filename": image.Filename,
		"size":     strconv.Itoa(int(image.Size)),
	})
}
