# Microservice

- A microservice template for building robust, scalable applications in Golang.

## Implementation Details

### Packages

- [SQLC](https://github.com/sqlc-dev/sqlc) – Generates type-safe Go code from SQL queries, improving database interactions.
- [Viper](https://github.com/spf13/viper) – A powerful configuration management tool for environment variables, config files, and flags.
- [Echo](https://github.com/labstack/echo) – A high-performance, extensible, and minimalist web framework for building APIs.
- [Ent](https://github.com/ent/ent) – A Go entity framework for struct-based ORM, simplifying database modeling.
- [GraphQL](https://github.com/99designs/gqlgen) – A Go GraphQL server implementation with a focus on code generation.
- [WebSocket](https://github.com/gobwas/ws) – A [highly optimized](https://github.com/akash-aman/microservice/issues/21) [WebSocket implementation](https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/) for real-time communication.
- [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-go) – Provides observability (tracing, metrics, logs) for distributed systems.
- [Zap](https://github.com/uber-go/zap) – A high-performance, structured logging library by Uber.
- [Consul](https://github.com/hashicorp/consul) – A service discovery and configuration tool for dynamic service registration and health checking.
- [Fx](https://github.com/uber-go/fx) – A dependency injection framework for simplifying application initialization and lifecycle management.
- [Ants](https://github.com/panjf2000/ants) – A high-performance goroutine pool in Go.

### Custom
- Workerpool : Added Implementation adaptive and fix workerpool.

### Folder Structure

#### Top Level View

```
├── pkg
└── services
    ├── inventory
    └── products
```

- **pkg**: Contains shared packages for different services.

#### Product Service
```
├── app
│   ├── apis
│   │   ├── coupons
│   │   └── products
│   │       ├── get_all
│   │       └── get_by_id
│   │           ├── v1
│   │           │   ├── dtos
│   │           │   ├── endpoints
│   │           │   ├── handler
│   │           │   └── model
│   │           └── v2
│   ├── core
│   │   ├── constants
│   │   ├── contracts
│   │   ├── errors
│   │   ├── helpers
│   │   ├── middlewares
│   │   └── models
│   ├── data
│   ├── infra
│   │   ├── bulk-factory
│   │   ├── scheduler
│   │   └── states
│   ├── inits
│   └── view
├── cmd
├── conf
├── docs
├── logs
├── cgfx
├── server
└── sqlc
```

- **app**: Contains the main application code.
  - **apis**: API related code.
    - **coupons**: API endpoints and handlers for coupons.
    - **products**: API endpoints and handlers for products.
      - **get_all**: Endpoint to get all products.
      - **get_by_id**: Endpoint to get a product by ID.
        - **v1**: Version 1 of the get_by_id endpoint.
          - **dtos**: Data Transfer Objects for v1.
          - **endpoints**: Endpoint definitions for v1.
          - **handler**: Handlers for v1.
          - **model**: Models for v1.
        - **v2**: Version 2 of the get_by_id endpoint.
  - **core**: Core functionalities and utilities.
    - **constants**: Constant values used across the application.
    - **contracts**: Interface definitions and contracts.
    - **errors**: Error handling and definitions.
    - **helpers**: Helper functions and utilities.
    - **middlewares**: Middleware functions.
    - **models**: Data models and schemas.
  - **data**: Data access layer.
  - **grpc**: Grpc server implementation, and grpc client client setup.
  - **infra**: Infrastructure related code.
    - **bulk-factory**: Bulk processing utilities.
    - **scheduler**: Task scheduling utilities.
    - **states**: State management utilities.
  - **inits**: Initialization scripts and configurations.
  - **view**: View layer code.
- **cmd**: Command line interface related code.
- **conf**: Configuration files.
- **docs**: Documentation files.
- **cgfx**: Code generation frameworks like entity framework, sqlc, gqlgen.
- **logs**: Error & access logs.
- **server**: Server related code.
- **sqlc**: SQL code and database migrations.
