package tools

import "testing"

func TestFilterSensitiveWords(t *testing.T) {
	got := FilterSensitiveWords("这条评论包含敏感词和傻逼内容")
	want := "这条评论包含***和**内容"
	if got != want {
		t.Fatalf("FilterSensitiveWords() = %q, want %q", got, want)
	}
}

func TestFilterSensitiveWordsLeavesNormalTextUnchanged(t *testing.T) {
	const content = "这是一条正常评论"
	if got := FilterSensitiveWords(content); got != content {
		t.Fatalf("FilterSensitiveWords() = %q, want %q", got, content)
	}
}
