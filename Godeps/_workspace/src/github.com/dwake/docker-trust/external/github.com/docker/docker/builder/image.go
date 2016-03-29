package builder

import "github.com/dwake/docker-trust/external/github.com/docker/engine-api/types/container"

// Image represents a Docker image used by the builder.
type Image interface {
	ID() string
	Config() *container.Config
}
