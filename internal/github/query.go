package github

// struct->queryの自動生成はできそうだが、複雑なクエリ(repositories(first: 100))に対応するのは難しそう
const RepositoriesQuery = `
	query {
		viewer {
			login
			repositories(first: 100, ownerAffiliations: OWNER, isFork: false) {
				nodes {
					name
					languages(first: 20, orderBy: {field: SIZE, direction: DESC}) {
						edges {
							size
							node {
								name
								color
							}
						}
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	}`
