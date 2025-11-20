package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "protobuf-grpc-reference/phase6/15/go/storagepb" // Assume generated package for storage.proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Simplified in-memory metadata store
var metadataStore = make(map[string]pb.FileMetadata)
var mu sync.Mutex

// HARDCODED_STORAGE_ADDRESS must match the Storage Node's address
const HARDCODED_STORAGE_ADDRESS = "localhost:50052"

type masterServer struct {
	pb.UnimplementedMetaServiceServer
}

// 15.2. Master Logic: RequestUpload (Unary)
func (*masterServer) RequestUpload(ctx context.Context, meta *pb.FileMetadata) (*pb.FileMetadata, error) {
	if meta.Filename == "" || meta.SizeBytes <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Filename and size are required.")
	}

	// 1. Generate unique file ID and simple auth token
	fileId := fmt.Sprintf("file_%d", time.Now().UnixNano())
	authToken := fmt.Sprintf("TOKEN-%d", time.Now().Unix()) // Simple token generation

	// 2. Prepare and store metadata
	newMeta := *meta
	newMeta.FileId = fileId
	newMeta.StorageNodeAddress = HARDCODED_STORAGE_ADDRESS
	newMeta.AuthToken = authToken

	mu.Lock()
	metadataStore[fileId] = newMeta
	mu.Unlock()

	log.Printf("Master: Registered new upload request for '%s'. Assigned ID: %s", meta.Filename, fileId)
	return &newMeta, nil
}

// ... GetFileLocation and GetSystemStatus implementation goes here ...
// Use GetSystemStatus for Server Streaming example (return all metadata in the map)

func main() {
	// ... gRPC Server setup (same as previous chapters) ...
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMetaServiceServer(s, &masterServer{})
	log.Println("Master Node (MetaService) listening on :50051")

	// Note: mTLS setup (15.5) will replace grpc.NewServer()
	// with credentials loading in the final step.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
