package app

import (
	"fmt"
	"github.com/Sawawa42/go-readme-stats/internal/model"
	"github.com/Sawawa42/go-readme-stats/internal/option"
	"github.com/Sawawa42/go-readme-stats/internal/service"
	"github.com/Sawawa42/go-readme-stats/internal/svg"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func Run() error {
	_ = godotenv.Load()

	opts, err := option.Parse(os.Args)
	if err != nil {
		return fmt.Errorf("error parsing options: %w", err)
	}

	if opts.Help {
		opts.FlagSet.Usage()
		return nil
	}

	stats, err := service.FetchAndBuildStats(opts.ExcludePatterns)
	if err != nil {
		return fmt.Errorf("error fetching and building stats: %w", err)
	}

	fmt.Print(printStatsToConsole(stats))

	config := svg.DefaultConfig()
	svgOutput := svg.Generate(stats, config)

	err = saveSVGToFile(svgOutput, "./generated/language-stats.svg")
	if err != nil {
		return fmt.Errorf("error saving SVG file: %w", err)
	}

	return nil
}

func printStatsToConsole(stats []model.LanguageStats) string {
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

func saveSVGToFile(svgContent, filePath string) error {
	err := os.WriteFile(filePath, []byte(svgContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing SVG file: %w", err)
	}
	fmt.Fprintln(os.Stderr, "SVG file generated:", filePath)
	return nil
}
