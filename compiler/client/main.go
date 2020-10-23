package main

import (
	"context"
	"fmt"
	"io"

	// "io"
	compiler_proto "learn-grpc/compiler/proto"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Start client")

	connect, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer connect.Close()

	client := compiler_proto.NewCompileServiceClient(connect)
	maxSizeOption := grpc.MaxCallRecvMsgSize(64 * 10e6)
	responseFromServer, err := client.Compile(context.Background(), &compiler_proto.CompileRequest{
		Code: "123",
	}, maxSizeOption)
	if err != nil {
		log.Println(err)
	}

	// result, err := responseFromServer.Recv()
	// log.Println("response from stream", result.Result)

	for {
		result, err := responseFromServer.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return
		}

		log.Println("response from stream", result.Result)
	}
}
