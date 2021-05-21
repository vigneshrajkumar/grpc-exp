package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"petstore-client/proto/pb"
	"sync"
	"time"

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

	case "get-list":
		stream, err := client.ListPets(context.TODO(), &emptypb.Empty{})
		if err != nil {
			log.Fatal(err)
		}

		for {
			animal, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			log.Println("recv: ", animal.Name, ": ", animal.Cost)
		}

	case "send-list":
		stream, err := client.OfferPets(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		catalogue := map[string]int32{
			"dog":  500,
			"cat":  200,
			"cow":  1000,
			"lion": 10000,
		}

		for pet, cost := range catalogue {
			time.Sleep(1 * time.Second)
			animal := &pb.Animal{Name: pet, Cost: cost}
			log.Println("sending: ", animal.Name, ": ", animal.Cost)
			if err := stream.Send(animal); err != nil {
				log.Fatal(err)
			}
		}
		offer, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("offer: ", offer.Message)

	case "negotiate":
		stream, err := client.Negotiate(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("negotiation started from client...")

		wg := sync.WaitGroup{}

		go func() {
			for {
				msg, err := stream.Recv()
				if err != nil {
					log.Fatal(err)
				}
				if msg.Contents == "exit" {
					wg.Done()
				}

				log.Println("store:", msg.Contents)
			}
		}()
		wg.Add(1)

		for {
			var message string
			fmt.Scanln(&message)
			if err := stream.Send(&pb.Message{Contents: message}); err != nil {
				log.Fatal(err)
			}
		}
		wg.Wait()

	default:
		log.Println("improper usage")
	}
}
