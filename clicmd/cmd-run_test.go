package clicmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestCmdRunCreatesOutputDir(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "rrdreader-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Define a non-existent output directory path
	outputDir := filepath.Join(tmpDir, "output")
	rootDir := filepath.Join(tmpDir, "collectd")

	// Create dummy directories and files
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		t.Fatalf("Failed to create root dir: %v", err)
	}

	// Create a new cobra command and set the flags
	cmd := &cobra.Command{}
	cmd.Flags().String("output", outputDir, "")
	cmd.Flags().String("rootdir", rootDir, "")
	// we need to set address and port to avoid issues with the listener in test
	cmd.Flags().String("address", "127.0.0.1", "")
	cmd.Flags().Int("port", 0, "") // use port 0 to let the OS choose a free port

	// Run CmdRun in a separate goroutine
	go func() {
		// a trick to avoid the infinite loop in the test
		// we are only testing the directory creation
		// which happens before the loop
		defer func() {
			if r := recover(); r != nil {
				// ignore the panic from log.Fatalf
			}
		}()
		CmdRun(cmd, []string{})
	}()

	// Wait for a short period to allow directory creation
	time.Sleep(200 * time.Millisecond)

	// Check if the output directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Errorf("Expected output directory '%s' to be created, but it was not", outputDir)
	}
}
