package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mttohmd/converter"
	"mttohmd/entry"
	"mttohmd/generator"
)

func main() {
	filename := "blog.basyura.org.export.txt"

	// ファイルの存在確認
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("エラー: %s が見つかりません\n", filename)
		os.Exit(1)
	}

	// エントリーの解析
	entries, err := entry.ParseEntries(filename)
	if err != nil {
		fmt.Printf("ファイル解析エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("解析完了: %d個のエントリーが見つかりました\n", len(entries))

	// 出力ディレクトリの作成
	mtsDir := "mts"
	mdsDir := "mds"

	if err := os.MkdirAll(mtsDir, 0755); err != nil {
		fmt.Printf("MTSディレクトリ作成エラー: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(mdsDir, 0755); err != nil {
		fmt.Printf("MDSディレクトリ作成エラー: %v\n", err)
		os.Exit(1)
	}

	// テスト用に最初の10件のみ処理
	testEntries := entries
	if len(entries) > 10 {
		testEntries = entries[:10]
		fmt.Printf("テストモード: 最初の10件のみ処理します\n")
	}

	// 各エントリーを2つの形式で出力
	for i, e := range testEntries {
		filename := generator.GenerateFilename(e)

		// MT形式でmtsフォルダに出力
		mtFilename := strings.Replace(filename, ".md", ".txt", 1)
		mtFilepath := filepath.Join(mtsDir, mtFilename)
		mtContent := generator.GenerateMTContent(e)

		if err := os.WriteFile(mtFilepath, []byte(mtContent), 0644); err != nil {
			fmt.Printf("MTファイル書き込みエラー (%s): %v\n", mtFilename, err)
		} else {
			fmt.Printf("%d: MTS/%s を作成しました\n", i+1, mtFilename)
		}

		// Markdown形式でmdsフォルダに出力
		mdFilepath := filepath.Join(mdsDir, filename)
		mdContent := converter.ToMarkdown(e)

		if err := os.WriteFile(mdFilepath, []byte(mdContent), 0644); err != nil {
			fmt.Printf("Markdownファイル書き込みエラー (%s): %v\n", filename, err)
		} else {
			fmt.Printf("%d: MDS/%s を作成しました\n", i+1, filename)
		}
	}

	fmt.Println("変換完了！")
}
