package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "golang-distributed-parallel-image-processing/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/profiling/proto"

	// register transports
	"go.nanomsg.org/mangos/protocol/respondent"
	_ "go.nanomsg.org/mangos/transport/all"
)

var (
	defaultRPCPort = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	proto.UnimplementedProfilingServer
}

var (
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	usage             = 0
	status            = "Idle"
	port              = 0
)

/* System Functions */

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller address")
	flag.StringVar(&workerName, "worker-name", "hard-worker", "Worker Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "Comma-separated worker tags")
}

/* Response Functions */

func (s *server) CreateJob(ctx context.Context, in *pb.JobRequest) (*pb.JobReply, error) {
	switch in.Msg {
	case "test":
		return &pb.JobReply{Msg: "Ping Pong!"}, nil
	default:
		return &pb.JobReply{Msg: "RPC not valid"}, nil
	}
}

func joinCluster() {
	errorMessage := "[ERR] Worker: (" + workerName + ") -> "
	workerData := workerName + "|" + status + "|" + strconv.Itoa(usage) + "|" + tags + "|" + strconv.Itoa(port)
	var err error
	socket, err := respondent.NewSocket()
	if err != nil {
		die(errorMessage + err.Error())
	}

	err = socket.Dial(controllerAddress)
	if err != nil {
		die(errorMessage + err.Error())
	}
	seconds := 0
	for {
		_, err = socket.Recv()
		if err != nil {
			die(errorMessage + "Error while Recv() ->" + err.Error())
		}
		if seconds == 9 {
			log.Printf("[Worker] %v: I've been pinged! ", workerName)
		}
		if seconds < 10 {
			seconds++
		} else {
			seconds = 0
		}

		err = socket.Send([]byte(workerData))
		if err != nil {
			die(errorMessage + err.Error())
		}
	}

}

func getAvailablePort() int {
	port := defaultRPCPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = port + 1
			continue
		}
		ln.Close()
		break
	}
	return port
}

func main() {
	flag.Parse()

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
	rpcPort := getAvailablePort()
	port = rpcPort
	log.Printf("[W] "+workerName+": Starting RPC Service on localhost:%v", rpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterControllerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
