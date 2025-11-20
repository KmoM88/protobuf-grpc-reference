package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "protobuf-grpc-reference/phase3/7/go/chatpb" // Adjust path

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// clientChannel is used to send messages to a specific client goroutine.
type clientChannel struct {
	stream pb.ChatService_JoinChatServer
	id     string
}

// Global state to manage connected clients and broadcast messages.
var (
	clients      = make(map[string]*clientChannel)
	mu           sync.Mutex
	messageQueue = make(chan *pb.ChatMessage, 100) // Channel for broadcast messages
)

type server struct {
	pb.UnimplementedChatServiceServer
}

// Start a goroutine to continuously broadcast messages from the queue to all connected clients.
func init() {
	go broadcaster()
}

func broadcaster() {
	for msg := range messageQueue {
		mu.Lock()
		// Iterate over all connected client channels and send the message
		for _, client := range clients {
			if err := client.stream.Send(msg); err != nil {
				log.Printf("Broadcaster: Error sending message to %s: %v", client.id, err)
				// Handle disconnection: client will be removed by the main stream goroutine
			}
		}
		mu.Unlock()
	}
}

// 7.2. Implementation: Bidirectional Streaming Method (Go)
func (*server) JoinChat(stream pb.ChatService_JoinChatServer) error {
	// Assign a unique ID to the new client
	clientName := fmt.Sprintf("User_%d", time.Now().UnixNano())
	log.Printf("Server: New client connected with ID: %s", clientName)

	client := &clientChannel{
		stream: stream,
		id:     clientName,
	}

	// Add client to the global map
	mu.Lock()
	clients[clientName] = client
	mu.Unlock()

	// --- 1. Reader Goroutine (Handles incoming client messages) ---
	// This goroutine blocks on Recv() and pushes messages to the global queue.
	go func() {
		for {
			req, err := stream.Recv()

			if err == io.EOF {
				log.Printf("Server: Client %s finished sending.", clientName)
				return // End of client's stream
			}
			if err != nil {
				// Handle client disconnection or error
				log.Printf("Server: Error reading from client %s: %v", clientName, err)
				return
			}

			// Augment message with Server info and push to broadcast queue
			req.Timestamp = time.Now().Unix()
			messageQueue <- req
			log.Printf("Server: Received from %s: %s", req.GetUserId(), req.GetText())
		}
	}()

	// --- 2. Keep the main stream function alive ---
	// Wait for the stream context to be cancelled (client disconnects or error)
	<-stream.Context().Done()

	// --- Clean up ---
	mu.Lock()
	delete(clients, clientName)
	mu.Unlock()
	log.Printf("Server: Client %s disconnected and removed.", clientName)

	return status.Error(codes.Canceled, "Stream closed by client or server error")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fatal: failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})

	log.Println("✅ Go Chat Server listening on :50051 (Broadcasting enabled)")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fatal: failed to serve: %v", err)
	}
}
