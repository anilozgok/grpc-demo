package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/anilozgok/grpc-demo/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedDemoServiceServer
}

func (s *server) UnaryCall(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Message: "Unary response: " + req.Message}, nil
}

func (s *server) ClientStreamingCall(stream pb.DemoService_ClientStreamingCallServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Response{Message: "Client streaming response: " + fmt.Sprint(messages)})
		}
		if err != nil {
			return err
		}
		messages = append(messages, req.Message)
	}
}

func (s *server) ServerStreamingCall(req *pb.Request, stream pb.DemoService_ServerStreamingCallServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.Response{Message: fmt.Sprintf("Server streaming response %d: %s", i, req.Message)}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) BidirectionalStreamingCall(stream pb.DemoService_BidirectionalStreamingCallServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.Response{Message: "Bidirectional response: " + req.Message}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDemoServiceServer(s, &server{})

	log.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
