package lrc2srt

import (
	"fmt"
	"testing"
)

func TestGet163Lyric(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"1423123512"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("id=%s", tt.input), func(t *testing.T) {
			l, lt := Get163Lyric(tt.input)
			if l == "" || lt == "" {
				t.Errorf("get cloud lyric faild, id = %s", tt.input)
			}
		})
	}
}
