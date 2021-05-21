package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"petstore-server/proto/pb"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (p *petstore) ListPets(none *emptypb.Empty, stream pb.Petstore_ListPetsServer) error {
	log.Println("request for streaming pets received")
	catalogue := map[string]int32{
		"dog": 500,
		"cat": 200,
		"cow": 1000,
	}
	for pet, cost := range catalogue {
		time.Sleep(1 * time.Second)
		if err := stream.Send(&pb.Animal{Cost: cost, Name: pet}); err != nil {
			return err
		}
	}
	return nil
}

func (p *petstore) OfferPets(stream pb.Petstore_OfferPetsServer) error {
	log.Println("request for offering pets received")
	rejectedAnimals := make([]string, 0)
	for {
		animal, err := stream.Recv()
		if err == io.EOF {
			msg := "we accept it all"
			if len(rejectedAnimals) > 0 {
				msg += " but " + strings.Join(rejectedAnimals, ", ")
			}
			return stream.SendAndClose(&pb.Offer{Message: msg})
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println("offer for ", animal.Name, " @", animal.Cost)
		if animal.Name == "lion" {
			rejectedAnimals = append(rejectedAnimals, animal.Name)
		}
	}
	return nil
}

func (p *petstore) Negotiate(stream pb.Petstore_NegotiateServer) error {
	log.Println("negtation called from client")

	for {
		message, err := stream.Recv()
		if err != nil {
			return err
		}
		if message.Contents == "exit" {
			if err := stream.Send(&pb.Message{Contents: "exit"}); err != nil {
				return err
			}
		}
		log.Println("client: ", message.Contents)
		time.Sleep(1 * time.Second)

		var resp string
		fmt.Scanln(&resp)
		if err := stream.Send(&pb.Message{Contents: resp}); err != nil {
			log.Fatal(err)
		}

	}
	return nil
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
