//go:build windows

package daemon

import (
	"os"
	"syscall"
)

func detachedSysProcAttr() *syscall.SysProcAttr {
	const (
		createNewProcessGroup = 0x00000200
		detachedProcess       = 0x00000008
	)
	return &syscall.SysProcAttr{
		CreationFlags: createNewProcessGroup | detachedProcess,
		HideWindow:    true,
	}
}

func stopProcess(process *os.Process) error {
	return process.Kill()
}
