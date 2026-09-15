package github

// ownerAffiliationsで指定した関係を持つリポジトリを取得するクエリ(ページネーション対応)
// OWNER: 自分が所有, ORGANIZATION_MEMBER: 所属するOrganizationのもの, COLLABORATOR: コラボレーターとして参加
const RepositoriesQuery = `
	query($cursor: String, $ownerAffiliations: [RepositoryAffiliation]) {
		viewer {
			repositories(first: 100, after: $cursor, ownerAffiliations: $ownerAffiliations, isFork: false) {
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
