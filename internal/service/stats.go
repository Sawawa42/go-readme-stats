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

	repos, err := fetchAllRepos(client)
	if err != nil {
		return nil, err
	}

	statsmap := aggregateStats(repos)
	stats := filterAndSortStats(statsmap, excludePatterns)
	return stats, nil
}

func fetchAllRepos(client *gqlclient.Client) ([]github.RepoNode, error) {
	seen := make(map[string]struct{})
	var all []github.RepoNode

	ownerRepos, err := fetchOwnerRepos(client)
	if err != nil {
		return nil, err
	}
	for _, r := range ownerRepos {
		if _, ok := seen[r.NameWithOwner]; ok {
			continue
		}
		seen[r.NameWithOwner] = struct{}{}
		all = append(all, r)
	}

	contributedRepos, err := fetchContributedRepos(client)
	if err != nil {
		return nil, err
	}
	for _, r := range contributedRepos {
		if _, ok := seen[r.NameWithOwner]; ok {
			continue
		}
		seen[r.NameWithOwner] = struct{}{}
		all = append(all, r)
	}

	return all, nil
}

func fetchOwnerRepos(client *gqlclient.Client) ([]github.RepoNode, error) {
	var repos []github.RepoNode
	var cursor any = nil
	for {
		req, err := client.NewRequest(github.OwnerRepositoriesQuery, map[string]any{"cursor": cursor})
		if err != nil {
			return nil, err
		}
		var resp github.OwnerRepositoriesResponse
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

func fetchContributedRepos(client *gqlclient.Client) ([]github.RepoNode, error) {
	var repos []github.RepoNode
	var cursor any = nil
	for {
		req, err := client.NewRequest(github.ContributedRepositoriesQuery, map[string]any{"cursor": cursor})
		if err != nil {
			return nil, err
		}
		var resp github.ContributedRepositoriesResponse
		if err := client.Do(req, &resp); err != nil {
			return nil, err
		}
		repos = append(repos, resp.Viewer.RepositoriesContributedTo.Nodes...)
		if !resp.Viewer.RepositoriesContributedTo.PageInfo.HasNextPage {
			break
		}
		cursor = resp.Viewer.RepositoriesContributedTo.PageInfo.EndCursor
	}
	return repos, nil
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

	// サイズでソート(降順)
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalSize > stats[j].TotalSize
	})
	return stats
}
