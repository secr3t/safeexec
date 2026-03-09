//go:build linux || darwin

package main

import (
	"os/signal"
	"syscall"
)

func ignoreSignals() {
	signal.Ignore(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
}

func killProcessGroup(pgid int) {
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
}
