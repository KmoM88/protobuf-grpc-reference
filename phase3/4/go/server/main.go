package main

import (
	"context"
	"log"
	"net"

	// Import the protobuf module
	pb "protobuf-grpc-reference/phase3/4/go/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server structure implements the generated CalculatorServiceServer interface.
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

// 4.2. Implementation: Add Method
func (*server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Example: Check context for potential timeout (crucial in distributed systems)
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Errorf(codes.DeadlineExceeded, "Request timed out")
	}

	num1 := req.GetNum1()
	num2 := req.GetNum2()
	result := num1 + num2

	log.Printf("Go Server: Received Add request: %d + %d = %d", num1, num2, result)

	// Create and return the response message
	res := &pb.AddResponse{
		Result: result,
	}
	return res, nil
}

// Implementation: Subtract Method
func (*server) Subtract(ctx context.Context, req *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	result := num1 - num2

	log.Printf("Go Server: Received Subtract request: %d - %d = %d", num1, num2, result)

	res := &pb.SubtractResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	// Server listens on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fatal: failed to listen: %v", err)
	}

	// Create a new gRPC server instance (insecure for development)
	s := grpc.NewServer()

	// Register the implemented service
	pb.RegisterCalculatorServiceServer(s, &server{})

	log.Println("✅ Go Server started successfully on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fatal: failed to serve: %v", err)
	}
}
