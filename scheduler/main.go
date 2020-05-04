package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	"golang-distributed-parallel-image-processing/models"
	pb "golang-distributed-parallel-image-processing/proto"

	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
	Data    string
}

func schedule(job Job) {
	/* Load Distribution */

	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewControllerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = c.CreateJob(ctx, &pb.JobRequest{Msg: job.RPCName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// log.Printf("Scheduler -> Task: %+v was completed by: %+v", job.RPCName, r.Msg)
}

func Start(jobs chan Job, DB map[string]models.Worker) error {
	for {
		job := <-jobs

		lowestUsage := 99999
		lowestPort := 0
		for _, data := range DB {
			if data.Usage < lowestUsage {
				lowestPort = data.Port
				lowestUsage = data.Usage
			}
		}

		if lowestPort == 0 {
			return nil
		}

		job.Address = "localhost:" + strconv.Itoa(lowestPort)

		go schedule(job)
	}
	return nil
}
