package main

import (
	"context"
	"fmt"
	"net"

	"learn-grpc/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	message := req.String()
	return &greetpb.GreetResponse{
		Resutl: message,
	}, nil
}

func main() {
	fmt.Println("Greet server running")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		fmt.Println(err)
	}

	grpcServer := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
