// +build windows

package windows

import (
	"fmt"

	"github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/execdriver"
)

// Update updates resource configs for a container.
func (d *Driver) Update(c *execdriver.Command) error {
	return fmt.Errorf("Windows: Update not implemented")
}
