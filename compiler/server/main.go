package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	compiler_proto "learn-grpc/compiler/proto"

	"google.golang.org/grpc"
)

var bridge chan string

type compiler struct {
	compiler_proto.UnimplementedCompileServiceServer
}

func (c compiler) Compile(req *compiler_proto.CompileRequest, stream compiler_proto.CompileService_CompileServer) error {
	code := req.Code
	file, err := os.Create("index.js")
	if err != nil {
		fmt.Println(err)
	} else {
		file.WriteString(code)
	}
	command := "docker run --rm -v $PWD/index.js:/index.js node:13-alpine node index.js"

	bridge = make(chan string)

	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			bridge <- m
		}

		close(bridge)
	}()

	for elem := range bridge {
		err := stream.Send(&compiler_proto.CompileResponse{
			Result: elem,
		})

		if err != nil {
			log.Println(err)
			return err
		}

		// select {
		// case <-bridge:
		//     err := stream.Send(&compiler_proto.CompileResponse{
		// 		Result: elem,
		// 	})

		// 	if err != nil {
		// 		log.Println(err)
		// 		return err
		// 	}
		// default:
		//     break;
		// }
	}

	cmd.Wait()
	return nil
}

func main() {
	log.Println("Start server")

	listen, err := net.Listen("tcp", "0.0.0.0:9001")
	if err != nil {
		log.Println(err)
	}

	grpcServer := grpc.NewServer()
	compiler_proto.RegisterCompileServiceServer(grpcServer, compiler{})

	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
