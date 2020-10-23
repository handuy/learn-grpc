package main

import (
	"context"
	"fmt"
	"io"

	// "io"
	compiler_proto "learn-grpc/compiler/proto"
	"log"

	"google.golang.org/grpc"

	"github.com/kataras/iris"
)

type Code struct {
	Code string `json:"Code"`
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

	app.Post("/get-code", func(ctx iris.Context){
		var req Code
		err := ctx.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(req.Code)
		ctx.JSON(req.Code)
	})

	app.Post("/compile", func(ctx iris.Context){
		var req Code
		err := ctx.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(req.Code)
		
		connect, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
		}
		defer connect.Close()

		client := compiler_proto.NewCompileServiceClient(connect)
		maxSizeOption := grpc.MaxCallRecvMsgSize(64 * 10e6)
		responseFromServer, err := client.Compile(context.Background(), &compiler_proto.CompileRequest{
			Code: req.Code,
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
	})

	app.Run(iris.Addr(":8087"), iris.WithoutServerError(iris.ErrServerClosed))
}
