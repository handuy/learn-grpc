package main

import (
	"context"
	"fmt"
	"learn-grpc/sum/sumpb"
	"net"

	"google.golang.org/grpc"
)

type sumServer struct {
	sumpb.UnimplementedSumServiceServer
}

func (s sumServer) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	firstNum := req.GetNum1()
	secondNum := req.GetNum2()
	result := firstNum + secondNum
	return &sumpb.SumResponse{
		Resutl: result,
	}, nil
}

func main() {
	fmt.Println("Server running")

	listen, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		fmt.Println(err)
	}

	grpcServer := grpc.NewServer()
	sumpb.RegisterSumServiceServer(grpcServer, sumServer{})

	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
