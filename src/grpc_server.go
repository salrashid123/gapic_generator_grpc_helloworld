package main

import (
	"echo"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	grpcport = flag.String("grpcport", "", "grpcport")
	hs       *health.Server

	workRequestMap map[string]workRequest
)

const (
	address string = ":50051"
)

type server struct{}

type operationsServer struct{}

type workRequest struct {
	ID  string
	Req echo.EchoRequest
}

type contextKey string

func (s *operationsServer) GetOperation(ctx context.Context, in *lropb.GetOperationRequest) (*lropb.Operation, error) {
	log.Println("GetOperation: ", in.Name)

	var answer *lropb.Operation
	chances := rand.Intn(100)

	if wr, ok := workRequestMap[in.Name]; ok {
		if chances >= 70 {
			log.Printf("LRO Complete %v", in.Name)
			answer = &lropb.Operation{
				Name: in.Name,
				Done: true,
			}
			resp, _ := ptypes.MarshalAny(
				&echo.EchoReply{Message: "Hello Callback " + wr.Req.Name},
			)
			delete(workRequestMap, in.Name)
			answer.Result = &lropb.Operation_Response{Response: resp}
		} else {
			answer = &lropb.Operation{
				Name: in.Name,
				Done: false,
			}
		}
	} else {
		return &lropb.Operation{
			Name: in.Name,
			Done: false,
		}, status.Errorf(codes.NotFound, "Operation %q not found.", in.Name)
	}
	return answer, nil
}

func (s operationsServer) CancelOperation(ctx context.Context, in *lropb.CancelOperationRequest) (*empty.Empty, error) {
	log.Println("CancelOperation: ", in.Name)
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot cancel operation without a name.")
	}
	if _, ok := workRequestMap[in.Name]; ok {
		// ..just delete it to cancel
		delete(workRequestMap, in.Name)
	} else {
		return nil, status.Errorf(codes.NotFound, "cannot cancel unknown entry %v", in.Name)
	}
	return &empty.Empty{}, nil
}

func (s operationsServer) DeleteOperation(ctx context.Context, in *lropb.DeleteOperationRequest) (*empty.Empty, error) {
	log.Println("DeleteOperation: ", in.Name)
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot delete operation without a name.")
	}
	if _, ok := workRequestMap[in.Name]; ok {
		delete(workRequestMap, in.Name)
	} else {
		return nil, status.Errorf(codes.NotFound, "cannot delete unknown entry %v", in.Name)
	}
	return &empty.Empty{}, nil
}

func (s operationsServer) ListOperations(ctx context.Context, in *lropb.ListOperationsRequest) (*lropb.ListOperationsResponse, error) {
	// unimplemented, really ...w'ere just returning all the operations as falures
	log.Println("ListOperations: ")
	var operations []*lropb.Operation
	for k := range workRequestMap {
		operations = append(operations, &lropb.Operation{
			Name: k,
			Done: false,
		})
	}
	return &lropb.ListOperationsResponse{
		Operations: operations,
	}, nil
}

// WaitOperation randomly waits and returns an operation with the same name
// https://github.com/googleapis/gapic-showcase/blob/master/server/services/operations_service.go#L178
func (s operationsServer) WaitOperation(ctx context.Context, in *lropb.WaitOperationRequest) (*lropb.Operation, error) {
	if in.Name == "" {
		return nil, status.Error(codes.NotFound, "cannot wait on a operation without a name.")
	}
	if _, ok := workRequestMap[in.Name]; ok {
		num := rand.Intn(500)
		time.Sleep(time.Duration(num) * time.Millisecond)
		var result *lropb.Operation_Response
		if num%2 == 0 {
			result = &lropb.Operation_Response{}
		}
		return &lropb.Operation{
			Name:   in.Name,
			Done:   result != nil,
			Result: result,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "cannot wait on unknown operation %v", in.Name)

}

func (s *server) SayHello(ctx context.Context, in *echo.EchoRequest) (*echo.EchoReply, error) {

	log.Println("Got SayHello --> ", in.Name)

	var h, err = os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname %v", err)
	}
	return &echo.EchoReply{Message: "Hello " + in.Name + "  from hostname " + h}, nil
}

func (s *server) SayHelloLRO(ctx context.Context, in *echo.EchoRequest) (*lropb.Operation, error) {

	log.Println("Got SayHelloLRO --> ", in.Name)

	uid, _ := uuid.NewUUID()
	answer := &lropb.Operation{
		Name: uid.String(),
		Done: false,
	}

	work := workRequest{ID: uid.String(), Req: *in}
	workRequestMap[uid.String()] = work
	log.Printf("Work request queued %v", uid.String())

	return answer, nil
}

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	log.Printf("Handling grpc Check request for Service")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {

	flag.Parse()

	workRequestMap = make(map[string]workRequest)
	rand.Seed(time.Now().UnixNano())

	if *grpcport == "" {
		fmt.Fprintln(os.Stderr, "missing -grpcport flag (:50051)")
		flag.Usage()
		os.Exit(2)
	}

	ce, err := credentials.NewServerTLSFromFile("server_crt.pem", "server_key.pem")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
	sopts = append(sopts, grpc.Creds(ce))

	s := grpc.NewServer(sopts...)

	echo.RegisterEchoServerServer(s, &server{})
	healthpb.RegisterHealthServer(s, &healthServer{})
	lropb.RegisterOperationsServer(s, &operationsServer{})

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		s.GracefulStop()
		log.Printf("Inflight tasks %v ", workRequestMap)
		log.Printf("Shutting down..")
		close(idleConnsClosed)
	}()

	log.Printf("Starting gRPC sever on port %v", *grpcport)
	s.Serve(lis)
}
