package safeexec

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed internal/assets/*
var assets embed.FS

var extractedWatchdogPath string

func init() {
	// Identify the correct binary for the current platform
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	binName := fmt.Sprintf("watchdog-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
	assetPath := "internal/assets/" + binName

	// Get cache directory
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	appCacheDir := filepath.Join(cacheDir, "safeexec")
	if err := os.MkdirAll(appCacheDir, 0755); err != nil {
		return
	}

	targetPath := filepath.Join(appCacheDir, binName)
	extractedWatchdogPath = targetPath

	// Check if already exists
	if _, err := os.Stat(targetPath); err == nil {
		// Already exists, just ensure it's executable
		_ = os.Chmod(targetPath, 0755)
		return
	}

	// Extract from embed.FS
	src, err := assets.Open(assetPath)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return
	}
}

// Cmd is a wrapper around exec.Cmd that provides group killing and parent watching capabilities.
type Cmd struct {
	*exec.Cmd
	Process      *Process
	watchdogCmd  *exec.Cmd
	pipeWriter   *os.File
	WatchdogPath string // Path to the watchdog binary (overrides automatic extraction)
}

// Process is a wrapper around os.Process that overrides Kill to kill the entire process group.
type Process struct {
	*os.Process
}

// Kill terminates the process and all its children by killing the process group.
func (p *Process) Kill() error {
	return killProcess(p)
}

// Command returns a Cmd struct to execute the named program with the given arguments.
func Command(name string, arg ...string) *Cmd {
	cmd := exec.Command(name, arg...)
	setupCmd(cmd)
	return &Cmd{Cmd: cmd}
}

// CommandContext is like Command but includes a context.
func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	setupCmd(cmd)
	// Ensure that when context is cancelled, the entire process group is killed.
	cmd.Cancel = func() error {
		return killProcess(&Process{Process: cmd.Process})
	}
	return &Cmd{Cmd: cmd}
}

// Start starts the specified command but does not wait for it to complete.
func (c *Cmd) Start() error {
	if err := c.Cmd.Start(); err != nil {
		return err
	}
	c.Process = &Process{Process: c.Cmd.Process}

	// Start watchdog
	if err := c.startWatchdog(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to start watchdog: %v\n", err)
	}

	return nil
}

func (c *Cmd) startWatchdog() error {
	path := c.WatchdogPath
	if path == "" {
		path = extractedWatchdogPath
	}

	if path == "" {
		return fmt.Errorf("watchdog binary not found or failed to extract")
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("watchdog binary not found at %s", path)
	}

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	c.pipeWriter = w

	watchdog := exec.Command(path, "-pgid", fmt.Sprintf("%d", c.Cmd.Process.Pid))
	setupWatchdogCmd(watchdog)
	watchdog.Stdin = r
	if err := watchdog.Start(); err != nil {
		r.Close()
		w.Close()
		return err
	}

	c.watchdogCmd = watchdog
	r.Close()

	return nil
}

// Wait waits for the command to exit and waits for any copying to stdin or copying from stdout or stderr to complete.
func (c *Cmd) Wait() error {
	err := c.Cmd.Wait()
	if c.pipeWriter != nil {
		c.pipeWriter.Close()
	}
	if c.watchdogCmd != nil {
		_ = c.watchdogCmd.Wait()
	}
	return err
}

// Run starts the specified command and waits for it to complete.
func (c *Cmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}
