package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pgid := flag.Int("pgid", 0, "Process Group ID to kill on parent death")
	flag.Parse()

	if *pgid <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: watchdog -pgid <PGID>\n")
		os.Exit(1)
	}

	// We ignore common signals to stay alive until the parent bridge breaks
	signal.Ignore(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// Watchdog logic: Wait for Stdin to close.
	// When the parent dies, the write end of the pipe (inherited as Stdin) will close.
	_, _ = io.Copy(io.Discard, os.Stdin)

	// Stdin closed -> Parent is gone.
	// Kill the process group.
	_ = syscall.Kill(-*pgid, syscall.SIGKILL)

	os.Exit(0)
}
