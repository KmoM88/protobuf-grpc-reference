package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	pb "protobuf-grpc-reference/phase3/5/go/tickerpb" // Assuming tickerpb is generated here

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
