package converter

import (
	"strings"
	"testing"

	"mttohmd/entry"
)

func TestToMarkdown(t *testing.T) {
	// テスト用のエントリーを作成
	testEntry := entry.Entry{
		Author:   "test_author",
		Title:    "Test Blog Post",
		Basename: "test_blog_post",
		Status:   "Publish",
		Date:     "01/15/2023 12:00:00 AM",
		Category: "Technology, Go, Testing",
		Body:     "これはテスト用の本文です。\n<strong>太字</strong>のテストも含まれます。",
		ImageURL: "https://example.com/test.jpg",
	}

	result := ToMarkdown(testEntry)

	// フロントマターの確認
	if !strings.Contains(result, "---") {
		t.Error("Expected frontmatter delimiters not found")
	}
	if !strings.Contains(result, "Title: Test Blog Post") {
		t.Error("Expected title not found in result")
	}
	if !strings.Contains(result, "Date: 01/15/2023 12:00:00 AM") {
		t.Error("Expected date not found in result")
	}

	// カテゴリーの確認
	if !strings.Contains(result, "Category:") {
		t.Error("Expected category section not found")
	}
	if !strings.Contains(result, "- Technology") {
		t.Error("Expected Technology category not found")
	}
	if !strings.Contains(result, "- Go") {
		t.Error("Expected Go category not found")
	}
	if !strings.Contains(result, "- Testing") {
		t.Error("Expected Testing category not found")
	}

	// 本文の確認
	if !strings.Contains(result, "これはテスト用の本文です。") {
		t.Error("Expected body content not found")
	}
	if !strings.Contains(result, "**太字**") {
		t.Error("Expected converted bold text not found")
	}

	// 画像の確認
	if !strings.Contains(result, "![Test Blog Post](https://example.com/test.jpg)") {
		t.Error("Expected image markdown not found")
	}
}

func TestToMarkdownWithoutCategory(t *testing.T) {
	testEntry := entry.Entry{
		Title: "No Category Post",
		Body:  "Simple post without category",
	}

	result := ToMarkdown(testEntry)

	if strings.Contains(result, "Category:") {
		t.Error("Category section should not be present")
	}
}

func TestToMarkdownWithoutImage(t *testing.T) {
	testEntry := entry.Entry{
		Title: "No Image Post",
		Body:  "Simple post without image",
	}

	result := ToMarkdown(testEntry)

	if strings.Contains(result, "![") {
		t.Error("Image markdown should not be present")
	}
}

func TestConvertMTToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基本的なテキスト",
			input:    "普通のテキストです。",
			expected: "普通のテキストです。",
		},
		{
			name:     "改行の正規化",
			input:    "テスト\r\nテスト\rテスト",
			expected: "テスト\nテスト\nテスト",
		},
		{
			name:     "空行の整理",
			input:    "テスト\n\n\n\nテスト",
			expected: "テスト\n\nテスト",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertMTToMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("convertMTToMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConvertHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "br タグ",
			input:    "行1<br>行2<BR/>行3",
			expected: "行1\n行2\n行3",
		},
		{
			name:     "p タグ",
			input:    "<p>段落1</p><p>段落2</p>",
			expected: "段落1\n\n段落2\n\n",
		},
		{
			name:     "strong と b タグ",
			input:    "<strong>強調</strong>と<b>太字</b>",
			expected: "**強調**と**太字**",
		},
		{
			name:     "em と i タグ",
			input:    "<em>斜体</em>と<i>イタリック</i>",
			expected: "*斜体*と*イタリック*",
		},
		{
			name:     "a タグ",
			input:    "<a href=\"https://example.com\">リンク</a>",
			expected: "[リンク](https://example.com)",
		},
		{
			name:     "img タグ",
			input:    "<img src=\"test.jpg\" alt=\"テスト画像\" />",
			expected: "![テスト画像](test.jpg)",
		},
		{
			name:     "h1-h6 タグ",
			input:    "<h1>見出し1</h1><h2>見出し2</h2><h3>見出し3</h3>",
			expected: "# 見出し1## 見出し2### 見出し3",
		},
		{
			name:     "blockquote タグ",
			input:    "<blockquote>これは引用文です</blockquote>",
			expected: "> これは引用文です",
		},
		{
			name:     "ul と li タグ",
			input:    "<ul><li>項目1</li><li>項目2</li></ul>",
			expected: "- 項目1- 項目2\n",
		},
		{
			name:     "複合的なHTML",
			input:    "<p><strong>重要:</strong> <a href=\"#\">リンク</a>です。</p>",
			expected: "**重要:** [リンク](#)です。\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertHTMLToMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("convertHTMLToMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestConvertHTMLToMarkdownComplexBlockquote(t *testing.T) {
	input := "<blockquote>複数行の\n引用文\nテストです</blockquote>"
	expected := "> 複数行の\n> 引用文\n> テストです"

	result := convertHTMLToMarkdown(input)
	if result != expected {
		t.Errorf("convertHTMLToMarkdown() = %q, want %q", result, expected)
	}
}
