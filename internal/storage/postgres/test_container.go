package postgres

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var hostName = "localhost"

const defaultPostgresPort = "5432/tcp"

type Container struct {
	resource *dockertest.Resource
}

func NewContainer(config *Config, connectFn func() error) (*Container, error) {
	hostPort, err := getFreePort()
	if err != nil {
		return nil, fmt.Errorf("could not get free hostPort: %w", err)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15-alpine",
			Env: []string{
				"POSTGRES_USER=su",
				"POSTGRES_PASSWORD=su",
				"POSTGRES_DB=postgres",
				"listen_addresses = '*'",
			},
			PortBindings: map[docker.Port][]docker.PortBinding{
				defaultPostgresPort: {{
					HostIP:   hostName,
					HostPort: strconv.Itoa(hostPort),
				}},
			},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
	if err != nil {
		return nil, fmt.Errorf("could not create a container: %w", err)
	}

	container := &Container{
		resource: resource,
	}
	config.Port = resource.GetPort(defaultPostgresPort)
	config.Host = hostName
	if err := pool.Retry(func() error {
		ConnectionString(config)
		return connectFn()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return container, nil
}

func (c *Container) Purge() error {
	return c.resource.Close()
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
