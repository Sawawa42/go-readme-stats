package main

import (
	"fmt"
	"os"
	"github.com/Sawawa42/go-readme-stats/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
