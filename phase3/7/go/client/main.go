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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
