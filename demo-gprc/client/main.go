package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"github.com/saragithubrit/golang-test/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	// 1. Purchase Ticket
	ticket, err := client.PurchaseTicket(context.TODO(), &train.Ticket{
		From:   "London",
		To:     "France",
		User:   &train.User{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		Price:  20.0,
		Section: "A", 
	})
	if err != nil {
		log.Fatalf("PurchaseTicket failed: %v", err)
	}
	fmt.Printf("Purchase successful. Receipt:\n%s\n", ticket)

	// 2. Get Receipt
	receipt, err := client.GetReceipt(context.TODO(), &train.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	})
	if err != nil {
		log.Fatalf("GetReceipt failed: %v", err)
	}
	fmt.Printf("Receipt:\n%s\n", receipt)

	// 3. Get Users By Section
	stream, err := client.GetUsersBySection(context.TODO(), &train.GetUsersBySectionRequest{Section: "A"})
	if err != nil {
		log.Fatalf("GetUsersBySection failed: %v", err)
	}

	fmt.Println("Users in Section A:")
	for {
		user, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", user)
	}

	// 4. Remove User
	removedUser, err := client.RemoveUser(context.TODO(), &train.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	})
	if err != nil {
		log.Fatalf("RemoveUser failed: %v", err)
	}
	fmt.Printf("Removed user:\n%s\n", removedUser)

	// 5. Modify Seat
	modifiedTicket, err := client.ModifySeat(context.TODO(), &train.Mod
