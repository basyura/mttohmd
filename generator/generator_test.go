package generator

import (
	"strings"
	"testing"

	"mttohmd/entry"
)

func TestGenerateFilename(t *testing.T) {
	tests := []struct {
		name     string
		entry    entry.Entry
		expected string
	}{
		{
			name: "基本的なタイトル",
			entry: entry.Entry{
				Title: "Test Blog Post",
			},
			expected: "Test_Blog_Post.md",
		},
		{
			name: "危険な文字を含むタイトル",
			entry: entry.Entry{
				Title: "Test<>:\"/\\|?*Post",
			},
			expected: "Test_________Post.md",
		},
		{
			name: "スペースを含むタイトル",
			entry: entry.Entry{
				Title: "My First Blog Post",
			},
			expected: "My_First_Blog_Post.md",
		},
		{
			name: "日付形式のBasenameあり",
			entry: entry.Entry{
				Title:    "Blog Post",
				Basename: "2023/01/15/blog-post",
			},
			expected: "2023-01-15-blog-post_Blog_Post.md",
		},
		{
			name: "DATE形式から日付抽出",
			entry: entry.Entry{
				Title: "Another Post",
				Date:  "01/15/2023 14:30:45 PM",
			},
			expected: "2023-01-15-143045_Another_Post.md",
		},
		{
			name: "BasenameとDATEの両方がある場合（Basenameが優先）",
			entry: entry.Entry{
				Title:    "Priority Test",
				Basename: "2023/01/15/priority",
				Date:     "01/16/2023 10:00:00 AM",
			},
			expected: "2023-01-15-priority_Priority_Test.md",
		},
		{
			name: "日付形式ではないBasename",
			entry: entry.Entry{
				Title:    "No Date Format",
				Basename: "simple-basename",
				Date:     "01/15/2023 12:00:00 PM",
			},
			expected: "2023-01-15-120000_No_Date_Format.md",
		},
		{
			name: "日付情報なし",
			entry: entry.Entry{
				Title: "Simple Post",
			},
			expected: "Simple_Post.md",
		},
		{
			name: "日本語タイトル",
			entry: entry.Entry{
				Title: "日本語のタイトル テスト",
			},
			expected: "日本語のタイトル_テスト.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateFilename(tt.entry)
			if result != tt.expected {
				t.Errorf("GenerateFilename() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGenerateMTContent(t *testing.T) {
	// 全フィールドを含むテスト
	fullEntry := entry.Entry{
		Author:   "test_author",
		Title:    "Test Entry",
		Basename: "test_entry",
		Status:   "Publish",
		Date:     "01/15/2023 12:00:00 AM",
		Category: "Technology",
		Body:     "これはテスト用の本文です。\n複数行のテストです。",
		ImageURL: "https://example.com/image.jpg",
	}

	result := GenerateMTContent(fullEntry)

	// 各フィールドの存在を確認
	expectedFields := []string{
		"AUTHOR: test_author",
		"TITLE: Test Entry",
		"BASENAME: test_entry",
		"STATUS: Publish",
		"DATE: 01/15/2023 12:00:00 AM",
		"CATEGORY: Technology",
		"IMAGE: https://example.com/image.jpg",
		"-----",
		"BODY:",
		"これはテスト用の本文です。",
		"複数行のテストです。",
	}

	for _, field := range expectedFields {
		if !strings.Contains(result, field) {
			t.Errorf("Expected field %q not found in result", field)
		}
	}

	// 構造の確認
	if !strings.HasSuffix(result, "-----\n") {
		t.Error("Result should end with -----")
	}
}

func TestGenerateMTContentMinimal(t *testing.T) {
	// 最小限のフィールドのみのテスト
	minimalEntry := entry.Entry{
		Title: "Minimal Entry",
		Body:  "Simple body",
	}

	result := GenerateMTContent(minimalEntry)

	// 必須フィールドの確認
	expectedFields := []string{
		"TITLE: Minimal Entry",
		"-----",
		"BODY:",
		"Simple body",
	}

	for _, field := range expectedFields {
		if !strings.Contains(result, field) {
			t.Errorf("Expected field %q not found in result", field)
		}
	}

	// 空のフィールドが出力されていないことを確認
	unexpectedFields := []string{
		"AUTHOR:",
		"BASENAME:",
		"STATUS:",
		"DATE:",
		"CATEGORY:",
		"IMAGE:",
	}

	for _, field := range unexpectedFields {
		if strings.Contains(result, field) {
			t.Errorf("Unexpected field %q found in result", field)
		}
	}
}

func TestGenerateMTContentEmpty(t *testing.T) {
	// 空のエントリーのテスト
	emptyEntry := entry.Entry{}

	result := GenerateMTContent(emptyEntry)

	// TITLEは空でも出力される
	if !strings.Contains(result, "TITLE: ") {
		t.Error("TITLE field should be present even when empty")
	}

	// BODYセクションの確認
	if !strings.Contains(result, "BODY:") {
		t.Error("BODY section should be present")
	}

	// 構造の確認
	if !strings.Contains(result, "-----") {
		t.Error("Delimiter should be present")
	}
}

func TestGenerateMTContentSpecialCharacters(t *testing.T) {
	// 特殊文字を含むテスト
	specialEntry := entry.Entry{
		Title:    "Special: Characters & Symbols <test>",
		Body:     "Body with\nnewlines and\ttabs",
		Category: "Test, Special Characters",
	}

	result := GenerateMTContent(specialEntry)

	// 特殊文字がそのまま保持されることを確認
	if !strings.Contains(result, "TITLE: Special: Characters & Symbols <test>") {
		t.Error("Special characters in title should be preserved")
	}
	if !strings.Contains(result, "CATEGORY: Test, Special Characters") {
		t.Error("Special characters in category should be preserved")
	}
	if !strings.Contains(result, "Body with\nnewlines and\ttabs") {
		t.Error("Newlines and tabs in body should be preserved")
	}
}
