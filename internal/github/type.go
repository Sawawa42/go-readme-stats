package github

type LanguageEdge struct {
	Size int `json:"size"`
	Node struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"node"`
}

type RepoNode struct {
	NameWithOwner string `json:"nameWithOwner"`
	Languages     struct {
		Edges []LanguageEdge `json:"edges"`
	} `json:"languages"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type OwnerRepositoriesResponse struct {
	Viewer struct {
		Login        string `json:"login"`
		Repositories struct {
			Nodes    []RepoNode `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"repositories"`
	} `json:"viewer"`
}

type ContributedRepositoriesResponse struct {
	Viewer struct {
		RepositoriesContributedTo struct {
			Nodes    []RepoNode `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"repositoriesContributedTo"`
	} `json:"viewer"`
}
