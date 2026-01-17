package main

import (
	"fmt"
	"github.com/Sawawa42/go-readme-stats/internal/option"
	"os"
)

func main() {
	opts, err := option.Parse(os.Args)
	if err != nil {
		fmt.Println("Error parsing options:", err)
		os.Exit(1)
	}
	fmt.Printf("Parsed Options: %+v\n", opts)
}
