package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/Sawawa42/go-readme-stats/internal/github"
	"github.com/Sawawa42/go-readme-stats/internal/gqlclient"
	"github.com/Sawawa42/go-readme-stats/internal/option"
	"github.com/Sawawa42/go-readme-stats/internal/svg"
	"github.com/joho/godotenv"
	"github.com/Sawawa42/go-readme-stats/internal/model"
)

func main() {
	_ = godotenv.Load()

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

	statsmap := make(map[string]*model.LanguageStats)

	for _, repo := range respData.Viewer.Repositories.Nodes {
		for _, langEdge := range repo.Languages.Edges {
			langName := langEdge.Node.Name
			if _, exists := statsmap[langName]; !exists {
				statsmap[langName] = &model.LanguageStats{
					Name:      langName,
					TotalSize: 0,
					Color:     langEdge.Node.Color,
				}
			}
			statsmap[langName].TotalSize += langEdge.Size
		}
	}

	var stats []model.LanguageStats
	for _, stat := range statsmap {
		if slices.Contains(opts.ExcludePatterns, stat.Name) {
			continue
		}
		stats = append(stats, *stat)
	}

	// サイズでソート（降順）
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalSize > stats[j].TotalSize
	})

	// テキスト出力（既存）
	fmt.Println(FormatLanguageStats(stats))

	// SVG出力
	svgStats := make([]model.LanguageStats, len(stats))
	for i, stat := range stats {
		svgStats[i] = model.LanguageStats{
			Name:      stat.Name,
			TotalSize: stat.TotalSize,
			Color:     stat.Color,
		}
	}

	config := svg.DefaultConfig()
	svgOutput := svg.Generate(svgStats, config)

	// SVGファイルに書き込む
	err = os.WriteFile("./generated/language-stats.svg", []byte(svgOutput), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing SVG file:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "SVG file generated: language-stats.svg")
}

func FormatLanguageStats(stats []model.LanguageStats) string {
	var total int
	for _, stat := range stats {
		total += stat.TotalSize
	}

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
