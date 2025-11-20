package main

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "protobuf-grpc-reference/phase6/15/go/storagepb" // Assume generated package for storage.proto

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Storage state (in-memory, simulates file system interaction)
var ongoingUploads = make(map[string]int64) // file_id -> bytes_received
var uploadsMutex sync.Mutex

const STORAGE_ROOT = "./storage_data"
const HARDCODED_STORAGE_ADDRESS = "localhost:50052"

type storageServer struct {
	pb.UnimplementedFileServiceServer
}

// 15.3. Storage Logic: StreamFile (Bidirectional Streaming)
func (*storageServer) StreamFile(stream pb.FileService_StreamFileServer) error {
	log.Println("Storage: New Bidirectional stream established.")
	var fileId string
	var totalBytes int64
	var file *os.File

	// 1. Initial Authentication and File Setup
	// The client's first message should contain the file ID and token (in metadata or first chunk)
	firstChunk, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Stream started without initial chunk: %v", err)
	}

	fileId = firstChunk.GetFileId()
	if fileId == "" {
		return status.Errorf(codes.InvalidArgument, "File ID missing from first chunk.")
	}

	// --- Simplified Auth Check (15.2) ---
	// In a real system, you'd check context metadata for the token.
	// For simplicity, we just check if fileId is somewhat valid.
	// log.Printf("Storage: Checking Auth Token from metadata...")
	// if token != "EXPECTED_TOKEN" { return status.Errorf(codes.Unauthenticated, "Invalid token") }

	// --- 15.6. Resume/Append Logic ---
	filePath := filepath.Join(STORAGE_ROOT, fileId)
	// Open file in append mode. If it doesn't exist, it's created.
	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to open file for writing: %v", err)
	}
	defer file.Close()

	// Check current file size for resume offset
	stat, _ := file.Stat()
	totalBytes = stat.Size()
	log.Printf("Storage: Started writing file %s. Current size (for resume): %d bytes.", fileId, totalBytes)

	// Send initial status back to client for resume check
	statusMsg := &pb.FileChunk{
		FileId: fileId,
		Offset: totalBytes,
		// Data and ChunkIndex are usually empty for a status response
	}
	if err := stream.Send(statusMsg); err != nil {
		log.Printf("Storage: Failed to send initial status for resume: %v", err)
		return err
	}

	// 2. Main Upload Loop (Remaining Chunks)
	// We handle the first chunk already received, then loop for the rest.
	currentChunk := firstChunk
	for {
		// Only write data if the current offset is correct (handling the resume gap)
		if currentChunk.GetOffset() >= totalBytes {
			n, writeErr := file.Write(currentChunk.GetData())
			if writeErr != nil {
				return status.Errorf(codes.Internal, "File write failed: %v", writeErr)
			}
			totalBytes += int64(n)
		} else {
			// This chunk was already received (part of the resume gap)
			log.Printf("Storage: Skipping chunk %d (offset %d) as it's already received.", currentChunk.GetChunkIndex(), currentChunk.GetOffset())
		}

		// Send a periodic status update back (optional, for highly interactive feedback)
		if totalBytes%(1024*1024) == 0 { // Every 1MB
			updateMsg := &pb.FileChunk{
				FileId: fileId,
				Offset: totalBytes,
			}
			stream.Send(updateMsg)
		}

		// Get the next chunk
		currentChunk, err = stream.Recv()
		if err == io.EOF {
			log.Printf("Storage: Successfully finished upload for %s. Total bytes: %d", fileId, totalBytes)
			break
		}
		if err != nil {
			// Non-EOF error (e.g., network break).
			log.Printf("Storage: Network error during upload of %s: %v", fileId, err)
			return err // Returns the error to the client
		}
	}

	return nil
}

func main() {
	if _, err := os.Stat(STORAGE_ROOT); os.IsNotExist(err) {
		os.Mkdir(STORAGE_ROOT, 0755)
	}

	lis, err := net.Listen("tcp", HARDCODED_STORAGE_ADDRESS)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &storageServer{})
	log.Println("Storage Node (FileService) listening on :50052")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
