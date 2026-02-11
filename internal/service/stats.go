package service

import (
	"slices"
	"sort"
	"github.com/Sawawa42/go-readme-stats/internal/github"
	"github.com/Sawawa42/go-readme-stats/internal/gqlclient"
	"github.com/Sawawa42/go-readme-stats/internal/model"
)

func FetchAndBuildStats(excludePatterns []string) ([]model.LanguageStats, error) {
	client := gqlclient.NewClient("https://api.github.com/graphql")
	req, err := client.NewRequest(github.RepositoriesQuery)
	if err != nil {
		return nil, err
	}

	var respData github.RepositoriesResponse
	err = client.Do(req, &respData)
	if err != nil {
		return nil, err
	}

	statsmap := aggregateStats(respData)
	stats := filterAndSortStats(statsmap, excludePatterns)
	return stats, nil
}

func aggregateStats(respData github.RepositoriesResponse) map[string]*model.LanguageStats {
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
	return statsmap
}

func filterAndSortStats(statsmap map[string]*model.LanguageStats, excludePatterns []string) []model.LanguageStats {
	var stats []model.LanguageStats
	for _, stat := range statsmap {
		if slices.Contains(excludePatterns, stat.Name) {
			continue
		}
		stats = append(stats, *stat)
	}

	// サイズでソート（降順）
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalSize > stats[j].TotalSize
	})
	return stats
}
