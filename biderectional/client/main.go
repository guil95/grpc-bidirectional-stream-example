package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/guil95/grpc-streams-example/biderectional/pb/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}

	client := chat.NewServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	ctx := stream.Context()

	var chatMessage struct {
		message string
		name    string
	}

	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(fmt.Sprintf("\n%sSay: %s", resp.Name, resp.Message))
		}
	}()

	for {
		if chatMessage.name == "" && chatMessage.message == "" {
			fmt.Println("Your Name: ")
			in := bufio.NewReader(os.Stdin)
			chatMessage.name, err = in.ReadString('\n')
			fmt.Println("Your message: ")
		}

		in := bufio.NewReader(os.Stdin)
		chatMessage.message, err = in.ReadString('\n')

		if chatMessage.message == "/quit" {
			err := stream.CloseSend()
			<-ctx.Done()
			if err != nil {
				log.Fatal(err)
			}
		}

		req := chat.Request{Name: chatMessage.name, Message: chatMessage.message}

		err := stream.Send(&req)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
