// Package main provides the entry point for the Godot Language Server (gdls).
package main

import (
	"fmt"
	"os"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"

	"github.com/andresperezl/gdls/internal/lsp"
)

const name = "gdls"

// Build-time variables injected via -ldflags.
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	// Handle version flag
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("%s %s\n", name, version)
			fmt.Printf("commit: %s\n", commit)
			fmt.Printf("built:  %s\n", buildTime)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		}
	}

	// Configure logging - verbosity can be increased for debugging
	commonlog.Configure(1, nil)

	server := lsp.NewServer(name, version)

	// Run the server on stdio
	if err := server.RunStdio(); err != nil {
		commonlog.GetLogger(name).Errorf("Server error: %v", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`%s - Godot Language Server

A Language Server Protocol implementation for Godot files.

Supported file types:
  - Text Scene files (.tscn, .escn)
  - Shader files (.gdshader, .gdshaderinc)

Usage:
  %s [options]

Options:
  -v, --version    Print version information
  -h, --help       Print this help message

The server communicates via stdio using the Language Server Protocol.
`, name, name)
}
