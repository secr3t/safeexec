//go:build linux

package safeexec

import "syscall"

func setPlatformSpecificSysProcAttr(attr *syscall.SysProcAttr) {
	// SIGKILL to child when parent dies.
	attr.Pdeathsig = syscall.SIGKILL
}
