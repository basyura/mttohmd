package converter

import (
	"fmt"
	"regexp"
	"strings"

	"mttohmd/entry"
)

// ToMarkdown エントリーをHatena Blog形式のMarkdownに変換
func ToMarkdown(e entry.Entry) string {
	var md strings.Builder

	// はてなブログ用のメタデータ（フロントマター形式）
	md.WriteString("---\n")
	md.WriteString("Title: ")
	md.WriteString(e.Title)
	md.WriteString("\n")

	if e.Category != "" {
		md.WriteString("Category:\n")
		categories := strings.Split(e.Category, ",")
		for _, cat := range categories {
			cat = strings.TrimSpace(cat)
			if cat != "" {
				md.WriteString("- ")
				md.WriteString(cat)
				md.WriteString("\n")
			}
		}
	}

	if e.Date != "" {
		md.WriteString("Date: ")
		md.WriteString(e.Date)
		md.WriteString("\n")
	}

	md.WriteString("---\n\n")

	// 記事本文
	// MovableType形式からMarkdown/HTML混在形式へ変換
	body := convertMTToMarkdown(e.Body)
	md.WriteString(body)

	// 画像がある場合は記事の最後に追加
	if e.ImageURL != "" {
		md.WriteString("\n\n")
		md.WriteString("![")
		md.WriteString(e.Title)
		md.WriteString("](")
		md.WriteString(e.ImageURL)
		md.WriteString(")")
	}

	return md.String()
}

// convertMTToMarkdown はMovableType形式のテキストをMarkdown形式に変換する
func convertMTToMarkdown(body string) string {
	// 基本的な変換処理
	result := body

	// 改行の正規化
	result = strings.ReplaceAll(result, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")

	// MovableTypeの基本的なHTMLタグをMarkdownに変換
	result = convertHTMLToMarkdown(result)

	// 空行の整理
	result = regexp.MustCompile(`\n\n+`).ReplaceAllString(result, "\n\n")
	result = strings.TrimSpace(result)

	return result
}

// convertHTMLToMarkdown は基本的なHTMLタグをMarkdown記法に変換する
func convertHTMLToMarkdown(text string) string {
	result := text

	// <br> タグを改行に変換
	result = regexp.MustCompile(`<br\s*/?>|<BR\s*/?>`).ReplaceAllString(result, "\n")

	// <p> タグを段落に変換
	result = regexp.MustCompile(`<p[^>]*>`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`</p>`).ReplaceAllString(result, "\n\n")

	// <strong> や <b> タグを太字に変換
	result = regexp.MustCompile(`<(?:strong|b)[^>]*>(.*?)</(?:strong|b)>`).ReplaceAllString(result, "**$1**")

	// <em> や <i> タグを斜体に変換
	result = regexp.MustCompile(`<(?:em|i)[^>]*>(.*?)</(?:em|i)>`).ReplaceAllString(result, "*$1*")

	// <a> タグをMarkdownリンクに変換
	result = regexp.MustCompile(`<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`).ReplaceAllString(result, "[$2]($1)")

	// <img> タグをMarkdown画像に変換
	result = regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*(?:alt=["']([^"']*)["'][^>]*)?/?>`).ReplaceAllString(result, "![$2]($1)")

	// <h1> から <h6> タグをMarkdownヘッダーに変換
	for i := 1; i <= 6; i++ {
		headerTag := fmt.Sprintf("h%d", i)
		headerMark := strings.Repeat("#", i)
		pattern := fmt.Sprintf(`<%s[^>]*>(.*?)</%s>`, headerTag, headerTag)
		replacement := fmt.Sprintf("%s $1", headerMark)
		result = regexp.MustCompile(pattern).ReplaceAllString(result, replacement)
	}

	// <blockquote> タグを引用に変換
	result = regexp.MustCompile(`<blockquote[^>]*>(.*?)</blockquote>`).ReplaceAllStringFunc(result, func(match string) string {
		content := regexp.MustCompile(`<blockquote[^>]*>(.*?)</blockquote>`).ReplaceAllString(match, "$1")
		lines := strings.Split(strings.TrimSpace(content), "\n")
		var quotedLines []string
		for _, line := range lines {
			quotedLines = append(quotedLines, "> "+strings.TrimSpace(line))
		}
		return strings.Join(quotedLines, "\n")
	})

	// <ul> と <li> タグをMarkdownリストに変換
	result = regexp.MustCompile(`<ul[^>]*>`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`</ul>`).ReplaceAllString(result, "\n")
	result = regexp.MustCompile(`<li[^>]*>(.*?)</li>`).ReplaceAllString(result, "- $1")

	// <ol> タグを番号付きリストに変換（簡易版）
	result = regexp.MustCompile(`<ol[^>]*>`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`</ol>`).ReplaceAllString(result, "\n")

	return result
}