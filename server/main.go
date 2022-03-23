package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/guil95/grpc-bidirectional-stream-example/pb/chat"
	"google.golang.org/grpc"
)

type server struct {
	chat.UnimplementedServiceServer
}

func (s *server) Chat(srv chat.Service_ChatServer) error {
	ctx := srv.Context()

	var chatMessage struct {
		message string
		name    string
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Server end")
			break
		default:
		}

		req, err := srv.Recv()
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Println(fmt.Sprintf("%sSay: %s", req.Name, req.Message))

		if chatMessage.name == "" && chatMessage.message == "" {
			fmt.Println("Your Name: ")
			in := bufio.NewReader(os.Stdin)
			chatMessage.name, err = in.ReadString('\n')
			fmt.Println("Your message: ")
			in = bufio.NewReader(os.Stdin)
			chatMessage.message, err = in.ReadString('\n')
		} else {
			in := bufio.NewReader(os.Stdin)
			chatMessage.message, err = in.ReadString('\n')
		}

		resp := chat.Response{Name: chatMessage.name, Message: chatMessage.message}
		if err := srv.Send(&resp); err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50001")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	chat.RegisterServiceServer(s, &server{})

	log.Println("Server running")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
