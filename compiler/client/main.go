package main

import (
	"context"
	"fmt"
	"io"
	// "time"

	// "io"
	compiler_proto "learn-grpc/compiler/proto"
	"log"

	"google.golang.org/grpc"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
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

	ws := websocket.New(websocket.Config{})

	ws.OnConnection(func(c websocket.Connection) {
		log.Printf("[%s] Connected to server!", c.ID())
		c.OnMessage(func(code []byte) {
			log.Println("hello", string(code))

			connect, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
			}
			defer connect.Close()

			client := compiler_proto.NewCompileServiceClient(connect)
			maxSizeOption := grpc.MaxCallRecvMsgSize(64 * 10e6)
			responseFromServer, err := client.Compile(context.Background(), &compiler_proto.CompileRequest{
				Code: string(code),
			}, maxSizeOption)
			if err != nil {
				log.Println(err)
			}

			// result, err := responseFromServer.Recv()
			// log.Println("response from stream", result.Result)

			for {
				result, err := responseFromServer.Recv()
				if err == io.EOF {
					// ctx.Application().Logger().Errorf("stream: %v", err)
					break
				}

				if err != nil {
					// ctx.Application().Logger().Errorf("stream: %v", err)
					log.Println(err)
					return
				}

				log.Println("response from stream", result.Result)

				c.To(websocket.Broadcast).EmitMessage([]byte(result.Result))
			}
		})
	
	})

	app.Get("/my_endpoint", ws.Handler())

	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Post("/get-code", func(ctx iris.Context) {
		var req Code
		err := ctx.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(req.Code)
		ctx.JSON(req.Code)
	})

	app.Post("/compile", func(ctx iris.Context) {
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.ContentType("text/html")

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
				// ctx.Application().Logger().Errorf("stream: %v", err)
				break
			}

			if err != nil {
				// ctx.Application().Logger().Errorf("stream: %v", err)
				log.Println(err)
				return
			}

			log.Println("response from stream", result.Result)

			ctx.StreamWriter(func(io.Writer) bool {
				ctx.Writef("Message number %d<br>", result.Result)
				return true
			})

			
		}
	})

	app.Run(iris.Addr(":8087"), iris.WithoutServerError(iris.ErrServerClosed))
}
