package svg

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sawawa42/go-readme-stats/internal/model"
)

type SVGConfig struct {
	Width      int
	Height     int
	Padding    int
	BarHeight  int
	BarSpacing int
	FontSize   int
	Title      string
	UpdatedAt  time.Time // ゼロ値なら更新日を表示しない
}

func DefaultConfig() SVGConfig {
	return SVGConfig{
		Width:      600,
		Height:     0,
		Padding:    20,
		BarHeight:  25,
		BarSpacing: 8,
		FontSize:   14,
		Title:      "Most Used Languages",
	}
}

func Generate(stats []model.LanguageStats, config SVGConfig) string {
	if len(stats) == 0 {
		return ""
	}

	var maxPercentage float64
	for _, stat := range stats {
		maxPercentage = max(maxPercentage, stat.Percentage)
	}

	// 最長の言語名の長さを計算 (1文字あたり約8ピクセルと仮定)
	maxNameLength := 0
	for _, stat := range stats {
		if len(stat.Name) > maxNameLength {
			maxNameLength = len(stat.Name)
		}
	}
	nameWidth := maxNameLength * 8 + 10 // 余白を追加

	titleHeight := 40
	config.Height = titleHeight + len(stats)*(config.BarHeight+config.BarSpacing) + config.Padding*2

	var builder strings.Builder

	// <svg>
	fmt.Fprintf(&builder, `<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`,
		config.Width, config.Height)
	builder.WriteString("\n")

	// <rect> 背景
	fmt.Fprintf(&builder, `  <rect width="%d" height="%d" fill="#ffffff" rx="10"/>`,
		config.Width, config.Height)
	builder.WriteString("\n")

	// <text> タイトル
	fmt.Fprintf(&builder, `  <text x="%d" y="%d" font-family="'Segoe UI', Ubuntu, Sans-Serif" font-size="18" font-weight="bold" fill="#2f80ed">%s</text>`,
		config.Padding, config.Padding + 20, config.Title)
	builder.WriteString("\n")

	// <text> 更新日 (タイトル行の右端)
	if !config.UpdatedAt.IsZero() {
		fmt.Fprintf(&builder, `  <text x="%d" y="%d" text-anchor="end" font-family="'Segoe UI', Ubuntu, Sans-Serif" font-size="11" fill="#888">Updated %s</text>`,
			config.Width - config.Padding, config.Padding + 20, config.UpdatedAt.Format("2006-01-02 UTC"))
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	y := titleHeight + config.Padding

	// レイアウト計算: パーセンテージ表示に必要な幅
	// monospaceフォールバック時 "100.0%" で約45px + バーとの間隔
	infoWidth := 60
	barX := config.Padding + nameWidth
	maxBarWidth := config.Width - config.Padding * 2 - nameWidth - infoWidth

	// 各言語のバーを描画
	for _, stat := range stats {
		// 最大の言語をバーの最大幅とし、他はその比率で描画する
		barWidth := 0
		if maxPercentage > 0 {
			barWidth = int(float64(maxBarWidth) * stat.Percentage / maxPercentage)
		}

		fmt.Fprintf(&builder, `  <text x="%d" y="%d" font-family="'Segoe UI', Ubuntu, monospace" font-size="%d" fill="#333" font-weight="600">%s</text>`,
			config.Padding, y + config.BarHeight / 2 + 5, config.FontSize, stat.Name)
		builder.WriteString("\n")

		fmt.Fprintf(&builder, `  <rect x="%d" y="%d" width="%d" height="%d" fill="#e1e4e8" rx="3"/>`,
			barX, y, maxBarWidth, config.BarHeight - 5)
		builder.WriteString("\n")

		color := stat.Color
		if color == "" {
			color = "#586069"
		}
		if barWidth > 0 {
			fmt.Fprintf(&builder, `  <rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="3"/>`,
				barX, y, barWidth, config.BarHeight - 5, color)
			builder.WriteString("\n")
		}

		textX := config.Width - config.Padding
		fmt.Fprintf(&builder, `  <text x="%d" y="%d" text-anchor="end" font-family="'Segoe UI', Ubuntu, monospace" font-size="%d" fill="#666">%.1f%%</text>`,
			textX, y + config.BarHeight / 2 + 5, config.FontSize - 2, stat.Percentage)
		builder.WriteString("\n\n")

		y += config.BarHeight + config.BarSpacing
	}

	// </svg>
	builder.WriteString("</svg>")

	return builder.String()
}
