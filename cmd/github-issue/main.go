package main

import (
	"os"

	"github.com/shichao402/github-issue-pack/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
