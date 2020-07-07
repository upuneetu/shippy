// shippy-service-consignment/main.go

package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/upuneetu/shippy/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - dummy structure that is used for a datastore
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create - creating a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil

}

//service - implements all of the methods to satisfy the definition in protobuf
type service struct {
	repo repository
}

//CreateConsignment - creating a method on our service: takes context & request as argument - handled by grpc server
func (s *service) CreateConsignment(cntx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}
	// Setup grpc server

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	server := grpc.NewServer()

	// register out service with grpc server,
	// this ties our implementation into the auto generated interface code for our protobug definition
	pb.RegisterShippingServiceServer(server, &service{repo})
	//register the reflection on gRPC server
	reflection.Register(server)

	log.Println("Running on port:", port)
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
