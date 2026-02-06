package main

import (
	"fmt"
	"github.com/Sawawa42/go-readme-stats/internal/gqlclient"
	"github.com/Sawawa42/go-readme-stats/internal/option"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	opts, err := option.Parse(os.Args)
	if err != nil {
		fmt.Println("Error parsing options:", err)
		os.Exit(1)
	}
	fmt.Printf("Parsed Options: %+v\n", opts)

	// クエリ例
	client := gqlclient.NewClient("https://api.github.com/graphql")
	query := `
	{
		viewer {
			login
		}
	}`
	req, err := client.NewRequest(query)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	var respData map[string]interface{}
	err = client.Do(req, &respData)
	if err != nil {
		fmt.Println("Error executing request:", err)
		os.Exit(1)
	}

	fmt.Printf("Response Data: %+v\n", respData)
}
