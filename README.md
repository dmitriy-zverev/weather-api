# Weather API

A high-performance weather API service built with Go that provides current weather, daily forecasts, and hourly weather data for any city worldwide. The service uses Redis for caching to ensure fast response times and reduce external API calls.

## Features

- **Current Weather**: Get real-time weather data for any city
- **Daily Forecasts**: Retrieve weather forecasts for multiple days
- **Hourly Weather**: Access detailed hourly weather data for today
- **Redis Caching**: Fast response times with intelligent caching
- **Rate Limiting**: Built-in request rate limiting (7 requests per minute)
- **Docker Support**: Complete containerization with Docker and Docker Compose
- **Environment Configuration**: Flexible configuration for development and production
- **Connection Pooling**: Efficient Redis connection management
- **RESTful API**: Clean and intuitive API endpoints
- **Error Handling**: Comprehensive error responses and logging

## Tech Stack

- **Language**: Go 1.24.5
- **Cache**: Redis 7
- **HTTP Router**: Go standard library
- **External APIs**: Open-Meteo API for weather data
- **Configuration**: Environment variables with .env support

## Prerequisites

- Go 1.24.5 or higher
- Docker (for Redis)
- Git

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/dmitriy-zverev/weather-api.git
cd weather-api
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory:

```env
REDIS_URL=localhost
REDIS_PORT=6379
PORT=8080
PLATFORM=dev
```

### 3. Start Redis

```bash
chmod +x setup_redis.sh
./setup_redis.sh
```

This will start a Redis container on port 6379.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Start the Server

```bash
chmod +x start_server.sh
./start_server.sh
```

The server will start on `http://localhost:8080`.

## Docker Deployment

For production deployment using Docker:

### 1. Create Docker Network

```bash
chmod +x scripts/create_network.sh
./scripts/create_network.sh
```

### 2. Build Docker Image

```bash
chmod +x scripts/build_docker.sh
./scripts/build_docker.sh
```

### 3. Run with Docker

```bash
chmod +x scripts/run_docker.sh
./scripts/run_docker.sh
```

This will start both Redis and the Weather API in Docker containers with proper networking.

### 4. Restart Services

```bash
chmod +x scripts/restart_docker.sh
./scripts/restart_docker.sh
```

## API Documentation

### Base URL

```
http://localhost:8080/v1
```

### Endpoints

#### Health Check

```http
GET /v1/
```

**Response:**
```
OK
```

#### Weather Data

```http
GET /v1/weather
```

**Request Body:**
```json
{
  "city": "London",
  "forecast_type": "current"
}
```

**Forecast Types:**

1. **Current Weather** (`current`)
   ```json
   {
     "city": "London",
     "forecast_type": "current"
   }
   ```

2. **N-Day Forecast** (`n_days`)
   ```json
   {
     "city": "London",
     "forecast_type": "n_days",
     "days": 5
   }
   ```

3. **Today's Hourly Forecast** (`today_hourly`)
   ```json
   {
     "city": "London",
     "forecast_type": "today_hourly"
   }
   ```

### Response Examples

#### Current Weather Response
```json
{
  "current": {
    "temperature_2m": 22.5
  }
}
```

#### Daily Weather Response
```json
{
  "daily": {
    "time": ["2025-09-25", "2025-09-26", "2025-09-27"],
    "temperature_2m_max": [25.2, 23.8, 21.5],
    "temperature_2m_min": [15.1, 14.2, 12.8]
  }
}
```

#### Hourly Weather Response
```json
{
  "hourly": {
    "time": ["2025-09-25T00:00", "2025-09-25T01:00", "2025-09-25T02:00"],
    "temperature_2m": [18.5, 17.8, 17.2]
  }
}
```

## Project Structure

```
weather-api/
├── main.go                     # Application entry point
├── routes.go                   # Route constants and configuration
├── Dockerfile                  # Docker container configuration
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── .env                        # Environment variables (create this)
├── setup_redis.sh              # Redis setup script
├── start_server.sh             # Server startup script
├── scripts/                    # Docker deployment scripts
│   ├── build_docker.sh         # Build Docker image
│   ├── create_network.sh       # Create Docker network
│   ├── restart_docker.sh       # Restart Docker services
│   └── run_docker.sh           # Run Docker containers
└── internal/
    ├── handlers/               # HTTP request handlers
    │   ├── handlers.go         # Main handler implementations
    │   ├── models.go           # Request/response models
    │   └── errors.go           # Error handling utilities
    ├── api_handler/            # External API integration
    │   ├── fetch.go            # Weather data fetching logic
    │   ├── models.go           # API response models
    │   └── consts.go           # API constants
    └── cache/                  # Redis caching layer
        ├── cache.go            # Cache operations
        └── keys.go             # Cache key management
```

## Configuration

The application uses environment variables for configuration:

| Variable     | Description                     | Default   | Required |
|--------------|---------------------------------|-----------|----------|
| `REDIS_URL`  | Redis server address            | localhost | Yes      |
| `REDIS_PORT` | Redis server port               | 6379      | Yes      |
| `PORT`       | Application server port         | 8080      | Yes      |
| `PLATFORM`   | Environment platform (dev/prod) | prod      | Yes      |

## Development

### Running Tests

```bash
go test ./...
```

### Building the Application

```bash
go build -o weather-api
```

### Running with Custom Port

The application runs on port 8080 by default. To change this, modify the `PORT` constant in `routes.go`.

## Caching Strategy

The application implements intelligent caching using Redis:

- Weather data is cached to reduce external API calls
- Cache keys are structured for efficient retrieval
- Automatic cache invalidation ensures data freshness
- Configurable cache timeouts for different data types

## Error Handling

The API provides comprehensive error responses:

- **400 Bad Request**: Invalid request parameters
- **500 Internal Server Error**: Server-side errors
- **Detailed Error Messages**: Clear descriptions of what went wrong

## Rate Limiting

The API implements built-in rate limiting to ensure fair usage and prevent abuse:

- **Request Limit**: 7 requests per minute per client
- **Burst Capacity**: Up to 10 requests can be made in quick succession
- **Implementation**: Uses Go's `golang.org/x/time/rate` package
- **Scope**: Rate limiting is applied globally across all endpoints

When rate limits are exceeded, the API will return appropriate HTTP status codes and error messages.

## Performance

- **Redis Caching**: Reduces response times and external API calls
- **Connection Pooling**: Efficient Redis connection management with pool size of 10
- **Timeouts**: Configured read/write timeouts (1 second each) prevent hanging requests
- **Concurrent Handling**: Go's goroutines handle multiple requests efficiently
- **Rate Limiting**: Prevents API abuse while maintaining performance

## Roadmap

### Completed Features

- [x] **Docker Support**: Complete containerization with Docker networking
- [x] **Rate Limiting**: Request rate limiting (7 requests per minute)
- [x] **Environment Configuration**: Flexible dev/prod configuration
- [x] **Connection Pooling**: Redis connection pooling with timeouts

### Space for Improvement

- [ ] **Docker Compose**: Multi-service orchestration with docker-compose.yml
- [ ] **Metrics**: Prometheus metrics and monitoring
- [ ] **Logging**: Structured logging with different levels
- [ ] **Health Checks**: Advanced health check endpoints
- [ ] **Historical Data**: Access to historical weather data

### Infrastructure Improvements

- [ ] **Kubernetes Deployment**: K8s manifests and Helm charts
- [ ] **CI/CD Pipeline**: Automated testing and deployment
- [ ] **Load Balancing**: Multi-instance deployment support
- [ ] **Backup Strategy**: Redis data backup and recovery

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is free to use.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/dmitriy-zverev/weather-api/issues) page
2. Create a new issue with detailed information
3. Contact the maintainers

## Acknowledgments

- [Open-Meteo](https://open-meteo.com/) for providing free weather data
- [Redis](https://redis.io/) for the excellent caching solution
- [Go](https://golang.org/) for the robust programming language
- [Roadmap](https://roadmap.sh/projects/weather-api-wrapper-service) for inspiring
