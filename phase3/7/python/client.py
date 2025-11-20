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