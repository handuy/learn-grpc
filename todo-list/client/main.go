package main

import (
	"context"
	"log"

	todo_proto "learn-grpc/todo-list/proto"

	"google.golang.org/grpc"
	"github.com/kataras/iris"
)

type formData struct {
	todo string `form:"todo"`
}

func main() {
	log.Println("Start client")
	app := iris.New()
	// // Đăng ký thư mục chứa HTML
	tmpl := iris.HTML("./view", ".html").Reload(true)
	app.RegisterView(tmpl)

	app.Get("/", func(ctx iris.Context){
		ctx.View("index.html")
	})

	app.Post("/create-todo", func(ctx iris.Context) {
		newTodo := formData{}
		err := ctx.ReadForm(&newTodo)

		log.Println("newTodo", newTodo)

		connection, err := grpc.Dial("localhost:8088", grpc.WithInsecure())
		if err != nil {
			log.Println(err)
		}

		grpcClient := todo_proto.NewTodoServiceClient(connection)
		rsp, err := grpcClient.CreateTodo(context.Background(), &todo_proto.CreateTodoRequest{
			Title: newTodo.todo,
		})
		if err != nil {
			log.Println(err)
		}

		log.Println("Server response", rsp.Id)

		ctx.StatusCode(200)
		ctx.JSON(rsp.Id)
	})

	app.Run(iris.Addr(":8087"), iris.WithoutServerError(iris.ErrServerClosed))
}
