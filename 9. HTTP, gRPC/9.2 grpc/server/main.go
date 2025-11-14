package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	pb "tempgrpc/generated"
)

type server struct {
	pb.UnimplementedCalculatorServiceServer
}

func (s *server) Calculate(_ context.Context, req *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	var result float64
	switch req.Operation {
	case "add":
		result = req.A + req.B
	case "subtract":
		result = req.A - req.B
	case "multiply":
		result = req.A * req.B
	case "divide":
		if req.B == 0 {
			return &pb.CalculateResponse{Error: "Division by zero"}, nil
		}

		result = req.A / req.B
	default:
		return &pb.CalculateResponse{Error: "Unknown operation"}, nil
	}

	print(result)

	return &pb.CalculateResponse{Result: result}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	s.Serve(lis)
}
