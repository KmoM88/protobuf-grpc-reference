# Solution: Exercise (Image Uploader)
This solution demonstrates Client Streaming with a Python Server and a Go Client.

## 1. Protobuf Definition
`proto/uploader.proto`

```protobuf
syntax = "proto3";

package uploader;

option go_package = "./uploaderpb";

service UploaderService {
  // Defines Client Streaming: The request is a stream of ImageChunk messages.
  rpc UploadImage (stream ImageChunk) returns (UploadSummary);
}

message ImageChunk {
  // Raw byte data for a file chunk
  bytes data = 1;
}

message UploadSummary {
  int32 total_chunks = 1;
  int64 total_bytes = 2;
}
```
## 2. Python Server (`python/server.py`)
This server counts the chunks and the total size by iterating over the client's request stream.

```python
import grpc
import time
from concurrent import futures

# Import generated modules
import uploader_pb2 as pb
import uploader_pb2_grpc as rpc

class UploaderServicer(rpc.UploaderServiceServicer):
    
    # 6.5. Implementation: Client Streaming Method
    # The first argument is the request_iterator
    def UploadImage(self, request_iterator, context):
        total_bytes = 0
        total_chunks = 0
        
        print("Server: Starting upload process...")

        # Iterate over the incoming stream (Python handles EOF automatically)
        for chunk in request_iterator:
            chunk_size = len(chunk.data)
            total_bytes += chunk_size
            total_chunks += 1
            print(f"Server: Received chunk #{total_chunks} ({chunk_size} bytes)")
            
        print(f"Server: Finished receiving stream.")
        
        # Return the single final response (UploadSummary)
        return pb.UploadSummary(
            total_chunks=total_chunks, 
            total_bytes=total_bytes
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    rpc.add_UploaderServiceServicer_to_server(UploaderServicer(), server)
    
    server.add_insecure_port('[::]:50051')
    server.start()
    print("✅ Python Server listening on :50051")
    
    try:
        while True:
            time.sleep(86400) 
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
```
## 3. Go Client (go/client/main.go)
This client generates random data, streams it to the Python server, and receives the final summary.

```go
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
```

