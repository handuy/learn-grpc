package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	stream_proto "learn-grpc/server-stream/proto"
)

func main() {
	log.Println("Start client")

	connection, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}

	grpcClient := stream_proto.NewStreamServiceClient(connection)

	streamFromServer, err := grpcClient.Stream(context.Background(), &stream_proto.StreamRequest{
		Number: 20000,
	})
	if err != nil {
		log.Println(err)
		return
	}

	for {
		responseFromStream, err := streamFromServer.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			return
		}

		log.Println("response from stream", responseFromStream.Number)
	}

}
