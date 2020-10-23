package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"

	compiler_proto "learn-grpc/compiler/proto"

	"google.golang.org/grpc"
)

var bridge chan string

type compiler struct {
	compiler_proto.UnimplementedCompileServiceServer
}

func (c compiler) Compile(req *compiler_proto.CompileRequest, stream compiler_proto.CompileService_CompileServer) error {
	command := req.Code
	bridge = make(chan string)

	// command := "docker run compiler 'for (i = 0; i < 1000000; i++) {console.log(\"i bằng\", i)}'"
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	go func() {
		scanner := bufio.NewScanner(stdout)
		// scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()
			bridge <- m
		}
	
		// r := bufio.NewReader(stdout)
		// line, _, err := r.ReadLine()
		// if err != nil {
		// 	log.Println(err)
		// }
		// bridge <- string(line)
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
