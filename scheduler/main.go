package scheduler

import (
	"context"
	"log"
	"sort"
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
	r, err := c.CreateJob(ctx, &pb.JobRequest{Msg: job.RPCName, Args: job.Data})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Scheduler -> Task: %+v was completed by: %+v", job.RPCName, r.GetMsg())
}

func Start(jobs chan Job, DB map[string]models.Worker) error {

	type UsageMonitor struct {
		Data  models.Worker
		Usage int
	}

	for {
		job := <-jobs

		ss := []UsageMonitor{}
		for _, data := range DB {
			ss = append(ss, UsageMonitor{Data: data, Usage: data.Usage})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Usage > ss[j].Usage
		})

		thePort := 0
		TheChoosenWorker := models.Worker{}

		for _, monitor := range ss {
			if monitor.Data.Status == "Overload" {
				continue
			}
			TheChoosenWorker = monitor.Data
		}

		thePort = TheChoosenWorker.Port
		if thePort == 0 {
			return nil
		}

		job.Address = "localhost:" + strconv.Itoa(thePort)

		go schedule(job)
	}
	return nil
}
