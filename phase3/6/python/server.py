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