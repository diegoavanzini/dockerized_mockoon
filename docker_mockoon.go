package dockerized_mockoon

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
)

type IDockerMockoon interface {
	Start(serviceToMock string) (IDockerMockoon, *string, error)
	TearDown() error
	AddService(service Service) error
}

func (dm *DockerMockoon) startDocker() (*dockertest.Pool, error) {
	var pool *dockertest.Pool
	var err error
	if dockerHost := os.Getenv("DOCKER_HOST"); dockerHost != "" {
		pool, err = dockertest.NewPool(dockerHost)
	} else {
		// Try npipe first (Windows Docker Desktop)
		pool, err = dockertest.NewPool("npipe:////./pipe/docker_engine")
		if err != nil {
			// Fallback to default
			pool, err = dockertest.NewPool("")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to start Dockertest: %+v", err)
	}

	// Try to connect to docker
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("Could not connect to docker: %s", err)
	}

	return pool, nil
}

func (dm *DockerMockoon) AddService(service Service) error {
	if _, exists := dm.services[service.Name]; exists {
		return fmt.Errorf("Service with name '%s' already exists", service.Name)
	}
	service.dataPath = fmt.Sprintf("%s/%s", dm.rootDataPath, service.Name)
	dm.services[service.Name] = service
	return nil
}

type DockerMockoon struct {
	rootDataPath string
	pool         *dockertest.Pool
	resource     *dockertest.Resource
	services     map[string]Service
}

type Service struct {
	Name        string
	ExposedPort int
	dataPath    string
}

func NewDockerMockoon() (IDockerMockoon, error) {
	absPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	absPath = strings.ReplaceAll(absPath, "\\", "/")
	// Ensure the path ends with /data
	absPath = fmt.Sprintf("%s/data", absPath)
	return &DockerMockoon{
		rootDataPath: absPath,
		services:     map[string]Service{},
	}, nil
}

func (dm *DockerMockoon) Start(serviceNameToMock string) (IDockerMockoon, *string, error) {
	serviceToMock := dm.services[serviceNameToMock]
	// Try to use Docker endpoint from environment or default to npipe for Windows
	pool, err := dm.startDocker()
	if err != nil {
		return nil, nil, err
	}
	opt := &dockertest.RunOptions{
		Repository: "mockoon/cli",
		Tag:        "latest",
		Cmd: []string{
			"-d", "/data/" + serviceNameToMock + ".yaml",
			"-p", fmt.Sprintf("%d", serviceToMock.ExposedPort),
		},
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", serviceToMock.ExposedPort)},
		Mounts: []string{
			fmt.Sprintf("%s:/data", serviceToMock.dataPath),
		},
	}

	resource, err := pool.RunWithOptions(opt)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to start mockoon/cli: %+v", err)
	}

	// determine the port the container is listening on
	addr := fmt.Sprintf("http://%s/api/v3", net.JoinHostPort("localhost", resource.GetPort(fmt.Sprintf("%d/tcp", serviceToMock.ExposedPort))))

	// wait for the container to be ready
	pool.MaxWait = time.Minute * 2
	err = pool.Retry(func() error {
		url := addr
		resp, e := http.Get(url)
		if e != nil {
			return e
		}
		defer resp.Body.Close()
		// Mockoon returns 404 for root path, which means it's working
		if resp.StatusCode >= 500 {
			return fmt.Errorf("service not ready, status: %d", resp.StatusCode)
		}
		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to ping Mockoon: %+v", err)
	}

	return &DockerMockoon{
		pool:     pool,
		resource: resource,
	}, &addr, nil
}

func (dm *DockerMockoon) TearDown() error {
	// You can't defer rtu because os.Exit doesn't care for defer
	if err := dm.pool.Purge(dm.resource); err != nil {
		err = errors.Wrap(err, "Could not purge docker resource")
		return err
	}
	return nil
}
