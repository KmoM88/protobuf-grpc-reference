# Solution: Exercise (Stock Ticker)

## 1. Protobuf Definition
`proto/ticker.proto`
```protobuf
syntax = "proto3";

package ticker;

option go_package = "./tickerpb";

service TickerService {
  // Defines Server Streaming: The response is a stream of StockPrice messages.
  rpc GetStockPrices (TickerRequest) returns (stream StockPrice);
}

message TickerRequest {
  // Symbol to monitor, e.g., "TSLA"
  string symbol = 1;
}

message StockPrice {
  string symbol = 1;
  double price = 2;
  int64 timestamp = 3;
}
```
## 2. Go Server (go/server/main.go)
This server implements the streaming logic.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "protobuf-grpc-reference/phase3/5/go/tickerpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedTickerServiceServer
}

// Helper to generate a random price based on a base seed
func generateRandomPrice(symbol string) float64 {
	// Base price seed (e.g., TSLA might be around 250)
	base := 250.0 
	
	// Add a random fluctuation (0 to 10)
	fluctuation := rand.Float64() * 10 
	
	return base + fluctuation
}

// 5.3. Implementation: Server Streaming Method
func (*server) GetStockPrices(req *pb.TickerRequest, stream pb.TickerService_GetStockPricesServer) error {
	symbol := req.GetSymbol()
	log.Printf("Server: Starting stream for symbol: %s", symbol)

	// Stream 10 prices
	for i := 0; i < 10; i++ {
		// Check for client cancellation
		if stream.Context().Err() == context.Canceled {
			log.Printf("Server: Client cancelled stream for %s.", symbol)
			return status.Error(codes.Canceled, "Stream cancelled by client")
		}
		
		priceValue := generateRandomPrice(symbol)

		priceMsg := &pb.StockPrice{
			Symbol:    symbol,
			Price:     priceValue,
			Timestamp: time.Now().Unix(),
		}

		// 1. Send the message
		if err := stream.Send(priceMsg); err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}

		log.Printf("Server: Sent price #%d: %.2f", i+1, priceValue)
		
		// 2. Pause for 1 second
		time.Sleep(1 * time.Second)
	}

	log.Printf("Server: Stream finished for symbol: %s", symbol)
	return nil // Returning nil closes the stream successfully
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano()) 

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTickerServiceServer(s, &server{})
	
	log.Println("Go Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```
## 3. Python Client (python/client.py)
This client connects to the Go server and reads the streaming response.

```python
import grpc
import time

# Import generated modules (adjust import based on your generation path)
import ticker_pb2 as pb
import ticker_pb2_grpc as rpc


def run():
    # Connect to the Go server
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = rpc.TickerServiceStub(channel)
        
        request = pb.TickerRequest(symbol="TSLA")
        
        print("--- Initiating Server Stream: TSLA Prices ---")
        
        try:
            # Call the Server Streaming RPC. This returns an iterable response object.
            price_stream = stub.GetStockPrices(request)
            
            # Iterate over the stream response. Each iteration blocks until a message arrives.
            for price in price_stream:
                print(f"💰 Received: Symbol={price.symbol} | Price={price.price:.2f} | Time={time.ctime(price.timestamp)}")
            
            print("--- Stream Complete ---")

        except grpc.RpcError as e:
            if e.code() == grpc.StatusCode.UNAVAILABLE:
                print("Error: Server is unavailable.")
            else:
                print(f"An RPC error occurred: {e.details()}")

if __name__ == '__main__':
    run()
```
## 4. Execution Order
- Generate Go and Python Code (using the commands learned in Chapter 3/4).
- If needed clean go cache and modcache with `go clean -cache -modcache`.
- Start Go Server (`cd go/server && go run main.go`).
- Run Python Client (`cd python && python client.py`).

The client will print 10 prices, one every second, demonstrating the continuous streaming of data.