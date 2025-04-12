# Email Breach Checker

A scalable web application that allows users to check if their email addresses have been compromised in data breaches.

## Features

- Web interface for users to check their email addresses
- RESTful API for email validation (supports both GET and POST methods)
- PostgreSQL database for storing compromised emails
- Redis cache for improved performance
- Nginx for reverse proxy and load balancing
- Containerized deployment with Docker Compose
- Horizontal scaling support for web and API services
- Security features including rate limiting, CORS, and secure headers

## Architecture

This application follows a microservices architecture with the following components:

### Components

1. **Web Frontend**
   - Go-based web server rendering HTML templates
   - Provides user interface for email breach checking
   - Displays container information for load balancing visualization
   - Scales horizontally to handle increased user traffic

2. **API Service**
   - RESTful API for checking email breach status
   - Connects to PostgreSQL database for breach data
   - Uses Redis for caching frequent queries
   - Scales horizontally to handle increased request load

3. **Nginx**
   - Acts as a reverse proxy and load balancer
   - Distributes traffic across multiple web and API instances
   - Provides a unified entry point for the application

4. **PostgreSQL Database**
   - Stores the database of compromised emails
   - Persistent storage with volume mounting

5. **Redis Cache**
   - Caches API responses to reduce database load
   - Improves response times for frequently checked emails

### Architecture Diagram

```
                   ┌─────────────┐
                   │    Users    │
                   └──────┬──────┘
                          │
                          ▼
                   ┌─────────────┐
                   │    Nginx    │
                   │Load Balancer│
                   └──────┬──────┘
                          │
              ┌───────────┴───────────┐
              │                       │
              ▼                       ▼
┌─────────────────────────┐   ┌─────────────────────┐
│     Web Frontend        │   │     API Service     │
│  (Multiple Instances)   │──►|(Multiple Instances) │
└─────────────────────────┘   └─────────┬───────────┘
                                        │
                             ┌──────────┴──────────┐
                             │                     │
                             ▼                     ▼
                      ┌─────────────┐      ┌─────────────┐
                      │  PostgreSQL │      │    Redis    │
                      │  Database   │      │    Cache    │
                      └─────────────┘      └─────────────┘
```

### Data Flow

1. User requests arrive at the Nginx reverse proxy
2. Web interface requests are forwarded to the Web Frontend service
3. Email check requests are sent from the Web Frontend to the API Service
4. The API Service first checks the Redis Cache for the email
5. If not found in cache, the API Service queries the PostgreSQL Database
6. Results are cached in Redis for future requests
7. The response is returned to the user through the Web Frontend

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Running the Application

1. Clone the repository:
   ```
   git clone https://github.com/zawmyolatt/breach-checker.git
   cd breach-checker
   ```

2. Start the services:
   ```
   docker-compose up -d
   ```

3. Access the application:
   - Web interface: http://localhost
   - API endpoint: http://localhost/api/check

### API Usage

The API supports both GET and POST methods:

#### GET Request
```
GET /api/check?email=test@example.com
```

#### POST Request
```
POST /api/check
Content-Type: application/json

{
  "email": "test@example.com"
}
```

### Testing

To test if an email is compromised, use one of the following sample emails:
- test@example.com
- compromised@example.com
- breach@example.com

## Security Features

- Rate limiting to prevent brute force attacks
- CORS configuration for API security
- HTTP security headers
- Input validation to prevent injection attacks
- Containerization for isolation
- Nginx as a reverse proxy to hide backend services

## Auto-Scaling Configuration

The application is designed to scale horizontally to handle increased load. Here's how the auto-scaling is configured:

### Container Information Display

The web interface displays which container is serving your request. This helps visualize the load balancing in action:

- Each web container shows its container ID and name at the bottom of the page
- This allows you to see which container is handling your request as you refresh the page

### Scaling Services

You can scale the services using Docker Compose:

```bash
# Scale the API service to 3 instances
docker-compose up -d --scale api=3

# Scale the web service to 3 instances
docker-compose up -d --scale web=3
```

When you scale a service, Nginx will automatically distribute traffic across all available instances.

### Load Balancing Strategy

- **Web Service**: Nginx uses round-robin load balancing to distribute user requests across multiple web containers
- **API Service**: Nginx uses round-robin load balancing to distribute API requests across multiple API containers

### Resource Limits

Each service has configured resource limits to ensure optimal performance:

- **Web**: 0.5 CPU cores, 512MB memory per instance
- **API**: 0.5 CPU cores, 512MB memory per instance
- **Database**: 1.0 CPU cores, 1GB memory
- **Redis**: 0.5 CPU cores, 512MB memory
- **Nginx**: 0.5 CPU cores, 256MB memory

### Health Checks

All services include health checks to ensure that only healthy instances receive traffic:

- **API**: Checks the `/health` endpoint every 10 seconds
- **Web**: Checks the `/health` endpoint every 10 seconds
- **Database**: Checks connection availability every 10 seconds
- **Redis**: Checks connection availability every 10 seconds

If a service fails its health check, it will be automatically restarted.

### Container Optimization

The Dockerfiles for both web and API services have been optimized for auto-scaling:

1. **Multi-stage builds** to create smaller, more efficient containers
2. **Non-root user execution** for improved security
3. **Built-in health checks** for better orchestration
4. **Optimized dependencies** to reduce container size and startup time

### Testing Load Balancing

To see the load balancing in action:

1. Scale up the web service:
   ```bash
   docker-compose up -d --scale web=3
   ```

2. Refresh the page multiple times
   - You should see different container IDs handling your requests
   - This demonstrates that Nginx is distributing traffic across multiple containers

3. Simulate high load:
   ```bash
   # Using Apache Bench to send 1000 requests with 10 concurrent connections
   ab -n 1000 -c 10 http://localhost/
   ```

4. Monitor container performance:
   ```bash
   docker stats
   ```

## Security Risks and Mitigation

### Risks

1. **Data Exposure**: The database contains sensitive information about compromised emails.
   - **Mitigation**: Implement proper access controls, encryption at rest, and network isolation.

2. **API Abuse**: The API could be abused for email harvesting or enumeration.
   - **Mitigation**: Implement rate limiting, require authentication for bulk queries, and log suspicious activity.

3. **DDoS Attacks**: The service could be targeted by distributed denial-of-service attacks.
   - **Mitigation**: Use rate limiting, CDN services, and consider cloud-based DDoS protection.

4. **SQL Injection**: Improperly sanitized inputs could lead to SQL injection.
   - **Mitigation**: Use parameterized queries (already implemented) and input validation.

5. **Insufficient Logging**: Lack of proper logging could hinder incident response.
   - **Mitigation**: Implement comprehensive logging and monitoring.

6. **Insecure Communication**: Data transmitted in plaintext could be intercepted.
   - **Mitigation**: Enable HTTPS by configuring SSL certificates in Nginx.

7. **Container Vulnerabilities**: Outdated container images could contain vulnerabilities.
   - **Mitigation**: Regularly update container images and scan for vulnerabilities.

## Future Improvements

1. Enable HTTPS with proper SSL certificates
2. Implement user authentication for advanced features
3. Add more comprehensive logging and monitoring
4. Implement a notification system for newly discovered breaches
5. Add support for checking passwords against breach databases (with proper hashing)
6. Implement a CI/CD pipeline for automated testing and deployment
7. Add metrics collection and dashboards for performance monitoring
8. Implement database replication for high availability

## CI/CD Pipeline

This project includes a GitHub Actions workflow for continuous integration:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │
│  Git Push   │────►│    Tests    │────►│   Docker    │
│             │     │             │     │   Build     │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Workflow Steps

1. **Testing**: Runs Go tests for both API and Web services
2. **Building**: Builds Docker images for all services to verify they can be built successfully

### Running the CI Pipeline Locally

You can simulate the CI pipeline locally with these commands:

```bash
# Run tests
cd api && go test -v ./... && cd ..
cd web && go test -v ./... && cd ..

# Build Docker images
docker build -t breach-checker-api:latest ./api
docker build -t breach-checker-web:latest ./web
docker build -t breach-checker-nginx:latest ./nginx

# Run the application with docker-compose
docker-compose up -d
```

This CI pipeline ensures that your code is tested and can be built successfully, providing confidence that the application will work as expected.
