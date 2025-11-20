# gRPC & Protocol Buffers
Target Languages: Go, C++, Python

## Phase 1: Conceptual Foundations & The "Why"
Theory first. Understand the problem before writing the solution.

### 1. [Introduction to Modern RPC](phase1/1/1.md)
- [x] 1.1. The Evolution of APIs: SOAP vs. REST vs. GraphQL vs. gRPC.
- [x] 1.2. HTTP/2 Deep Dive: Multiplexing, Header Compression (HPACK), and binary framing.
- [x] 1.3. Why gRPC? (Strict contracts, Code generation, Low latency).
- [x] 1.4. Industry Use Cases: Microservices (internal traffic), Mobile, Browser (gRPC-Web).

### 2. [Protocol Buffers (The Data Layer)](phase1/2/2.md)
- [x] 2.1. Serialization concepts: Text (JSON/XML) vs. Binary (Protobuf).
- [x] 2.2. Anatomy of a .proto file (syntax = "proto3";, package).
- [x] 2.3. Scalar Value Types (int32, float, bool, string, bytes).
- [x] 2.4. Complex Types: message, enum, repeated (arrays), map.
- [x] 2.5. Handling Nullability and Optional fields in proto3.
- [x] 2.6. oneof: Handling union types.
- [x] 2.7. Any type: Embedding arbitrary messages.

## Phase 2: Environment & Tooling
Setting up a robust development workflow.

### 3. [Setup and Code Generation](phase2/3/3.md)
- [x] 3.1. Installing the Protocol Buffers Compiler (protoc).
- [x] 3.2. Go Setup: Installing protoc-gen-go and protoc-gen-go-grpc.
- [x] 3.3. Python Setup: Installing grpcio and grpcio-tools.
- [x] 3.4. C++ Setup: Configuring CMake/Bazel for gRPC.
- [x] 3.5. Modern Tooling (Recommended): Installing and configuring Buf (replaces complex protoc commands).
- [x] 3.6. Lab: Generate code for a simple User message in all three languages and inspect the output files.

## Phase 3: Core Implementation (Polyglot)
Mastering the 4 communication patterns. Focus on concurrency differences.

### 4. [Unary RPC (Simple Request-Response)](phase3/4/4.md)
- [x] 4.1. Defining a Unary service in .proto.
- [x] 4.2. Go Implementation:
	- [x] Implementing the Server interface.
	- [x] Using context.Context.
	- [x] Creating a Client connection (grpc.Dial).
- [x] 4.3. C++ Implementation:
	- [x] Synchronous Service implementation.
	- [x] Understanding grpc::Status and grpc::ServerContext.
	- [x] Managing string/bytes memory.
- [x] 4.4. Python Implementation:
	- [x] Implementing the Servicer class.
	- [x] Starting the grpc.server with a thread pool.
- [x] 4.5. Exercise: Build a "Calculator Service" (Add, Subtract) with a Go Server and Python Client.

### 5. [Server Streaming RPC (One Request, Many Responses)](phase3/5/5.md)
- [x] 5.1. Use cases: Real-time feeds, large dataset downloads.
- [x] 5.2. Defining stream in the return type.
- [x] 5.3. Go: Sending messages into the stream object.
- [x] 5.4. C++: Using ServerWriter loop.
- [x] 5.5. Python: Using generators (yield) for responses.
- [x] 5.6. Exercise: Build a "Stock Ticker" that streams random prices every second.

### 6. [Client Streaming RPC (Many Requests, One Response)](phase3/6/6.md)
- [x] 6.1. Use cases: File uploads, IoT sensor ingestion.
- [x] 6.2. Defining stream in the argument type.
- [x] 6.3. Go: Using Recv() inside a loop until EOF.
- [x] 6.4. C++: Using ServerReader to aggregate data.
- [x] 6.5. Python: Iterating over the request iterator.
- [x] 6.6. Exercise: Build an "Image Uploader" (upload chunks, return total size).

### 7. [Bidirectional Streaming RPC (Many Requests, Many Responses)](phase3/7/7.md)
- [x] 7.1. Use cases: Chat apps, Multiplayer games, Live synchronization.
- [x] 7.2. Go: Handling independent Send and Recv in goroutines.
- [x] 7.3. C++: Using ServerReaderWriter.
- [x] 7.4. Python: Consuming an iterator while yielding responses simultaneously.
- [x] 7.5. Exercise: Build a "Chat Room" (Python & Go Clients).

## Phase 4: Production Engineering & Best Practices
Moving from "it works" to "it scales".

### 8. [Schema Design & Evolution](phase4/8/8.md)
- [x] 8.1. Google's Style Guide for .proto.
- [x] 8.2. Field Numbering Rules: Why you never reuse numbers.
- [x] 8.3. Backward & Forward Compatibility strategies.
- [x] 8.4. Using reserved fields.
- [x] 8.5. Well-Known Types (Timestamp, Duration, Struct).

### 9. [Reliability & Error Handling](phase4/9/9.md)
- [x] 9.1. The gRPC Status Codes (OK, CANCELLED, DEADLINE_EXCEEDED, UNIMPLEMENTED, etc.).
- [x] 9.2. Go: Using the status and codes packages.
- [x] 9.3. C++: Catching exceptions vs checking Status.
- [x] 9.4. Deadlines & Timeouts:
	- [x] Setting deadlines in Go Contexts.
	- [x] Setting deadlines in C++ ClientContext.
	- [x] Why default timeouts are necessary.
- [x] 9.5. Retries: Configuring exponential backoff policies.

### 10. [Metadata & Security](phase4/10/10.md)
- [x] 10.1. What is Metadata? (Headers/Trailers).
- [x] 10.2. Sending Auth Tokens (JWT) via metadata.
- [x] 10.3. TLS/SSL: Setting up secure credentials.
- [x] 10.4. mTLS (Mutual TLS): Configuring certificates for Zero Trust.

### 11. Interceptors (Middleware)
- [ ] 11.1. Go: UnaryServerInterceptor and StreamServerInterceptor.
- [ ] 11.2. Python: Client and Server interceptors.
- [ ] 11.3. C++: AuthMetadataProcessor and generic interceptors.
- [ ] 11.4. Use cases: Logging, Tracing (OpenTelemetry), Authentication validation.

## Phase 5: Advanced Performance & Internals
Deep dive for high-performance requirements.

### 12. Advanced C++ & Go Tuning
- [ ] 12.1. C++ Arenas: Optimizing memory allocation for Protobuf.
- [ ] 12.2. C++ Move Semantics: Avoiding copies in message passing.
- [ ] 12.3. Go: Buffer reuse and Goroutine pooling.
- [ ] 12.4. Concurrency Models: Thread-per-request (C++) vs Async vs Goroutines.

### 13. Ecosystem & Integration
- [ ] 13.1. gRPC-Gateway: Exposing gRPC as REST/JSON automatically.
- [ ] 13.2. gRPC-Web: Connecting frontend apps (JS/TS).
- [ ] 13.3. Load Balancing strategies (L7 vs L4).
- [ ] 13.4. CLI Tools: grpcurl and Evans.

### 14. Internals (Under the Hood)
- [ ] 14.1. Wire Format: Understanding Base 128 Varints.
- [ ] 14.2. ZigZag Encoding for signed integers.
- [ ] 14.3. How Protobuf handles field tags and wire types binary-level.

## Phase 6: Capstone Project
Final Exam: Polyglot Distributed System.

### 15. Project: Distributed File Storage System
- [ ] 15.1. Define the Schema: FileService, MetaService, Blob messages.
- [ ] 15.2. Component A (Master Node - Go): Manage file metadata, locations, and auth tokens.
- [ ] 15.3. Component B (Storage Nodes - C++): High-performance streaming of file chunks to/from disk.
- [ ] 15.4. Component C (CLI Client - Python): Script to upload/download files and check status.
- [ ] 15.5. Requirement: Implement mTLS between all nodes.
- [ ] 15.6. Requirement: Handle network interruptions (resume upload).