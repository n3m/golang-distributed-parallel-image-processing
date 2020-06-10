package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

var maxJobs int = 100

var (
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	usage             = 0
	status            = "Idle"
	port              = 0
	jobsDone          = 0
	storeToken        = ""
	storeEndpoint     = ""
)

/* System Functions */

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller address")
	flag.StringVar(&workerName, "node-name", "hard-worker", "Worker Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "Comma-separated worker tags")
	flag.StringVar(&storeToken, "image-store-token", "token", "Image Store Token")
	flag.StringVar(&storeEndpoint, "image-store-endpoint", "url_endpoint", "Image Store URL")
}

/* Response Functions */

func (s *server) CreateJob(ctx context.Context, in *pb.JobRequest) (*pb.JobReply, error) {
	switch in.Msg {
	case "filter":
		jobsDone++
		usage++
		status = "Running"
		/* DO TASK */
		in.GetArgs()
		args := in.GetArgs()
		workID, filter, fileID, fileNameOnServer, fileExt := getArgsForFilter(args)
		log.Printf("[Worker] %+v: I've been called to do a filter: %+v", workerName, filter)

		/*Create Repo*/
		_ = os.MkdirAll(workerName+"/"+workID+"/", 0755)

		//Download file
		downloadURL := storeEndpoint + "/download"
		urlOnTheWorker := workerName + "/" + workID + "/" + fileID + "." + fileExt

		err := DownloadFileFromServer(workID, storeToken, fileNameOnServer, downloadURL, urlOnTheWorker)
		if err != nil {
			log.Printf("[ERR] Error downloading file: %+v", err.Error())
		}

		// Switch between the selected filters and process file
		switch filter {
		case "filtertest":
			err := SendPostRequest(workID, storeToken, storeEndpoint+"/upload", urlOnTheWorker, "image")
			if err != nil {
				log.Printf("%+v", err.Error())
			}

			break
		case "greyscale":
			cmd := exec.Command("filter", "-p "+urlOnTheWorker)

			if _, err := cmd.Output(); err != nil {
				log.Printf("%+v", err.Error())
			} else {
				//Upload the file
				err := SendPostRequest(workID, storeToken, storeEndpoint+"/upload", urlOnTheWorker, "image")
				if err != nil {
					log.Printf("%+v", err.Error())
				}
			}
			break
		default:
			break
		}

		/* END TASK*/
		response := &pb.JobReply{Msg: workerName, Args: "done"}
		usage--
		return response, nil
	case "test":
		jobsDone++
		log.Printf("[Worker] %+v: I've been called to do a test", workerName)
		usage++
		status = "Running"
		/* DO TASK */
		time.Sleep(time.Second * 5)
		/* END TASK*/
		response := &pb.JobReply{Msg: workerName}
		usage--
		return response, nil
	default:
		/*Log*/
		jobsDone++
		/*Task*/
		log.Printf("[Worker] %+v: I've been called by an RPC, but no task was received", workerName)
		usage++
		status = "Running"
		response := &pb.JobReply{Msg: "RPC not valid"}
		usage--
		status = "Idle"
		return response, nil
	}
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
	for {
		_, err = socket.Recv()
		if err != nil {
			die(errorMessage + "Error while Recv() ->" + err.Error())
		}

		err = socket.Send([]byte(createDataString()))
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

	rpcPort := getAvailablePort()
	port = rpcPort

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
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

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for range ticker.C {
			if usage > maxJobs {
				status = "Overload"
			} else {
				status = "Idle"
			}
		}
	}()
}

func createDataString() string {
	data := workerName + "|" + status + "|" + strconv.Itoa(usage) + "|" + tags + "|" + strconv.Itoa(port) + "|" + strconv.Itoa(jobsDone) + "|" + storeToken
	return data
}

func DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getArgsForFilter(args string) (string, string, string, string, string) {
	data := strings.Split(args, "|")
	return data[0], data[1], data[2], data[3], data[4]
}

func SendPostRequest(workID string, token string, urlToSend string, filename string, filetype string) error {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/*Initial Request Values*/
	requestBody := &bytes.Buffer{}
	multiPartWriter := multipart.NewWriter(requestBody)

	/* Insert File */
	part, err := multiPartWriter.CreateFormFile("image", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}

	/* Insert Token */
	tokenSender, err := multiPartWriter.CreateFormField("worker-token")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = tokenSender.Write([]byte(token))
	if err != nil {
		log.Fatalln(err)
	}

	/* Insert WorkloadID */
	workloadSender, err := multiPartWriter.CreateFormField("workload_id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = workloadSender.Write([]byte(workID))
	if err != nil {
		log.Fatalln(err)
	}

	/*Close RequestBody*/
	multiPartWriter.Close()

	/*Setup request*/
	request, err := http.NewRequest("POST", urlToSend, requestBody)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	/*Create Request and Send*/
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("[SendPostRequest()] Could not upload file.")
	}

}

func DownloadFileFromServer(workID string, token string, image_id string, urlToSend string, filepath string) error {
	/*Initial Request Values*/
	requestBody := &bytes.Buffer{}
	multiPartWriter := multipart.NewWriter(requestBody)

	/* Insert Token */
	tokenSender, err := multiPartWriter.CreateFormField("worker-token")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = tokenSender.Write([]byte(token))
	if err != nil {
		log.Fatalln(err)
	}

	/* Insert WorkloadID */
	workloadSender, err := multiPartWriter.CreateFormField("workload_id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = workloadSender.Write([]byte(workID))
	if err != nil {
		log.Fatalln(err)
	}

	/* Insert image_id */
	imageSender, err := multiPartWriter.CreateFormField("image_id")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = imageSender.Write([]byte(image_id))
	if err != nil {
		log.Fatalln(err)
	}

	/*Close RequestBody*/
	multiPartWriter.Close()

	/*Setup request*/
	request, err := http.NewRequest("POST", urlToSend, requestBody)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", multiPartWriter.FormDataContentType())

	/*Create Request and Send*/
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		// Create the file
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Writer the body to file
		_, err = io.Copy(out, response.Body)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("[DownloadFileFromServer()] Could not download file.")
	}

}
