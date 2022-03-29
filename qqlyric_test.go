package lrc2srt

import (
	"fmt"
	"testing"
)

func TestGetQQLyric(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"0002Jztl3eJKu0"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("id=%s", tt.input), func(t *testing.T) {
			l, lt := GetQQLyric(tt.input)
			if l == "" || lt == "" {
				t.Errorf("get cloud lyric faild, id = %s", tt.input)
			}
		})
	}
}
