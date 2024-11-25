# Distributed Task Queue in Go

This project implements a **distributed task queue** in Go with Redis as the backend. It features task retries, metrics collection using Prometheus, a REST API for task submission, and graceful shutdown handling.

## Features

- **Task Queueing:** Adds tasks to a queue using Redis.
- **REST API:** Provides an endpoint for task submission.
- **Task Processing:** Processes tasks concurrently with retry logic for failed tasks.
- **Metrics Collection:** Tracks task processing success, failure, retries, and processing time using Prometheus.
- **Graceful Shutdown:** Handles termination signals and cleans up resources safely.

## Requirements

- **Go**: Version 1.20 or higher.
- **Redis**: Installed and running on `localhost:6379`.
- **Prometheus**: To scrape metrics exposed on `:2112/metrics`.

## Installation

1. Clone this repository:

   ```bash
   git clone <repo-url>
   cd distributed-task-queue
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

### 1. Running Redis

Make sure Redis is running locally on `localhost:6379`. You can start Redis using:

```bash
redis-server
```

### 2. Run the Worker

To start the worker:

```bash
make run-worker
```

### 3. Run the API Server

To start the API server:

```bash
make run-api
```

### 4. Submit Tasks via the API

You can submit tasks to the queue using a REST API endpoint. Here's an example:

#### Endpoint:

```
POST http://127.0.0.1:8080/tasks
```

#### Request Body:

```json
{
  "payload": "task_payload_here"
}
```

#### Example cURL Command:

```bash
curl -X POST http://127.0.0.1:8080/tasks \
     -H "Content-Type: application/json" \
     -d '{"payload": "my_first_task"}'
```

#### Response:

```json
{
  "task_id": "task1699895678901234567"
}
```

### 5. Metrics

Visit [http://127.0.0.1:2112/metrics](http://127.0.0.1:2112/metrics) to view metrics.

## Testing

To run unit tests:

```bash
make test
```

## Project Structure

```
distributed-task-queue/
│
├── internal/
│   ├── queue/           # Redis queue implementation
│   ├── monitoring/      # Metrics setup
│   └── task/            # Task model
│
├── cmd/
│   ├── worker/main.go   # Main worker application
│   ├── api/main.go      # API server for task submission
│
├── go.mod               # Dependencies
├── README.md            # Documentation
└── Makefile             # Automation tasks
```

## Makefile Commands

- `make run-worker`: Runs the worker.
- `make run-api`: Runs the API server.
- `make build-worker`: Builds the worker binary.
- `make build-api`: Builds the API server binary.
- `make build-all`: Builds both worker and API server binaries.
- `make test`: Executes unit tests.
- `make clean`: Cleans up build artifacts.
- `make fmt`: Formats the codebase.
- `make lint`: Runs lint checks on the code.

## Future Improvements

- Support for task prioritization.
- Integration with Kubernetes for scalability.
- Implement a web interface for task submission and monitoring.

## License

Distributed under the MIT License. See `LICENSE` for more information.
