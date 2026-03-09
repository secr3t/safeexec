package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/secr3t/safeexec"
)

func main() {
	fmt.Println("Starting Watchdog & Group Kill Test...")

	// 1. Test Process Group Killing with Kill()
	fmt.Println("\n--- Test 1: Manual Process Group Killing ---")
	cmd := safeexec.Command("sh", "-c", "sleep 100 & sleep 100")
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting: %v\n", err)
		return
	}

	pid := cmd.Process.Pid
	fmt.Printf("Started process group leader with PID: %d\n", pid)
	printProcessTree(pid)

	time.Sleep(1 * time.Second)
	fmt.Println("Calling cmd.Process.Kill()...")
	if err := cmd.Process.Kill(); err != nil {
		fmt.Printf("Error killing: %v\n", err)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Checking if processes are gone...")
	printProcessTree(pid)

	// 2. Test Watchdog (Simulated by closing pipe or waiting for exit)
	fmt.Println("\n--- Test 2: Watchdog Verification ---")
	fmt.Println("This test will start a long-running process and then this program will exit.")
	fmt.Println("The Watchdog should then kill the process group.")

	cmdW := safeexec.Command("sh", "-c", "sleep 100 & sleep 100")
	// WatchdogPath is now automatically handled via init() and embed

	if err := cmdW.Start(); err != nil {
		fmt.Printf("Error starting with watchdog: %v\n", err)
		return
	}

	wPid := cmdW.Process.Pid
	fmt.Printf("Started process group %d with watchdog.\n", wPid)
	printProcessTree(wPid)

	// In a real scenario, if this program exits now, the watchdog will kill wPid.
	// For this test, we'll just wait a bit and then explain how to verify.
	fmt.Println("Verify manually by running 'ps -g <PGID>' after this program exits.")
	fmt.Printf("PGID to check: %d\n", wPid)

	time.Sleep(1 * time.Second)
	fmt.Println("Test finished. Exiting parent now...")
}

func printProcessTree(pid int) {
	out, _ := exec.Command("ps", "-o", "pid,ppid,pgid,comm", "-g", fmt.Sprintf("%d", pid)).Output()
	if len(out) == 0 {
		fmt.Println("No processes found in group.")
		return
	}
	fmt.Printf("Processes in group %d:\n%s", pid, string(out))
}
