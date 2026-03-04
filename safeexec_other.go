//go:build !linux

package safeexec

import "syscall"

func setPlatformSpecificSysProcAttr(attr *syscall.SysProcAttr) {
	// Other platforms (like Darwin) don't support Pdeathsig in the same way.
	// We rely on the context cancellation or manual Kill() calling for these platforms.
}
