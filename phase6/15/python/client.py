import grpc
import time
import os
import random
from typing import Iterator

# Import generated modules (Adjust path based on your generation)
import storage_pb2 as pb
import storage_pb2_grpc as rpc

# --- Configuration ---
MASTER_ADDRESS = 'localhost:50051'
CHUNK_SIZE = 1024 * 64  # 64 KB chunks

# --- Simplified mTLS Credentials (15.5) ---
# FIX: Since we are running WITHOUT mTLS, this function returns None.
def get_credentials():
    return None


def file_chunk_generator(file_path: str, file_id: str, offset: int) -> Iterator[pb.FileChunk]:
    """Reads a file from the specified offset and yields FileChunk messages."""
    
    with open(file_path, 'rb') as f:
        f.seek(offset) # Start reading from the resume offset
        chunk_index = offset // CHUNK_SIZE
        
        while True:
            data = f.read(CHUNK_SIZE)
            if not data:
                break # EOF

            yield pb.FileChunk(
                file_id=file_id,
                chunk_index=chunk_index,
                data=data,
                offset=f.tell() - len(data) # Current chunk's start offset
            )
            chunk_index += 1
            # Simulate a network interruption every 5th chunk for 15.6
            if chunk_index % 5 == 0:
                print("Client: !!! SIMULATING NETWORK BREAK !!!")
                # In a real app, this would be a disconnection error. 
                # Here, we can simulate by yielding a problematic chunk or exiting the generator.
                # To simulate resumption, we need to exit and reconnect (handled in the main run function).
                # For simplicity in the generator, we just show a message.
        
def run_upload(local_file_path: str):
    file_size = os.path.getsize(local_file_path)
    file_name = os.path.basename(local_file_path)

    # 1. Connect to Master Node and get Metadata
    with grpc.insecure_channel(MASTER_ADDRESS, get_credentials()) as master_channel:
        master_stub = rpc.MetaServiceStub(master_channel)
        
        initial_meta = pb.FileMetadata(filename=file_name, size_bytes=file_size)
        file_meta = master_stub.RequestUpload(initial_meta)
        
        file_id = file_meta.file_id
        storage_addr = file_meta.storage_node_address
        auth_token = file_meta.auth_token
        
        print(f"Master granted ID: {file_id}, Storage: {storage_addr}, Token: {auth_token}")
        
    # 2. Connect to Storage Node (with token in metadata if needed)
    with grpc.insecure_channel(storage_addr, get_credentials()) as storage_channel:
        storage_stub = rpc.FileServiceStub(storage_channel)
        
        # Set initial offset to 0
        current_offset = 0
        
        while current_offset < file_size:
            print(f"\n--- Starting Stream (Offset: {current_offset} bytes) ---")
            
            try:
                # Prepare metadata for the storage node (including Auth Token)
                # In a real implementation, you'd send the token here.
                metadata = (('authorization', f'Bearer {auth_token}'),)
                
                # Start Bidirectional Stream
                stream = storage_stub.StreamFile(file_chunk_generator(local_file_path, file_id, current_offset), metadata=metadata)
                
                # First response from the server is the current file offset for resumption (15.6)
                try:
                    initial_status = next(stream)
                    current_offset = initial_status.offset
                    print(f"Server reported existing bytes: {current_offset}. Resuming from here.")
                except StopIteration:
                    print("No initial status received.")
                    
                # Main receiver loop for periodic status updates or control messages
                for status_msg in stream:
                    if status_msg.offset > 0:
                        print(f"Server Update: {status_msg.offset} bytes received.")
                
                # If the generator finishes and the receiver loop completes, the file is uploaded.
                current_offset = file_size # Mark as complete
                print("--- Upload finished successfully. ---")
                
            except grpc.RpcError as e:
                if e.code() == grpc.StatusCode.UNAVAILABLE or e.code() == grpc.StatusCode.DEADLINE_EXCEEDED:
                    print(f"Network interruption or timeout: {e.details()}. Retrying in 2 seconds...")
                    time.sleep(2)
                    # Loop continues, and file_chunk_generator will pick up from current_offset
                else:
                    print(f"Fatal RPC Error: {e.details()}")
                    break

def main():
    # Create a dummy file for testing
    dummy_file_path = "test_file.bin"
    dummy_size = 1024 * 1024 * 5 # 5 MB
    if not os.path.exists(dummy_file_path):
        with open(dummy_file_path, "wb") as f:
            f.write(os.urandom(dummy_size))
        print(f"Created dummy file: {dummy_file_path} ({dummy_size / (1024*1024):.1f} MB)")
        
    run_upload(dummy_file_path)

if __name__ == '__main__':
    main()