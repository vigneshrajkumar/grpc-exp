package main

import (
	"context"
	"log"
	"os"
	"petstore-client/proto/pb"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

func main() {
	log.Println("GRPC client")

	con, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	client := pb.NewPetstoreClient(con)

	if len(os.Args) != 2 {
		log.Println("improper usage")
	}
	switch os.Args[1] {
	case "catalogue":

		catalogue, err := client.GetCatalogue(context.TODO(), &emptypb.Empty{})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("------ CATALOGUE --------")
		for pet, price := range catalogue.Pets {
			log.Println(pet, ":", price)
		}

	default:
		log.Println("improper usage")
	}
}
