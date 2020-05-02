package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

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

func (s *server) ResponseDetails(ctx context.Context, in *pb.DetailsRequest) (*pb.DetailsReply, error) {
	return &pb.DetailsReply{Status: status, Workername: workerName, Workload: int64(usage)}, nil
}

func (s *server) ResponseStatus(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	return &pb.StatusReply{Status: status}, nil
}

func (s *server) ResponseWorkload(ctx context.Context, in *pb.WorkloadRequest) (*pb.WorkloadReply, error) {
	return &pb.WorkloadReply{Workload: int64(usage)}, nil
}

func (s *server) ResponseWorkerName(ctx context.Context, in *pb.WorkerNameRequest) (*pb.WorkerNameReply, error) {
	return &pb.WorkerNameReply{Workername: workerName}, nil
}

func (s *server) ResponsePong(ctx context.Context, in *pb.PingRequest) (*pb.PongReply, error) {
	return &pb.PongReply{Msg: "Pong"}, nil
}

func joinCluster() {
	errorMessage := "[ERR] Worker: (" + workerName + ") -> "
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
			log.Printf("[Worker] %v:I've been pinged! ", workerName)
		}
		if seconds < 10 {
			seconds++
		} else {
			seconds = 0
		}

		err = socket.Send([]byte(workerName))
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
