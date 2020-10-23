package main

import (
	"context"
	"fmt"
	todo_proto "learn-grpc/todo-list/proto"
	"log"
	"net"

	"github.com/rs/xid"
	"google.golang.org/grpc"

	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type server struct {
	DB    *sql.DB
	todo_proto.UnimplementedTodoServiceServer
}

func (s server) CreateTodo(ctx context.Context, req *todo_proto.CreateTodoRequest) (*todo_proto.CreateTodoResponse, error) {
	var newTodoID = xid.New().String()
	sqlStatement := `
	INSERT INTO todo (id, title)
	VALUES ($1, $2)`
	_, err := s.DB.Exec(sqlStatement, newTodoID, req.Title)
	if err != nil {
		panic(err)
	}

	return &todo_proto.CreateTodoResponse{
		Id: newTodoID,
	}, nil
}

func main() {
	log.Println("Start server")

	listen, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		fmt.Println(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	todoService := server{}
	todoService.DB = db

	grpcServer := grpc.NewServer()
	todo_proto.RegisterTodoServiceServer(grpcServer, todoService)

	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
	}
}
