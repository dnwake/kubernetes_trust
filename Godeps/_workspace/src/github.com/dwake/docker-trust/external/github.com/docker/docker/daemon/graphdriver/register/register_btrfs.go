// +build !exclude_graphdriver_btrfs,linux

package register

import (
	// register the btrfs graphdriver
	_ "github.com/dwake/docker-trust/external/github.com/docker/docker/daemon/graphdriver/btrfs"
)
