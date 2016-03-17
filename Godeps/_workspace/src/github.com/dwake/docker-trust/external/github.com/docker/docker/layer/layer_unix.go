// +build linux freebsd darwin

package layer

import "github.com/dwake/docker-trust/external/github.com/docker/docker/pkg/stringid"

func (ls *layerStore) mountID(name string) string {
	return stringid.GenerateRandomID()
}
