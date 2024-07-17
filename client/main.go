package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/anilozgok/grpc-demo/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDemoServiceClient(conn)

	// Unary call
	unaryResp, err := c.UnaryCall(context.Background(), &pb.Request{Message: "Hello Unary"})
	if err != nil {
		log.Fatalf("Unary call failed: %v", err)
	}
	fmt.Println("Unary response:", unaryResp.Message)

	// Client streaming call
	clientStream, err := c.ClientStreamingCall(context.Background())
	if err != nil {
		log.Fatalf("Client streaming call failed: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := clientStream.Send(&pb.Request{Message: fmt.Sprintf("Message %d", i)}); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}
	clientStreamResp, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}
	fmt.Println("Client streaming response:", clientStreamResp.Message)

	// Server streaming call
	serverStream, err := c.ServerStreamingCall(context.Background(), &pb.Request{Message: "Hello Server Stream"})
	if err != nil {
		log.Fatalf("Server streaming call failed: %v", err)
	}
	for {
		resp, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive message: %v", err)
		}
		fmt.Println("Server streaming response:", resp.Message)
	}

	// Bidirectional streaming call
	bidiStream, err := c.BidirectionalStreamingCall(context.Background())
	if err != nil {
		log.Fatalf("Bidirectional streaming call failed: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			resp, err := bidiStream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message: %v", err)
			}
			fmt.Println("Bidirectional response:", resp.Message)
		}
	}()
	for i := 0; i < 5; i++ {
		if err := bidiStream.Send(&pb.Request{Message: fmt.Sprintf("Message %d", i)}); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		time.Sleep(time.Second)
	}
	bidiStream.CloseSend()
	<-waitc
}
