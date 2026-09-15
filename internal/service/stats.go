package service

import (
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/Sawawa42/go-readme-stats/internal/github"
	"github.com/Sawawa42/go-readme-stats/internal/gqlclient"
	"github.com/Sawawa42/go-readme-stats/internal/model"
)

type Params struct {
	Roles           []string // 対象とするリポジトリとの関係(ownerAffiliations)
	ExcludeRepos    []string // 除外するリポジトリ
	IncludePatterns []string // 表示対象とする言語
	SizeWeight      float64  // 指標におけるバイト数の重み
	CountWeight     float64  // 指標におけるリポジトリ数の重み
}

func FetchAndBuildStats(params Params) ([]model.LanguageStats, []model.RepoStats, error) {
	client := gqlclient.NewClient("https://api.github.com/graphql")

	repos, err := fetchRepos(client, params.Roles)
	if err != nil {
		return nil, nil, err
	}
	repos = excludeRepositories(repos, params.ExcludeRepos)

	statsmap := aggregateStats(repos)
	stats := filterAndSortStats(statsmap, params.IncludePatterns, params.SizeWeight, params.CountWeight)
	repoStats := buildRepoStats(repos, params.IncludePatterns)
	return stats, repoStats, nil
}

func fetchRepos(client *gqlclient.Client, roles []string) ([]github.RepoNode, error) {
	var repos []github.RepoNode
	var cursor any = nil
	for {
		req, err := client.NewRequest(github.RepositoriesQuery, map[string]any{"cursor": cursor, "ownerAffiliations": roles})
		if err != nil {
			return nil, err
		}
		var resp github.RepositoriesResponse
		if err := client.Do(req, &resp); err != nil {
			return nil, err
		}
		repos = append(repos, resp.Viewer.Repositories.Nodes...)
		if !resp.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}
		cursor = resp.Viewer.Repositories.PageInfo.EndCursor
	}
	return repos, nil
}

// 除外指定に一致するリポジトリを取り除く
// "owner/name" なら完全一致、"name" だけならどのオーナーのものでも一致
func excludeRepositories(repos []github.RepoNode, excludeRepos []string) []github.RepoNode {
	if len(excludeRepos) == 0 {
		return repos
	}
	return slices.DeleteFunc(repos, func(r github.RepoNode) bool {
		_, name, _ := strings.Cut(r.NameWithOwner, "/")
		return slices.Contains(excludeRepos, r.NameWithOwner) || slices.Contains(excludeRepos, name)
	})
}

func aggregateStats(repos []github.RepoNode) map[string]*model.LanguageStats {
	statsmap := make(map[string]*model.LanguageStats)

	for _, repo := range repos {
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
			statsmap[langName].RepoCount++
		}
	}
	return statsmap
}

// リポジトリごとに対象言語のサイズをまとめる(サイズ降順)
func buildRepoStats(repos []github.RepoNode, includePatterns []string) []model.RepoStats {
	var result []model.RepoStats
	for _, repo := range repos {
		rs := model.RepoStats{NameWithOwner: repo.NameWithOwner}
		for _, langEdge := range repo.Languages.Edges {
			if len(includePatterns) > 0 && !slices.Contains(includePatterns, langEdge.Node.Name) {
				continue
			}
			rs.TotalSize += langEdge.Size
			rs.Languages = append(rs.Languages, model.LanguageStats{
				Name:      langEdge.Node.Name,
				TotalSize: langEdge.Size,
				Color:     langEdge.Node.Color,
			})
		}
		if rs.TotalSize > 0 {
			result = append(result, rs)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalSize > result[j].TotalSize
	})
	return result
}

// 表示対象の言語に絞り込み、指標の降順に並べる
// 指標 = バイト数^sizeWeight × リポジトリ数^countWeight を、表示対象の言語内で合計100%に正規化したもの
// (github-stats-extended の size_weight / count_weight と同じ考え方)
func filterAndSortStats(statsmap map[string]*model.LanguageStats, includePatterns []string, sizeWeight, countWeight float64) []model.LanguageStats {
	var stats []model.LanguageStats
	scores := make(map[string]float64)
	var totalScore float64
	for _, stat := range statsmap {
		if len(includePatterns) > 0 && !slices.Contains(includePatterns, stat.Name) {
			continue
		}
		score := math.Pow(float64(stat.TotalSize), sizeWeight) * math.Pow(float64(stat.RepoCount), countWeight)
		scores[stat.Name] = score
		totalScore += score
		stats = append(stats, *stat)
	}

	for i := range stats {
		if totalScore > 0 {
			stats[i].Percentage = scores[stats[i].Name] / totalScore * 100
		}
	}

	// 指標でソート(降順)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Percentage > stats[j].Percentage
	})
	return stats
}
