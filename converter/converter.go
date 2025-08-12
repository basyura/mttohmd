package converter

import (
	"html"
	"regexp"
	"strings"

	"mttohmd/entry"
)

var (
	// ASIN詳細タグ用の正規表現
	asinDetailWithPTagRegex = regexp.MustCompile(`(?s)<p><div class="hatena-asin-detail">.*?href="https://www\.amazon\.co\.jp/dp/([A-Z0-9]+)[^"]*".*?</div></div></p>`)
	asinDetailNestedRegex   = regexp.MustCompile(`(?s)<div class="hatena-asin-detail">.*?href="https://www\.amazon\.co\.jp/dp/([A-Z0-9]+)[^"]*".*?</div></div>`)
	asinDetailSimpleRegex   = regexp.MustCompile(`(?s)<div class="hatena-asin-detail">.*?href="https://www\.amazon\.co\.jp/dp/([A-Z0-9]+)[^"]*".*?</div>`)

	// HTMLタグ変換用の正規表現
	brTagRegex            = regexp.MustCompile(`<br\s*/?>|<BR\s*/?>`)
	pTagOpenRegex         = regexp.MustCompile(`<p[^>]*>`)
	pTagCloseRegex        = regexp.MustCompile(`</p>`)
	strongTagRegex        = regexp.MustCompile(`<(?:strong|b)[^>]*>(.*?)</(?:strong|b)>`)
	emTagRegex            = regexp.MustCompile(`<(?:em|i)[^>]*>(.*?)</(?:em|i)>`)
	codeTagRegex          = regexp.MustCompile(`<code[^>]*>(.*?)</code>`)
	aTagRegex             = regexp.MustCompile(`<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	imgWithAltRegex       = regexp.MustCompile(`<img[^>]*src=["']([^"']+)["'][^>]*alt=["']([^"']*)["'][^>]*/?>`)
	imgAltSrcRegex        = regexp.MustCompile(`<img[^>]*alt=["']([^"']*)["'][^>]*src=["']([^"']+)["'][^>]*/?>`)
	imgSimpleRegex        = regexp.MustCompile(`<img[^>]*src=["']([^"']+)["'][^>]*/?>`)
	blockquoteRegex       = regexp.MustCompile(`(?s)<blockquote[^>]*>(.*?)</blockquote>`)
	blockquoteInnerRegex  = regexp.MustCompile(`(?s)<blockquote[^>]*>(.*?)</blockquote>`)
	ulOpenRegex           = regexp.MustCompile(`<ul[^>]*>`)
	ulCloseRegex          = regexp.MustCompile(`</ul>`)
	liTagRegex            = regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
	olOpenRegex           = regexp.MustCompile(`<ol[^>]*>`)
	olCloseRegex          = regexp.MustCompile(`</ol>`)
	newlineNormalizeRegex = regexp.MustCompile(`\n\n+`)

	// ヘッダータグ用の正規表現
	h1Regex = regexp.MustCompile(`<h1[^>]*>(.*?)</h1>`)
	h2Regex = regexp.MustCompile(`<h2[^>]*>(.*?)</h2>`)
	h3Regex = regexp.MustCompile(`<h3[^>]*>(.*?)</h3>`)
	h4Regex = regexp.MustCompile(`<h4[^>]*>(.*?)</h4>`)
	h5Regex = regexp.MustCompile(`<h5[^>]*>(.*?)</h5>`)
	h6Regex = regexp.MustCompile(`<h6[^>]*>(.*?)</h6>`)
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
	result = newlineNormalizeRegex.ReplaceAllString(result, "\n\n")
	result = strings.TrimSpace(result)

	return result
}

// convertHTMLToMarkdown は基本的なHTMLタグをMarkdown記法に変換する
func convertHTMLToMarkdown(text string) string {
	result := text

	// はてなブログのASIN詳細タグを変換（他の変換の前に実行）
	result = asinDetailWithPTagRegex.ReplaceAllString(result, "[asin:$1:detail]")
	result = asinDetailNestedRegex.ReplaceAllString(result, "[asin:$1:detail]")
	result = asinDetailSimpleRegex.ReplaceAllString(result, "[asin:$1:detail]")

	// <br> タグを改行に変換
	result = brTagRegex.ReplaceAllString(result, "\n")

	// <p> タグを段落に変換
	result = pTagOpenRegex.ReplaceAllString(result, "")
	result = pTagCloseRegex.ReplaceAllString(result, "\n\n")

	// <strong> や <b> タグを太字に変換
	result = strongTagRegex.ReplaceAllString(result, "**$1**")

	// <em> や <i> タグを斜体に変換
	result = emTagRegex.ReplaceAllString(result, "*$1*")

	// <code> タグを処理（HTMLエスケープを復元してバッククォートで囲む）
	result = codeTagRegex.ReplaceAllStringFunc(result, func(match string) string {
		content := codeTagRegex.ReplaceAllString(match, "$1")
		// HTMLエスケープを復元
		unescaped := html.UnescapeString(content)
		return "`" + unescaped + "`"
	})

	// <a> タグをMarkdownリンクに変換
	result = aTagRegex.ReplaceAllString(result, "[$2]($1)")

	// <img> タグをMarkdown画像に変換
	result = imgWithAltRegex.ReplaceAllString(result, "![$2]($1)")
	result = imgAltSrcRegex.ReplaceAllString(result, "![$1]($2)")
	result = imgSimpleRegex.ReplaceAllString(result, "![]($1)")

	// <h1> から <h6> タグをMarkdownヘッダーに変換
	result = h1Regex.ReplaceAllString(result, "# $1")
	result = h2Regex.ReplaceAllString(result, "## $1")
	result = h3Regex.ReplaceAllString(result, "### $1")
	result = h4Regex.ReplaceAllString(result, "#### $1")
	result = h5Regex.ReplaceAllString(result, "##### $1")
	result = h6Regex.ReplaceAllString(result, "###### $1")

	// <blockquote> タグを引用に変換
	result = blockquoteRegex.ReplaceAllStringFunc(result, func(match string) string {
		content := blockquoteInnerRegex.ReplaceAllString(match, "$1")
		lines := strings.Split(strings.TrimSpace(content), "\n")
		var quotedLines []string
		for _, line := range lines {
			quotedLines = append(quotedLines, "> "+strings.TrimSpace(line))
		}
		return strings.Join(quotedLines, "\n")
	})

	// <ul> と <li> タグをMarkdownリストに変換
	result = ulOpenRegex.ReplaceAllString(result, "")
	result = ulCloseRegex.ReplaceAllString(result, "\n")
	result = liTagRegex.ReplaceAllString(result, "- $1")

	// <ol> タグを番号付きリストに変換（簡易版）
	result = olOpenRegex.ReplaceAllString(result, "")
	result = olCloseRegex.ReplaceAllString(result, "\n")

	return result
}
