# Solution: Chapter 7 Exercise (Simple Chat Room)
This solution uses a Go Server and two clients (one Go and one Python) to demonstrate cross-language Bidirectional Streaming.

## 1. Protobuf Definition
`proto/chat.proto`

```protobuf
syntax = "proto3";

package chat;

option go_package = "./chatpb";

service ChatService {
  // Bidirectional Streaming: User sends messages, and receives messages from others.
  rpc JoinChat (stream ChatMessage) returns (stream ChatMessage);
}

message ChatMessage {
  string user_id = 1;
  string text = 2;
  int64 timestamp = 3; // For time-stamping messages
}
```
## 2. Go Server (`go/server/main.go`)
The server uses a global map (`clients`) to store an output channel for every connected client. Messages are broadcasted to all channels.

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "./chatpb" // Adjust path
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
	clients = make(map[string]*clientChannel)
	mu      sync.Mutex
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
```
## 3. Go Client (`go/client/main.go`)
This client uses two goroutines: one for reading (`recv`) and one for writing (`send`).

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "protobuf-grpc-reference/phase3/7/go/chatpb" // Adjust module path
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Unique ID for this client instance
const CLIENT_ID = "GoClient" 

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.JoinChat(ctx)
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}
	
	fmt.Printf("\n--- Go Client (%s) Joined Chat ---\n", CLIENT_ID)
	fmt.Println("Type messages and press Enter. Type 'exit' to quit.")

	// --- 1. Receiver Goroutine ---
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\nServer closed the stream.")
				cancel()
				return
			}
			if err != nil {
				if status.Code(err) == codes.Canceled {
					return // Stream cancelled by client
				}
				log.Printf("Receiver Error: %v", err)
				cancel()
				return
			}
			
			// Display the received message
			ts := time.Unix(msg.GetTimestamp(), 0).Format("15:04:05")
			fmt.Printf("[%s] %s: %s\n", ts, msg.GetUserId(), msg.GetText())
		}
	}()

	// --- 2. Sender Loop (Main Thread) ---
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		
		if strings.ToLower(text) == "exit" {
			// Client closes its side of the stream.
			stream.CloseSend()
			return
		}

		chatMsg := &pb.ChatMessage{
			UserId: CLIENT_ID,
			Text:   text,
		}

		if err := stream.Send(chatMsg); err != nil {
			log.Printf("Sender Error: %v", err)
			break
		}
	}
}
```
## 4. Python Client (`python/client.py`)
This client uses threads to manage the sending and receiving processes concurrently.

```python
import grpc
import time
import threading
import sys
from datetime import datetime

# Import generated modules
import chat_pb2 as pb
import chat_pb2_grpc as rpc

# Unique ID for this client instance
CLIENT_ID = "PythonClient"

def generate_requests():
    """Generator for the Sender Thread to stream messages to the server."""
    print("Type messages and press Enter. Type 'exit' to quit.")
    
    while True:
        try:
            # Read user input from the console
            text = input(f"[{CLIENT_ID}] > ")
        except EOFError:
            # Handle Ctrl+D or end of file
            break
        except KeyboardInterrupt:
            # Handle Ctrl+C
            break

        if text.lower() == "exit":
            break

        chat_msg = pb.ChatMessage(user_id=CLIENT_ID, text=text)
        yield chat_msg
        
        # Ensures the Python runtime is responsive
        time.sleep(0.1)

def receive_responses(stream_reader):
    """Function run in a separate thread to receive responses continuously."""
    try:
        # Iterate over the response stream (server's stream)
        for response in stream_reader:
            ts = datetime.fromtimestamp(response.timestamp).strftime('%H:%M:%S')
            
            # Print the received message above the input prompt
            sys.stdout.write(f"\r[{ts}] {response.user_id}: {response.text}\n[{CLIENT_ID}] > ")
            sys.stdout.flush()

    except grpc.RpcError as e:
        if e.code() == grpc.StatusCode.CANCELLED:
            print("\nClient Receiver: Stream cancelled by client.")
        else:
            print(f"\nClient Receiver: An RPC error occurred: {e.details()}")
    except Exception as e:
        print(f"\nClient Receiver: Unexpected error: {e}")
    finally:
        # Stop the main thread by exiting the process
        print("\nReceiver thread finished. Exiting...")
        os._exit(0) # Force exit if main thread is blocked

def run():
    # 1. Establish connection
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = rpc.ChatServiceStub(channel)
        
        print(f"\n--- Python Client ({CLIENT_ID}) Joined Chat ---")

        # 2. Initiate the Bidirectional RPC call
        # The stream_handler is an iterable used by the receiver thread.
        # The request generator is consumed by the client runtime for sending.
        stream_handler = stub.JoinChat(generate_requests())

        # 3. Start the Receiver Thread
        receiver_thread = threading.Thread(target=receive_responses, args=(stream_handler,))
        receiver_thread.daemon = True # Allows thread to exit when main thread does
        receiver_thread.start()

        # 4. Wait for the receiver thread to run indefinitely (user interaction loop)
        while receiver_thread.is_alive():
            time.sleep(1)

if __name__ == '__main__':
    run()
```
## 5. Execution Commands Summary
Here are the commands to set up, generate the code, and run the clients and server.

### 1. Generate Code
Run these commands from the root of your project.

#### Go Generation (Server and Go Client)
```bash
protoc --proto_path=proto \
       --go_out=go/chatpb --go_opt=paths=source_relative \
       --go-grpc_out=go/chatpb --go-grpc_opt=paths=source_relative \
       proto/chat.proto
```
#### Python Generation (Python Client)
```bash
python -m grpc_tools.protoc --proto_path=proto \
                           --python_out=python \
                           --pyi_out=python \
                           --grpc_python_out=python \
                           proto/chat.proto
```
### 2. Run Server and Clients
#### Start Go Server (Terminal 1)
```bash
cd go/server
go run main.go
# Output: ✅ Go Chat Server listening on :50051 (Broadcasting enabled)
```
#### Run Go Client (Terminal 2)
```bash
cd go/client
go run main.go
# (Client will prompt for input)
```
#### Run Python Client (Terminal 3)
```bash
cd python
python client.py
# (Client will prompt for input)
```
Testing: Messages typed in Terminal 2 (Go Client) should appear in Terminal 3 (Python Client) and vice-versa, demonstrating the server's broadcast logic and Bidirectional communication.