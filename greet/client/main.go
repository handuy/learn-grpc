package main

import (
	"context"
	"fmt"

	"learn-grpc/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Client running")

	connect, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer connect.Close()

	client := greetpb.NewGreetServiceClient(connect)

	req := &greetpb.GreetRequest{
		Greet: &greetpb.Greeting{
			FirstName: "123",
			LastName:  "456",
		},
	}

	res, err := client.Greet(context.Background(), req)
	fmt.Println("From server", res)
}
