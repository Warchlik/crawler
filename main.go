package main

import (
	"fmt"
	"os"

	process "siteprobe/internal/process"
)

func main() {
	if err := process.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
