package entry

import (
	"bufio"
	"os"
	"strings"
)

// Entry MovableType形式のエントリーを表現する構造体
type Entry struct {
	Author   string
	Title    string
	Basename string
	Status   string
	Date     string
	Category string
	Body     string
	ImageURL string
}

// ParseEntries MTファイルを解析してエントリー一覧を返す
func ParseEntries(filename string) ([]Entry, error) {
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
