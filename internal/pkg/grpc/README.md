gRPC supports four types of communication patterns, each suited for different use cases. These communication models leverage Protocol Buffers (protobuf) for efficient serialization and HTTP/2 for fast, low-latency communication.


### ✅ 1. Unary RPC (Request-Response)
- Definition:
- The simplest form of gRPC where the client sends one request and receives one response.
- Use Case:
- Simple database lookups, user authentication, or fetching single records.

Example:

rpc GetProduct (ProductRequest) returns (ProductResponse);

- Flow:

```
Client  ---- Request ----> Server  
Client  <---- Response --- Server
```


### ✅ 2. Server Streaming RPC
- Definition:
- The client sends a single request, but the server sends multiple responses (a stream).
- Use Case:
- Streaming logs, real-time market data, or sending large datasets in chunks.

Example:

rpc ListProducts (ProductRequest) returns (stream ProductResponse);

- Flow:

```
Client  ---- Request ----> Server  
Client  <---- Stream of Responses --- Server
```



### ✅ 3. Client Streaming RPC
- Definition:
- The client sends multiple requests (a stream) to the server, and the server sends one final response.
- Use Case:
- File uploads, IoT sensor data, or collecting telemetry data.

Example:

rpc UploadProductImages (stream ImageRequest) returns (UploadResponse);

- Flow:

```
Client  ---- Stream of Requests ----> Server  
Client  <---- Single Response --------- Server
```


### ✅ 4. Bidirectional Streaming RPC
- Definition:
- Both the client and the server send streams of messages to each other simultaneously.
- Use Case:
- Chat applications, real-time collaboration tools, or multiplayer games.

Example:

rpc Chat (stream ChatMessage) returns (stream ChatMessage);

- Flow:

```
Client <----> Stream of Messages <----> Server
```



### 🔍 Summary of gRPC Communication Types:

| **Type**             | **Request**            | **Response**             | **Use Case**                            |
|----------------------|------------------------|--------------------------|-----------------------------------------|
| **Unary**            | Single                 | Single                   | Simple CRUD operations                  |
| **Server Streaming** | Single                 | Multiple (Stream)        | Logs, file downloads                    |
| **Client Streaming** | Multiple (Stream)      | Single                   | File uploads, telemetry data            |
| **Bidirectional**    | Multiple (Stream)      | Multiple (Stream)        | Chat apps, real-time collaboration      |