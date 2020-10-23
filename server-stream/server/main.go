package main

import (
	"fmt"
	stream_proto "learn-grpc/server-stream/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type streamServer struct {
	stream_proto.UnimplementedStreamServiceServer
}

func (s streamServer) Stream(req *stream_proto.StreamRequest, stream stream_proto.StreamService_StreamServer) error {
	numberToSolve := req.Number
	divisor := int32(2)

	for numberToSolve > 1 {
		if numberToSolve%divisor == 0 {
			err := stream.Send(&stream_proto.StreamResponse{
				Number: divisor,
			})

			log.Println("send to client", divisor)

			if err != nil {
				log.Println(err)
				return err
			}

			numberToSolve = numberToSolve / divisor

			time.Sleep(500 * time.Millisecond)
		} else {
			divisor++
		}
	}

	return nil
}

func main() {
	log.Println("Start server")

	listen, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		fmt.Println(err)
	}

	grpcServer := grpc.NewServer()
	stream_proto.RegisterStreamServiceServer(grpcServer, streamServer{})

	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
