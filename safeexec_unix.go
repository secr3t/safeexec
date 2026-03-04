//go:build linux || darwin

package safeexec

import (
	"os/exec"
	"syscall"
)

func setupCmd(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Create a new process group for the child process.
	// This allows us to kill the process and all its descendants.
	cmd.SysProcAttr.Setpgid = true

	// On Linux, we can use Pdeathsig to ensure the child dies when the parent dies.
	// This is handled by the kernel.
	setPlatformSpecificSysProcAttr(cmd.SysProcAttr)
}
