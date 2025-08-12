package entry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEntries(t *testing.T) {
	// テスト用の一時ファイルを作成
	testContent := `AUTHOR: test_author
TITLE: Test Entry 1
BASENAME: test_entry_1
STATUS: Publish
DATE: 01/01/2023 12:00:00 AM
CATEGORY: テストカテゴリ
-----
BODY:
これはテスト用のエントリー本文です。
複数行のテストです。
-----
--------
AUTHOR: another_author
TITLE: Test Entry 2
BASENAME: test_entry_2
STATUS: Draft
DATE: 01/02/2023 12:00:00 AM
CATEGORY: 別のカテゴリ
IMAGE: https://example.com/image.jpg
-----
BODY:
2番目のテストエントリーです。
画像URLも含まれています。
-----
--------`

	// 一時ディレクトリとファイルを作成
	tmpDir, err := os.MkdirTemp("", "mttohmd_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test_entries.txt")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// ParseEntriesをテスト
	entries, err := ParseEntries(testFile)
	if err != nil {
		t.Fatalf("ParseEntries failed: %v", err)
	}

	// エントリー数の確認
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// 1番目のエントリーの確認
	if entries[0].Author != "test_author" {
		t.Errorf("Expected Author 'test_author', got '%s'", entries[0].Author)
	}
	if entries[0].Title != "Test Entry 1" {
		t.Errorf("Expected Title 'Test Entry 1', got '%s'", entries[0].Title)
	}
	if entries[0].Basename != "test_entry_1" {
		t.Errorf("Expected Basename 'test_entry_1', got '%s'", entries[0].Basename)
	}
	if entries[0].Status != "Publish" {
		t.Errorf("Expected Status 'Publish', got '%s'", entries[0].Status)
	}
	if entries[0].Date != "01/01/2023 12:00:00 AM" {
		t.Errorf("Expected Date '01/01/2023 12:00:00 AM', got '%s'", entries[0].Date)
	}
	if entries[0].Category != "テストカテゴリ" {
		t.Errorf("Expected Category 'テストカテゴリ', got '%s'", entries[0].Category)
	}
	if entries[0].Body != "これはテスト用のエントリー本文です。\n複数行のテストです。" {
		t.Errorf("Expected Body content, got '%s'", entries[0].Body)
	}
	if entries[0].ImageURL != "" {
		t.Errorf("Expected empty ImageURL, got '%s'", entries[0].ImageURL)
	}

	// 2番目のエントリーの確認
	if entries[1].Author != "another_author" {
		t.Errorf("Expected Author 'another_author', got '%s'", entries[1].Author)
	}
	if entries[1].Title != "Test Entry 2" {
		t.Errorf("Expected Title 'Test Entry 2', got '%s'", entries[1].Title)
	}
	if entries[1].Basename != "test_entry_2" {
		t.Errorf("Expected Basename 'test_entry_2', got '%s'", entries[1].Basename)
	}
	if entries[1].Status != "Draft" {
		t.Errorf("Expected Status 'Draft', got '%s'", entries[1].Status)
	}
	if entries[1].Date != "01/02/2023 12:00:00 AM" {
		t.Errorf("Expected Date '01/02/2023 12:00:00 AM', got '%s'", entries[1].Date)
	}
	if entries[1].Category != "別のカテゴリ" {
		t.Errorf("Expected Category '別のカテゴリ', got '%s'", entries[1].Category)
	}
	if entries[1].Body != "2番目のテストエントリーです。\n画像URLも含まれています。" {
		t.Errorf("Expected Body content, got '%s'", entries[1].Body)
	}
	if entries[1].ImageURL != "https://example.com/image.jpg" {
		t.Errorf("Expected ImageURL 'https://example.com/image.jpg', got '%s'", entries[1].ImageURL)
	}
}

func TestParseEntriesNonExistentFile(t *testing.T) {
	_, err := ParseEntries("non_existent_file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestParseEntriesEmptyFile(t *testing.T) {
	// 空ファイルのテスト
	tmpDir, err := os.MkdirTemp("", "mttohmd_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(testFile, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := ParseEntries(testFile)
	if err != nil {
		t.Fatalf("ParseEntries failed: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty file, got %d", len(entries))
	}
}

func TestParseEntriesSingleEntry(t *testing.T) {
	// 単一エントリーのテスト
	testContent := `AUTHOR: single_author
TITLE: Single Entry
BASENAME: single_entry
STATUS: Publish
DATE: 01/01/2023 12:00:00 AM
CATEGORY: シングル
-----
BODY:
単一エントリーのテストです。
-----`

	tmpDir, err := os.MkdirTemp("", "mttohmd_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "single_entry.txt")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	entries, err := ParseEntries(testFile)
	if err != nil {
		t.Fatalf("ParseEntries failed: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Title != "Single Entry" {
		t.Errorf("Expected Title 'Single Entry', got '%s'", entries[0].Title)
	}
	if entries[0].Body != "単一エントリーのテストです。" {
		t.Errorf("Expected Body content, got '%s'", entries[0].Body)
	}
}
