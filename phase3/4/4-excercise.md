# Solution: Go Server and Python Client
## 1. File Structure
Ensure your project structure is set up correctly, assuming the code generation step (Lab 3.6) was successful:

```
├── proto/
│   └── calculator.proto
├── go/
│   ├── calculatorpb/     # Generated Go code
│   └── server/           # Go Server code
└── python/
    ├── calculator_pb2.py      # Generated Python message code
    ├── calculator_pb2_grpc.py # Generated Python stub code
    └── client.py              # Python Client code
```
## 2. Go Server Implementation
This code implements both the Add and Subtract RPC methods and starts the server on port 50051.

`go/server/main.go`

```go
package main

import (
	"context"
	"log"
	"net"

	// Import the generated protobuf code (adjust path if needed)
	pb "./calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server structure implements the generated CalculatorServiceServer interface.
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

// 4.2. Implementation: Add Method
func (*server) Add(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	// Example: Check context for potential timeout (crucial in distributed systems)
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Errorf(codes.DeadlineExceeded, "Request timed out")
	}

	num1 := req.GetNum1()
	num2 := req.GetNum2()
	result := num1 + num2

	log.Printf("Go Server: Received Add request: %d + %d = %d", num1, num2, result)

	// Create and return the response message
	res := &pb.AddResponse{
		Result: result,
	}
	return res, nil
}

// Implementation: Subtract Method
func (*server) Subtract(ctx context.Context, req *pb.SubtractRequest) (*pb.SubtractResponse, error) {
	num1 := req.GetNum1()
	num2 := req.GetNum2()
	result := num1 - num2

	log.Printf("Go Server: Received Subtract request: %d - %d = %d", num1, num2, result)

	res := &pb.SubtractResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	// Server listens on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fatal: failed to listen: %v", err)
	}

	// Create a new gRPC server instance (insecure for development)
	s := grpc.NewServer()

	// Register the implemented service
	pb.RegisterCalculatorServiceServer(s, &server{})

	log.Println("✅ Go Server started successfully on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fatal: failed to serve: %v", err)
	}
}
```
### 3. Python Client Implementation
This client connects to the Go server on port 50051 and performs the two required RPC calls.

`python/client.py`


```python
import grpc

# Import generated modules
from calculator import calculator_pb2 as pb
from calculator import calculator_pb2_grpc as rpc

def run():
    # 4.2. Connection: Create an insecure channel to the Go server
    with grpc.insecure_channel('localhost:50051') as channel:
        
        # Create the client stub for the CalculatorService
        stub = rpc.CalculatorServiceClient(channel)
        
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

if __name__ == '__main '__main__':
    run()
```

## 4. Execution Steps

1. Generate Code:
```bash
protoc --proto_path=proto --go_out=go --go_opt=paths=source_relative proto/calculator.proto
python -m grpc_tools.protoc --proto_path=proto --python_out=python --pyi_out=python --grpc_out=python --plugin=protoc-gen-grpc=`which grpc_python_plugin` proto/calculator.proto
```
Note: We use the full cmd in python to generate the gRPC stubs.

2. Start the Go Server: In one terminal, navigate to go/server and run:

```bash
go run main.go
# Output should be: ✅ Go Server started successfully on port 50051
```
3. Run the Python Client: In a second terminal, navigate to python and run:

```bash
python client.py
```
### Expected Output:
Go Server Terminal:
```bash
Go Server: Received Add request: 10 + 5 = 15
Go Server: Received Subtract request: 20 - 7 = 13
```
Python Client Terminal:
```bash
➕ Add Call Success:
   Request: 10 + 5
   Result: 15
--------------------
➖ Subtract Call Success:
   Request: 20 - 7
   Result: 13
```