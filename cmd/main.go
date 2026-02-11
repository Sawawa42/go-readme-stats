package main

import (
	"fmt"
	"github.com/Sawawa42/go-readme-stats/internal/app"
	"os"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
