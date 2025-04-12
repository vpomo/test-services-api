This is a test assignment.

Goal:
To develop a simple but architecturally competent micro-service API on Go, including registration, user authorization, as well as work with posts and comments. The interaction between microservices should take place through gRPC, and the project itself should be fully ready for deployment in Docker.

### Architecture requirements
The general structure of microservices:
#### API Gateway
— HTTP service that accepts requests from the client and proxies calls to the corresponding gRPC services.

#### User Service
— gRPC-the service responsible for:

    - User registration
    - Authentication (with the issuance of a JWT token)
    - Token validation for subsequent requests

#### Article Service
— gRPC-the service responsible for:

    - Creating, editing, and receiving posts (articles)
    - Adding comments to posts (both your own and others')

### Documentation
All HTTP API endpoints must be documented using Swagger (swagger.yaml).
Documentation should be available on the path /swagger

### Docker
All microservices and gateway should be assembled into Docker containers.
Configure the launch using Docker Compose