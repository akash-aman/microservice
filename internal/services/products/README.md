
## Folder Structure

The following is the folder structure of the project:

```
├── app
│   ├── apis
│   ├── core
│   │   ├── constants
│   │   ├── contracts
│   │   ├── errors
│   │   ├── helpers
│   │   └── models
│   ├── data
│   ├── infra
│   │   ├── bulk-factory
│   │   ├── scheduler
│   │   └── states
│   └── inits
├── cmd
├── conf
├── docs
├── server
└── sqlc
```

### Description of Folders

- **app**: Contains the main application code.
    - **apis**: API related code.
    - **core**: Core functionalities and utilities.
        - **constants**: Constant values used across the application.
        - **contracts**: Interface definitions and contracts.
        - **errors**: Error handling and definitions.
        - **helpers**: Helper functions and utilities.
        - **models**: Data models and schemas.
    - **data**: Data access layer.
    - **infra**: Infrastructure related code.
        - **bulk-factory**: Bulk processing utilities.
        - **scheduler**: Task scheduling utilities.
        - **states**: State management utilities.
    - **inits**: Initialization scripts and configurations.
- **cmd**: Command line interface related code.
- **conf**: Configuration files.
- **docs**: Documentation files.
- **server**: Server related code.
- **sqlc**: SQL code and database migrations.