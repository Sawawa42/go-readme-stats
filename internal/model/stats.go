package model

type LanguageStats struct {
	Name       string  `json:"name"`
	TotalSize  int     `json:"totalSize"`
	RepoCount  int     `json:"repoCount"`  // この言語を含むリポジトリ数
	Percentage float64 `json:"percentage"` // よく使っている度合いの指標(表示対象の言語内で合計100)
	Color      string  `json:"color"`
}

// リポジトリ単位の内訳(集計値の検証用)
type RepoStats struct {
	NameWithOwner string          `json:"nameWithOwner"`
	TotalSize     int             `json:"totalSize"`
	Languages     []LanguageStats `json:"languages"`
}
