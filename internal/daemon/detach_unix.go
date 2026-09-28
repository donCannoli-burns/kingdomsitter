//go:build !windows

package daemon

import (
	"os"
	"syscall"
)

func detachedSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}

func stopProcess(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
