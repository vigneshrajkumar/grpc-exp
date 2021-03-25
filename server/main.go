package main

import (
	"context"
	"log"
	"net"
	"petstore-server/proto/pb"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

type petstore struct {
	pb.UnimplementedPetstoreServer
}

func (p *petstore) GetCatalogue(ctx context.Context, none *emptypb.Empty) (*pb.Catalogue, error) {
	log.Println("request for catalogue received")

	return &pb.Catalogue{
		Pets: map[string]int32{
			"dog": 500,
			"cat": 200,
			"cow": 1000,
		},
	}, nil
}

func main() {
	log.Println("GRPC server")

	lis, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()

	pb.RegisterPetstoreServer(server, &petstore{})

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
