## Folder Structure

The following is the folder structure of the project:

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

### Description of Folders

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

## Entity Framework Code Generation

- Ent Code & Gql Schema

```
go run ./cgfx/ent/entc.go

```

- GQL Resolvers

```
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen
```
- GRPC Code Gen

```
protoc --go_out=. --go-grpc_out=. ./app/grpc/server/proto/product.proto
```

## To Do

- [ ] Setup authentication or EF GraphQL.
- [ ] Setup Custom GraphQL Server & Resolvers.
- [ ] Setup Kafka Service for Communication between services.
- [ ] Setup casbin for authorization.
- [ ] Setup websockets & Grpc.
