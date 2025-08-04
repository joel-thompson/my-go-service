package main

import (
	"os"

	"github.com/joel-thompson/my-go-service/cmd/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
