//go:build windows

package safeexec

import (
	"os/exec"
	"syscall"
)

func setupCmd(cmd *exec.Cmd) {
	// On Windows, there is no Setpgid or Pdeathsig in SysProcAttr.
	// We rely on the watchdog process to clean up children.
}

func killProcess(p *Process) error {
	if p == nil || p.Process == nil {
		return nil
	}
	// On Windows, we just kill the process itself.
	// The watchdog or other mechanisms should handle child processes if necessary.
	return p.Process.Kill()
}

func setupWatchdogCmd(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// DETACHED_PROCESS (0x00000008) and CREATE_NEW_PROCESS_GROUP (0x00000200)
	// ensure the watchdog process doesn't get killed when the parent console/process is closed.
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008
	cmd.SysProcAttr.HideWindow = true
}
