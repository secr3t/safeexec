//go:build windows

package main

import (
	"os"
)

func ignoreSignals() {
	// Windows doesn't have SIGHUP in the same way.
	// We could ignore other signals if needed, but for now we'll do nothing.
}

func killProcessGroup(pgid int) {
	// On Windows, pgid is likely just a PID since we don't have process groups.
	// We try to kill the process.
	if p, err := os.FindProcess(pgid); err == nil {
		_ = p.Kill()
	}
}
