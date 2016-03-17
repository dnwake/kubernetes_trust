// +build !windows

package daemon

import (
	"os"
	"os/signal"
	"syscall"

	psignal "github.com/dwake/docker-trust/external/github.com/docker/docker/pkg/signal"
)

func setupDumpStackTrap() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	go func() {
		for range c {
			psignal.DumpStacks()
		}
	}()
}
