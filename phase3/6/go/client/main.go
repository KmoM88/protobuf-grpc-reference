package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	pb "protobuf-grpc-reference/phase3/6/go/uploaderpb" // Adjust module path

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const CHUNK_SIZE = 1024 // 1 KB per chunk

func createRandomChunk() []byte {
	b := make([]byte, CHUNK_SIZE)
	rand.Read(b)
	return b
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// 1. Establish connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// 2. Create client stub
	client := pb.NewUploaderServiceClient(conn)

	// Create context with a timeout for the entire RPC
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3. Initiate the streaming RPC call
	// The client receives the stream object to send messages on.
	stream, err := client.UploadImage(ctx)
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}

	totalChunks := 5 // We will send 5 chunks

	log.Printf("Client: Starting to stream %d chunks...", totalChunks)

	// 4. Client Streaming Logic: Send chunks in a loop
	for i := 1; i <= totalChunks; i++ {
		chunkData := createRandomChunk()

		chunkMsg := &pb.ImageChunk{
			Data: chunkData,
		}

		// Send the message to the server
		if err := stream.Send(chunkMsg); err != nil {
			if err == io.EOF {
				break // Server already closed the stream (unlikely here)
			}
			log.Fatalf("Error sending chunk %d: %v", i, err)
		}
		log.Printf("Client: Sent chunk #%d (%d bytes)", i, len(chunkData))
		time.Sleep(50 * time.Millisecond) // Small delay
	}

	// 5. Close the stream and receive the single response
	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving summary from server: %v", err)
	}

	// 6. Print the result
	fmt.Println("\n--- Upload Summary Received ---")
	fmt.Printf("Total Chunks Sent: %d (Expected: %d)\n", summary.GetTotalChunks(), totalChunks)
	fmt.Printf("Total Bytes Uploaded: %d bytes (Expected: %d bytes)\n", summary.GetTotalBytes(), totalChunks*CHUNK_SIZE)
	fmt.Println("-----------------------------")
}
