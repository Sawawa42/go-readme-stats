package github

type RepositoriesResponse struct {
	Viewer struct {
		Login        string `json:"login"`
		Repositories struct {
			Nodes []struct {
				Name      string `json:"name"`
				Languages struct {
					Edges []struct {
						Size int `json:"size"`
						Node struct {
							Name  string `json:"name"`
							Color string `json:"color"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"languages"`
			} `json:"nodes"`
			PageInfo struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"repositories"`
	} `json:"viewer"`
}
