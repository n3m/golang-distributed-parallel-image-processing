package status

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api/helpers"
	"net/http"
	"strconv"
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

	workers := []map[string]interface{}{}
	cc := c.(*helpers.CustomContext)
	for _, data := range cc.DB {
		if data.Active {
			// var jsonData map[string]interface{}
			// byteData, err := json.Marshal(data)
			// if err != nil {
			// 	return helpers.ReturnJSON(c, http.StatusConflict, "Error Marshaling data")
			// }
			// err = json.Unmarshal(byteData, &jsonData)
			// if err != nil {
			// 	return helpers.ReturnJSON(c, http.StatusConflict, "Error Unmarshaling data")
			// }
			// workers = append(workers, jsonData)
			workers = append(workers, map[string]interface{}{
				"Worker":    data.Name,
				"Tags":      data.Tags,
				"Status":    data.Status,
				"Usage":     strconv.Itoa(data.Usage) + "%",
				"Jobs Done": data.JobsDone,
				"Port":      data.Port,
			})
		}
	}

	return helpers.ReturnJSONMap(c, http.StatusOK, map[string]interface{}{
		"message": "Hi " + claims["user"].(string) + ", the Distributed Parallel Image Processing System is up and running!",
		"time":    time.Now().String(),
		"workers": workers,
	})
}

func StatusWorkerResponse(c echo.Context) error {
	fmt.Println("[ACCESS] New connection to:\t/status/:worker")
	user := c.Get("user").(*jwt.Token)
	token := user.Raw
	valid := helpers.IsTokenActive(token)

	if !valid {
		return helpers.ReturnJSON(c, http.StatusConflict, "Token is invalid or revoked")
	}

	cc := c.(*helpers.CustomContext)

	worker := c.Param("worker")
	if workerData, ok := cc.DB[worker]; ok {
		if workerData.Active {
			return helpers.ReturnJSONMap(cc.Context, http.StatusOK, map[string]interface{}{
				"Worker": workerData.Name,
				"Tags":   workerData.Tags,
				"Status": workerData.Status,
				"Usage":  strconv.Itoa(workerData.Usage) + "%",
			})
		} else {
			return helpers.ReturnJSON(c, http.StatusConflict, "Worker is not active anymore!")
		}

	}
	return helpers.ReturnJSON(c, http.StatusNotFound, "Worker wasn't found!")
}
