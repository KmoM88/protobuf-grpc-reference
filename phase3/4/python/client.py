import grpc

# Import generated modules
import calculator_pb2 as pb
import calculator_pb2_grpc as rpc

def run():
    # 4.2. Connection: Create an insecure channel to the Go server
    with grpc.insecure_channel('localhost:50051') as channel:
        
        # Create the client stub for the CalculatorService
        stub = rpc.CalculatorServiceStub(channel)
        
        # --- RPC Call: Add (10 + 5) ---
        
        # 4.2. Create the request message
        add_request = pb.AddRequest(num1=10, num2=5)
        
        try:
            # Make the Unary RPC call
            add_response = stub.Add(add_request)
            print(f"➕ Add Call Success:")
            print(f"   Request: {add_request.num1} + {add_request.num2}")
            print(f"   Result: {add_response.result}")
            
        except grpc.RpcError as e:
            # Handle gRPC errors (e.g., DeadlineExceeded, UNAVAILABLE)
            print(f"Error in Add call: {e.details()}")

        print("-" * 20)
        
        # --- RPC Call: Subtract (20 - 7) ---
        
        subtract_request = pb.SubtractRequest(num1=20, num2=7)
        
        try:
            subtract_response = stub.Subtract(subtract_request)
            print(f"➖ Subtract Call Success:")
            print(f"   Request: {subtract_request.num1} - {subtract_request.num2}")
            print(f"   Result: {subtract_response.result}")
            
        except grpc.RpcError as e:
            print(f"Error in Subtract call: {e.details()}")

if __name__ == '__main__':
    run()