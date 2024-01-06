package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/saragithubrit/golang-test/proto"
)

type server struct {
	mu    sync.Mutex
	users map[string]train.Ticket
}

func (s *server) PurchaseTicket(ctx context.Context, req *train.Ticket) (*train.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Assign a seat and store the user
	req.Section = "A" 
	s.users[req.User.Email] = *req

	return req, nil
}

func (s *server) GetReceipt(ctx context.Context, user *train.User) (*train.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, ok := s.users[user.Email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &ticket, nil
}

func (s *server) GetUsersBySection(section string, stream train.TrainService_GetUsersBySectionServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ticket := range s.users {
		if ticket.Section == section {
			if err := stream.Send(&ticket); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *server) RemoveUser(ctx context.Context, user *train.User) (*train.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, ok := s.users[user.Email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	delete(s.users, user.Email)
	return &ticket, nil
}

func (s *server) ModifySeat(ctx context.Context, user *train.User, section string) (*train.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticket, ok := s.users[user.Email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	ticket.Section = section
	s.users[user.Email] = ticket

	return &ticket, nil
}

func main() {
	s := &server{users: make(map[string]train.Ticket)}

	grpcServer := grpc.NewServer()
	train.RegisterTrainServiceServer(grpcServer, s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server is running on :50051")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	
	// Start HTTP server for gRPC-Gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = train.RegisterTrainServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

	log.Println("HTTP server is running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}

