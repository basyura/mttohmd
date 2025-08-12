package generator

import (
	"fmt"
	"regexp"
	"strings"

	"mttohmd/entry"
)

// GenerateFilename ファイル名を生成
func GenerateFilename(e entry.Entry) string {
	// 危険な文字を除去してファイル名を作成
	title := e.Title
	unsafe := regexp.MustCompile(`[<>:"/\\|?*]`)
	title = unsafe.ReplaceAllString(title, "_")

	// バックスラッシュからアンダースコアに変換後、空白をアンダースコアに
	title = strings.ReplaceAll(title, " ", "_")

	// 日付プレフィックスを生成
	var datePrefix string
	if e.Basename != "" && regexp.MustCompile(`^\d{4}/\d{2}/\d{2}/`).MatchString(e.Basename) {
		// Basenameが日付形式の場合
		basename := strings.ReplaceAll(e.Basename, "/", "-")
		datePrefix = basename
	} else {
		// DATEフィールドから日付を抽出 (MM/DD/YYYY HH:MM:SS → YYYY-MM-DD)
		if e.Date != "" {
			dateRegex := regexp.MustCompile(`(\d{2})/(\d{2})/(\d{4}) (\d{2}):(\d{2}):(\d{2})`)
			matches := dateRegex.FindStringSubmatch(e.Date)
			if len(matches) > 6 {
				month, day, year, hour, minute, second := matches[1], matches[2], matches[3], matches[4], matches[5], matches[6]
				datePrefix = fmt.Sprintf("%s-%s-%s-%s%s%s", year, month, day, hour, minute, second)
			}
		}
	}

	if datePrefix != "" {
		return fmt.Sprintf("%s_%s.md", datePrefix, title)
	}

	return fmt.Sprintf("%s.md", title)
}

// GenerateMTContent エントリーをMovableType形式のまま出力
func GenerateMTContent(e entry.Entry) string {
	var mt strings.Builder
	
	// MovableType形式のヘッダー
	if e.Author != "" {
		mt.WriteString("AUTHOR: ")
		mt.WriteString(e.Author)
		mt.WriteString("\n")
	}
	
	mt.WriteString("TITLE: ")
	mt.WriteString(e.Title)
	mt.WriteString("\n")
	
	if e.Basename != "" {
		mt.WriteString("BASENAME: ")
		mt.WriteString(e.Basename)
		mt.WriteString("\n")
	}
	
	if e.Status != "" {
		mt.WriteString("STATUS: ")
		mt.WriteString(e.Status)
		mt.WriteString("\n")
	}
	
	if e.Date != "" {
		mt.WriteString("DATE: ")
		mt.WriteString(e.Date)
		mt.WriteString("\n")
	}
	
	if e.Category != "" {
		mt.WriteString("CATEGORY: ")
		mt.WriteString(e.Category)
		mt.WriteString("\n")
	}
	
	if e.ImageURL != "" {
		mt.WriteString("IMAGE: ")
		mt.WriteString(e.ImageURL)
		mt.WriteString("\n")
	}
	
	mt.WriteString("-----\n")
	mt.WriteString("BODY:\n")
	mt.WriteString(e.Body)
	mt.WriteString("\n-----\n")
	
	return mt.String()
}