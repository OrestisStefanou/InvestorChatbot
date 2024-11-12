# InvestorChatbot
A chatbot focused to help investors make data driven decisions 

## Code structure
investbot/
├── cmd/
│   └── investbot/            # Main application entry point
│       └── main.go             # Initializes the app and starts the HTTP server
├── pkg/                        # Library code that can be imported by other applications
│   ├── domain/                 # Data structures and types (business domain models)
│   │   └── user.go             # Example: User struct
│   ├── handlers/               # HTTP handler functions for each route
│   │   └── user_handler.go     # Example: User HTTP handlers
│   ├── services/               # Core business logic (service layer)
│   │   └── user_service.go     # Example: User service functions
│   ├── repositories/           # Database access layer (repository pattern)
│   │   └── user_repository.go  # Example: User DB interactions
│   ├── utils/                  # Utility functions/helpers
│   │   └── response.go         # Example: Response formatting functions
│   └── config/                 # Configuration loading (e.g., environment variables)
│       └── config.go           # Configuration management
├── migrations/                 # Database migration files
├── docs/                       # API documentation (Swagger, Postman, etc.)
└── go.mod                      # Go module file
