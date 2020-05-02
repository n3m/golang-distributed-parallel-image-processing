package main

import (
	"fmt"
	"golang-distributed-parallel-image-processing/api"
	"golang-distributed-parallel-image-processing/api/helpers"
	"golang-distributed-parallel-image-processing/controller"
	"golang-distributed-parallel-image-processing/models"
	"log"

	"github.com/labstack/echo"
	"github.com/sonyarouje/simdb/db"
)

func main() {
	/* Setup Variables */
	ControllerConnectionURL := "tcp://localhost:40899"
	APIPort := ":8080"
	DBName := "workers"
	currentWorkers := map[string]models.Worker{}

	/* DB Setup */
	db, err := db.New(DBName)
	if err != nil {
		panic(err)
	}

	/* Controller Setup */
	go controller.Start(ControllerConnectionURL, currentWorkers, db)
	log.Printf("[SETUP] Controller Connection URL: %+v", ControllerConnectionURL)

	/* API EndPoint Setup */
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &helpers.CustomContext{c, currentWorkers}
			return next(cc)
		}
	})

	modules := api.LoadModules()
	fmt.Println("== URLs Loaded == ")
	for _, mod := range modules {
		switch mod.Method {
		case "GET":
			fmt.Println("\tGET:\t" + mod.Path)
			if mod.Middleware != nil {
				e.GET(mod.Path, mod.Function, *mod.Middleware)
			} else {
				e.GET(mod.Path, mod.Function)
			}
			break
		case "POST":
			fmt.Println("\tPOST:\t" + mod.Path)
			if mod.Middleware != nil {
				e.POST(mod.Path, mod.Function, *mod.Middleware)
			} else {
				e.POST(mod.Path, mod.Function)
			}
			break
		}
	}

	go e.Logger.Fatal(e.Start(APIPort))

	/* Scheduler Setup */
	// jobs := make(chan scheduler.Job)
	// go scheduler.Start(jobs)

	// sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "hello"}

	// for {
	// 	sampleJob.RPCName = fmt.Sprintf("hello-%v", rand.Intn(10000))
	// 	jobs <- sampleJob
	// 	time.Sleep(time.Second * 5)
	// }
}
