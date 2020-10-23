package main

import (
	"context"
	"fmt"
	"learn-grpc/sum/sumpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client running")

	connection, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}

	client := sumpb.NewSumServiceClient(connection)

	resp, err := client.Sum(context.Background(), &sumpb.SumRequest{
		Num1: 10,
		Num2: 20,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}
