//go:build windows

package safeexec

import (
	"os/exec"
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
