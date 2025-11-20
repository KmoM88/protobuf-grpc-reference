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