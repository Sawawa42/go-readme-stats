package main

import (
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/Sawawa42/go-readme-stats/internal/github"
	"github.com/Sawawa42/go-readme-stats/internal/gqlclient"
	"github.com/Sawawa42/go-readme-stats/internal/option"
	"github.com/joho/godotenv"
	"strings"
)

type LanguageStats struct {
	Name      string `json:"name"`
	TotalSize int    `json:"totalSize"`
	Color     string `json:"color"`
}

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

	if opts.Help {
		opts.FlagSet.Usage()
		return
	}

	client := gqlclient.NewClient("https://api.github.com/graphql")
	req, err := client.NewRequest(github.RepositoriesQuery)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	var respData github.RepositoriesResponse
	err = client.Do(req, &respData)
	if err != nil {
		fmt.Println("Error executing request:", err)
		os.Exit(1)
	}

	statsmap := make(map[string]*LanguageStats)

	for _, repo := range respData.Viewer.Repositories.Nodes {
		for _, langEdge := range repo.Languages.Edges {
			langName := langEdge.Node.Name
			if _, exists := statsmap[langName]; !exists {
				statsmap[langName] = &LanguageStats{
					Name:      langName,
					TotalSize: 0,
					Color:     langEdge.Node.Color,
				}
			}
			statsmap[langName].TotalSize += langEdge.Size
		}
	}

	var stats []LanguageStats
	for _, stat := range statsmap {
		if slices.Contains(opts.ExcludePatterns, stat.Name) {
			continue
		}
		stats = append(stats, *stat)
	}

	fmt.Println(ForatLanguageStats(stats))
}

func ForatLanguageStats(stats []LanguageStats) string {
	var total int
	for _, stat := range stats {
		total += stat.TotalSize
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalSize > stats[j].TotalSize
	})

	var builder strings.Builder
	for _, stat := range stats {
		const barLength = 20
		percentage := float64(stat.TotalSize) / float64(total) * 100
		filled := int(percentage / 100 * barLength)
		sizeKB := float64(stat.TotalSize) / 1024
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
		fmt.Fprintf(&builder, "%-12s %s %5.1f%% (%6.1f KB)\n", stat.Name, bar, percentage, sizeKB)
	}
	return builder.String()
}
