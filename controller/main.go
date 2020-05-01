package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sonyarouje/simdb/db"
	"go.nanomsg.org/mangos"

	// register transports
	"go.nanomsg.org/mangos/protocol/surveyor"
	_ "go.nanomsg.org/mangos/transport/all"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Start(controllerAddress string, currentWorkers map[string]interface{}, db *db.Driver) {
	errorMessage := "[ERR] Controller -> "
	socket, err := surveyor.NewSocket()
	if err != nil {
		die(errorMessage+"Couldn't get a socket connection -> %+v", err.Error())
	}

	err = socket.Listen(controllerAddress)
	if err != nil {
		die(errorMessage+"Couldn't listen on \""+controllerAddress+"\" -> %+v", err.Error())
	}
	err = socket.SetOption(mangos.OptionSurveyTime, time.Second/2)
	if err != nil {
		die(errorMessage+"SetOption(): %+v", err.Error())
	}
	seconds := 0

	for {
		time.Sleep(time.Second)
		if seconds == 9 {
			log.Printf("Current Online Workers: %+v, -> %v", len(currentWorkers), currentWorkers)
		}
		if seconds < 10 {
			seconds += 1
		} else {
			seconds = 0
		}

		err = socket.Send([]byte("Is anyone there?"))
		if err != nil {
			die(errorMessage+"The process for looking for workers failed! -> %+v", err.Error())
		}

		for {
			msg, err := socket.Recv()
			if err != nil {
				break
			}
			workerName := string(msg)
			if currentWorkers[workerName] == nil {
				log.Printf("[C] Controller: Found a new Worker")
				currentWorkers[workerName] = 1
			}

		}
	}
}
