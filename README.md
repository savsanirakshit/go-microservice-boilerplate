
**Project Overview**

This repository houses a microservice written in Go, following the Model-View-Controller (MVC) structure. The service handles API calls by routing them through designated routes, handlers, services, and repositories. Additionally, it includes configurations for logging, environment variables, database management, and migrations.

**Technologies Used**
- Language: Go
- Routing: Gorilla Mux
- ORM: pgmodel
- Database: PostgreSQL

**Structure**
The project follows the MVC pattern for better organization and scalability:

- Routes: Defines the endpoints and their corresponding handlers.
- Handlers: Processes incoming requests and invokes appropriate service methods.
- Services: Implements the business logic and interacts with repositories.
- Repositories: Handles data access and storage operations.

**Configuration**

- Logger: Configured in logger/servicelogger.go.
- Environment Variables: Managed in app.config.
- Database Configuration: Handled in db/config.go.
- Migrations: SQL migration files are available in db/migration/*.sql.

**Installation**

To set up the project locally, follow these steps:

  1.Clone this repository.
  
  2.Install dependencies using go mod tidy.
  
  3.Configure environment variables as specified in app.config.
  
  4.Configure database settings in db/config.go.

After installation, the microservice is ready to use. Start the server, and it will listen for incoming API requests based on the defined routes.

Note : import Go microservice.postman_collection.json in postman to test / debug route
