package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Entry MovableType形式のエントリーを表現する構造体
type Entry struct {
	Author     string
	Title      string
	Basename   string
	Status     string
	Date       string
	Category   string
	Body       string
	ImageURL   string
}

// parseEntries MTファイルを解析してエントリー一覧を返す
func parseEntries(filename string) ([]Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	var currentEntry Entry
	var inBody bool
	var bodyLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// エントリー区切り
		if line == "--------" {
			if currentEntry.Title != "" {
				currentEntry.Body = strings.Join(bodyLines, "\n")
				entries = append(entries, currentEntry)
			}
			currentEntry = Entry{}
			inBody = false
			bodyLines = []string{}
			continue
		}

		// ボディ開始
		if line == "BODY:" {
			inBody = true
			continue
		}

		// ヘッダー終了
		if line == "-----" && !inBody {
			continue
		}

		// ボディ部分
		if inBody {
			if line == "-----" {
				inBody = false
				continue
			}
			bodyLines = append(bodyLines, line)
			continue
		}

		// メタデータの解析
		if strings.HasPrefix(line, "AUTHOR: ") {
			currentEntry.Author = strings.TrimPrefix(line, "AUTHOR: ")
		} else if strings.HasPrefix(line, "TITLE: ") {
			currentEntry.Title = strings.TrimPrefix(line, "TITLE: ")
		} else if strings.HasPrefix(line, "BASENAME: ") {
			currentEntry.Basename = strings.TrimPrefix(line, "BASENAME: ")
		} else if strings.HasPrefix(line, "STATUS: ") {
			currentEntry.Status = strings.TrimPrefix(line, "STATUS: ")
		} else if strings.HasPrefix(line, "DATE: ") {
			currentEntry.Date = strings.TrimPrefix(line, "DATE: ")
		} else if strings.HasPrefix(line, "CATEGORY: ") {
			currentEntry.Category = strings.TrimPrefix(line, "CATEGORY: ")
		} else if strings.HasPrefix(line, "IMAGE: ") {
			currentEntry.ImageURL = strings.TrimPrefix(line, "IMAGE: ")
		}
	}

	// 最後のエントリーを追加
	if currentEntry.Title != "" {
		currentEntry.Body = strings.Join(bodyLines, "\n")
		entries = append(entries, currentEntry)
	}

	return entries, scanner.Err()
}

// generateMarkdown エントリーをHatena Blog形式のMarkdownに変換
func generateMarkdown(entry Entry) string {
	var md strings.Builder

	// タイトル
	md.WriteString("# ")
	md.WriteString(entry.Title)
	md.WriteString("\n\n")

	// メタデータ
	if entry.Category != "" {
		md.WriteString("カテゴリ: ")
		md.WriteString(entry.Category)
		md.WriteString("\n\n")
	}

	// 画像がある場合
	if entry.ImageURL != "" {
		md.WriteString("![")
		md.WriteString(entry.Title)
		md.WriteString("](")
		md.WriteString(entry.ImageURL)
		md.WriteString(")\n\n")
	}

	// ボディ（HTMLをそのまま保持、はてなブログのMarkdownはHTMLも混在可能）
	md.WriteString(entry.Body)

	return md.String()
}

// generateFilename ファイル名を生成
func generateFilename(entry Entry) string {
	// 危険な文字を除去してファイル名を作成
	title := entry.Title
	unsafe := regexp.MustCompile(`[<>:"/\\|?*]`)
	title = unsafe.ReplaceAllString(title, "_")
	
	// バックスラッシュからアンダースコアに変換後、空白をアンダースコアに
	title = strings.ReplaceAll(title, " ", "_")
	
	// 日付プレフィックスを生成
	var datePrefix string
	if entry.Basename != "" && regexp.MustCompile(`^\d{4}/\d{2}/\d{2}/`).MatchString(entry.Basename) {
		// Basenameが日付形式の場合
		basename := strings.ReplaceAll(entry.Basename, "/", "-")
		datePrefix = basename
	} else {
		// DATEフィールドから日付を抽出 (MM/DD/YYYY HH:MM:SS → YYYY-MM-DD)
		if entry.Date != "" {
			dateRegex := regexp.MustCompile(`(\d{2})/(\d{2})/(\d{4}) (\d{2}):(\d{2}):(\d{2})`)
			matches := dateRegex.FindStringSubmatch(entry.Date)
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

func main() {
	filename := "blog.basyura.org.export.txt"
	
	// ファイルの存在確認
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("エラー: %s が見つかりません\n", filename)
		os.Exit(1)
	}

	// エントリーの解析
	entries, err := parseEntries(filename)
	if err != nil {
		fmt.Printf("ファイル解析エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("解析完了: %d個のエントリーが見つかりました\n", len(entries))

	// 出力ディレクトリの作成
	outputDir := "entries"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("出力ディレクトリ作成エラー: %v\n", err)
		os.Exit(1)
	}

	// テスト用に最初の10件のみ処理
	testEntries := entries
	if len(entries) > 10 {
		testEntries = entries[:10]
		fmt.Printf("テストモード: 最初の10件のみ処理します\n")
	}

	// 各エントリーをMarkdownファイルとして出力
	for i, entry := range testEntries {
		filename := generateFilename(entry)
		filepath := filepath.Join(outputDir, filename)
		
		content := generateMarkdown(entry)
		
		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			fmt.Printf("ファイル書き込みエラー (%s): %v\n", filename, err)
			continue
		}
		
		fmt.Printf("%d: %s を作成しました\n", i+1, filename)
	}
	
	fmt.Println("変換完了！")
}
