# go-readme-stats

Github上のリポジトリの言語統計をSVGとして自動生成するCLIツール

<img src="./generated/language-stats.svg" />

## Quick Start

```sh
export GITHUB_TOKEN=ghp_xxxx...
go run ./cmd/main.go -x=Kotlin,PHP,HTML,CSS,PLpgSQL,Nix,Dockerfile,Dart
```

## Usage

```
go run ./cmd/main.go [OPTIONS]

OPTIONS:
  -x, --exclude  除外する言語(カンマ区切り)
  -h, --help     ヘルプを表示
```

## Features

- Github GraphQL APIで情報を取得
- SVG形式で出力
