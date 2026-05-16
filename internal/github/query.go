package github

// 自分が所有するリポジトリを取得するクエリ(ページネーション対応)
const OwnerRepositoriesQuery = `
	query($cursor: String) {
		viewer {
			login
			repositories(first: 100, after: $cursor, ownerAffiliations: OWNER, isFork: false) {
				nodes {
					nameWithOwner
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

// Organization含む、自分が貢献したリポジトリを取得するクエリ(ページネーション対応)
// includeUserRepositories: false で OWNER クエリと重複する自リポジトリを除外
const ContributedRepositoriesQuery = `
	query($cursor: String) {
		viewer {
			repositoriesContributedTo(first: 100, after: $cursor, contributionTypes: [COMMIT, PULL_REQUEST], includeUserRepositories: false) {
				nodes {
					nameWithOwner
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
