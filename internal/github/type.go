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

type RepositoriesResponse struct {
	Viewer struct {
		Repositories struct {
			Nodes    []RepoNode `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"repositories"`
	} `json:"viewer"`
}
