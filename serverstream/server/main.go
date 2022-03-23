package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/guil95/grpc-streams-example/serverstream/pb/products"
	"google.golang.org/grpc"
)

type server struct {
	products.UnsafeProductServiceServer
}

type ProductList struct {
	ProductList []Product `json:"products"`
}

type Product struct {
	Description string `json:"description"`
	Value       int64  `json:"value"`
}

func (s *server) ListProducts(req *emptypb.Empty, srv products.ProductService_ListProductsServer) error {
	jsonFile, err := os.Open("./serverstream/server/products.json")
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var productList ProductList

	err = json.Unmarshal(byteValue, &productList)
	if err != nil {
		return err
	}

	for i, p := range productList.ProductList {
		log.Println(fmt.Sprintf("Chave: %d produto: %s", i, p.Description))
		time.Sleep(time.Second * 2)
		srv.Send(&products.Product{Description: p.Description, Value: p.Value})
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50002")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	products.RegisterProductServiceServer(s, &server{})

	log.Println("Server running")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
