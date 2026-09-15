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
	"time"
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

	stats, repoStats, err := service.FetchAndBuildStats(service.Params{
		Roles:           opts.Roles,
		ExcludeRepos:    opts.ExcludeRepos,
		IncludePatterns: opts.IncludePatterns,
		SizeWeight:      opts.SizeWeight,
		CountWeight:     opts.CountWeight,
	})
	if err != nil {
		return fmt.Errorf("error fetching and building stats: %w", err)
	}

	fmt.Print(printStatsToConsole(stats))
	if opts.Verbose {
		fmt.Print(printRepoStatsToConsole(repoStats))
	}

	config := svg.DefaultConfig()
	config.UpdatedAt = time.Now().UTC()
	svgOutput := svg.Generate(stats, config)

	err = saveSVGToFile(svgOutput, "./generated/language-stats.svg")
	if err != nil {
		return fmt.Errorf("error saving SVG file: %w", err)
	}

	return nil
}

// 指標に加え、検証用に元になったバイト数とリポジトリ数も出力する
func printStatsToConsole(stats []model.LanguageStats) string {
	var builder strings.Builder
	for _, stat := range stats {
		const barLength = 20
		filled := int(stat.Percentage / 100 * barLength)
		sizeKB := float64(stat.TotalSize) / 1024
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
		fmt.Fprintf(&builder, "%-12s %s %5.1f%% (%6.1f KB, %2d repos)\n", stat.Name, bar, stat.Percentage, sizeKB, stat.RepoCount)
	}
	return builder.String()
}

func printRepoStatsToConsole(repoStats []model.RepoStats) string {
	var total int
	for _, rs := range repoStats {
		total += rs.TotalSize
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "\n--- Breakdown by repository (%d repos) ---\n", len(repoStats))
	for _, rs := range repoStats {
		var langs []string
		for _, lang := range rs.Languages {
			langs = append(langs, fmt.Sprintf("%s %.1f KB", lang.Name, float64(lang.TotalSize)/1024))
		}
		fmt.Fprintf(&builder, "%-45s %8.1f KB %5.1f%%  %s\n",
			rs.NameWithOwner, float64(rs.TotalSize)/1024,
			float64(rs.TotalSize)/float64(total)*100, strings.Join(langs, ", "))
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
