# Dockerized Mockoon

A Go library for testing HTTP APIs using [Mockoon](https://mockoon.com/) in Docker containers. This project provides an easy way to spin up mock services for testing your Go applications without depending on external APIs.

## Overview

This library wraps Mockoon's CLI in Docker containers, allowing you to:

* Start and stop mock HTTP services programmatically
* Use pre-configured mock data for testing
* Generate API clients from OpenAPI specifications
* Write integration tests with real HTTP calls to mock services

## Features

* 🐳 **Docker-based**: Isolated mock services in containers
* 🚀 **Easy setup**: Simple Go interface for starting/stopping mocks
* 🧪 **Test-friendly**: Designed for integration testing
* 🔧 **Configurable**: Support for multiple services and ports

## Prerequisites


1. **Docker** (or equivalent container runtime)
2. **Go** 1.23.0 or later
3. **Make** (optional, for using the Makefile)

## Installation

```bash
go get github.com/diegoavanzini/dockerized_mockoon
```

## Quick Start

### 1. Basic Usage

```go
package main

import (
    "context"
    "testing"
    "github.com/antihax/optional"
    "github.com/diegoavanzini/dockerized_mockoon"
    "github.com/diegoavanzini/dockerized_mockoon/clients/petstore"
)

func TestPetStoreAPI(t *testing.T) {
    // Create a new DockerMockoon instance
    serviceMockFactory, err := dockerized_mockoon.NewDockerMockoon()
    if err != nil {
        t.Fatalf("Failed to create DockerMockoon: %v", err)
    }
    
    // Add the petstore service
    serviceMockFactory.AddService(dockerized_mockoon.Service{
        Name:        "petstore",
        ExposedPort: 3005,
    })
    
    // Start the mock service
    mockedPetstore, addr, err := serviceMockFactory.Start("petstore")
    if err != nil {
        t.Fatalf("Failed to start mock service: %v", err)
    }
    defer mockedPetstore.TearDown()
    
    // Use the generated client to make API calls
    petstoreConfig := petstore.NewConfiguration()
    petstoreConfig.BasePath = *addr
    client := petstore.NewAPIClient(petstoreConfig)
    
    // Make API calls to the mock service
    pets, _, err := client.PetApi.FindPetsByStatus(context.Background(), &petstore.FindPetsByStatusOpts{
        Status: optional.NewString("available"),
    })
    if err != nil {
        t.Fatalf("Failed to call API: %v", err)
    }
    
    // Assert the response
    if len(pets) == 0 {
        t.Fatal("Expected at least one pet in response")
    }
}
```

### 2. Generate API Clients

Generate Go clients from OpenAPI specifications:

```bash
make client
```

This will:

* Download the OpenAPI generator CLI
* Generate a Go client from the petstore OpenAPI spec
* Place the generated client in `./clients/petstore/`

## Project Structure

```
├── docker_mockoon.go          # Main library implementation
├── docker_mockoon_test.go     # Tests demonstrating usage
├── go.mod                     # Go module definition
├── Makefile                   # Build automation
├── clients/                   # Generated API clients
│   └── petstore/             # Petstore API client
├── data/                     # Mock data configurations
│   └── petstore/
│       ├── petstore.json     # Mockoon service configuration
│       └── petstore.yaml     # OpenAPI specification
└── README.md                 # This file
```

## API Reference

### DockerMockoon Interface

```go
type IDockerMockoon interface {
    Start(serviceToMock string) (IDockerMockoon, *string, error)
    TearDown() error
    AddService(service Service) error
}
```

### Service Configuration

```go
type Service struct {
    Name        string  // Service name (must match directory in data/)
    ExposedPort int     // Port to expose the service on
}
```

### Main Functions

* `NewDockerMockoon()`: Creates a new DockerMockoon instance
* `AddService(service Service)`: Adds a service configuration
* `Start(serviceToMock string)`: Starts a specific mock service
* `TearDown()`: Stops and cleans up the Docker container

## Available Make Commands

```bash
make client         # Generate API clients from OpenAPI specs
make delete_client  # Remove generated client code
make test          # Run all tests
make build         # Build the Go module
make clean         # Clean up generated files
```

## Adding New Mock Services


1. **Create service data**: Add a new directory under `data/` with your service name
2. **Add Mockoon config**: Create a `.json` file with your Mockoon service configuration
3. **Add OpenAPI spec**: Include a `.yaml` file with your OpenAPI specification
4. **Generate client**: Run `make client` or manually generate using the OpenAPI generator
5. **Use in tests**: Add the service using `AddService()` and start it with `Start()`

## Example: Petstore Service

The project includes a complete example using the Swagger Petstore API:

* **Mock data**: `data/petstore/petstore.json` - Mockoon configuration with endpoints and responses
* **OpenAPI spec**: `data/petstore/petstore.yaml` - Complete API specification
* **Generated client**: `clients/petstore/` - Auto-generated Go client
* **Tests**: `docker_mockoon_test.go` - Integration tests demonstrating usage

## Dependencies

* **dockertest/v3**: For Docker container management
* **Mockoon CLI**: For running the mock services
* **OpenAPI Generator**: For generating API clients

## How It Works


1. **Service Configuration**: Define mock services with names and ports
2. **Docker Container**: Spins up a Mockoon container with your service configuration
3. **Mock Data**: Uses pre-configured JSON files with mock responses
4. **API Client**: Generated Go clients make real HTTP calls to the mock service
5. **Testing**: Write integration tests that call the mock APIs

The library handles:

* Container lifecycle (start/stop/cleanup)
* Port mapping and networking
* Service discovery and addressing
* Error handling and timeouts

## Contributing


1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass with `make test`
5. Submit a pull request

## License

This project is open source. Please check the license file for details.

## Troubleshooting

### Common Issues


1. **Docker not running**: Ensure Docker is installed and running
2. **Port conflicts**: Check that the specified ports are available
3. **Generated client errors**: Re-run `make client` to regenerate clients
4. **Test failures**: Ensure all dependencies are installed and Docker is accessible
5. **Container startup issues**: Check Docker logs for Mockoon service errors

### Debug Tips

* Use `docker ps` to see running containers
* Check container logs with `docker logs <container-id>`
* Verify mock data files are properly formatted JSON
* Ensure OpenAPI specs are valid YAML

## Examples in Action

The project demonstrates a complete testing workflow:


1. **Mock Service Setup**: Configure a petstore service with realistic endpoints
2. **Client Generation**: Auto-generate Go clients from OpenAPI specifications
3. **Integration Testing**: Write tests that make real HTTP calls to mock services
4. **Service Lifecycle**: Proper container management with cleanup

This approach enables:

* **Reliable Testing**: Consistent mock responses for predictable tests
* **Development Speed**: No need to wait for external services
* **Offline Development**: Work without internet connectivity
* **Realistic Testing**: Use actual HTTP calls rather than mocked functions


