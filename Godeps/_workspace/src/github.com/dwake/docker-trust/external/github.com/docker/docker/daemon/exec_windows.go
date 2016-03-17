package daemon

import (
	"github.com/dwake/docker-trust/external/github.com/docker/docker/container"
	"github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/execdriver"
	"github.com/dwake/docker-trust/external/github.com/docker/engine-api/types"
)

// setPlatformSpecificExecProcessConfig sets platform-specific fields in the
// ProcessConfig structure. This is a no-op on Windows
func setPlatformSpecificExecProcessConfig(config *types.ExecConfig, container *container.Container, pc *execdriver.ProcessConfig) {
}
