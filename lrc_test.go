package lrc2srt

import (
	"reflect"
	"testing"
)

func TestParseLRC(t *testing.T) {
	lrc := `[ar:artist]
[al:album]
[ti:title]
[by:author]
[00:24.83] 天涯的尽头 有谁去过
[00:28.53] 山水优雅着 保持沉默
[00:32.20] 我们的青春却热闹很多
[00:35.38] 而且是谁都 不准偷
`
	content := []string{
		"天涯的尽头 有谁去过", "山水优雅着 保持沉默", "我们的青春却热闹很多", "而且是谁都 不准偷",
	}
	l := ParseLRC(lrc)
	if l.Artist != "artist" {
		t.Errorf("LRC Artist=%s, want=%s", l.Artist, "artist")
	}
	if l.Album != "album" {
		t.Errorf("LRC Album=%s, want=%s", l.Album, "album")
	}
	if l.Title != "title" {
		t.Errorf("LRC Title=%s, want=%s", l.Title, "title")
	}
	if l.Author != "author" {
		t.Errorf("LRC Author=%s, want=%s", l.Author, "author")
	}
	lrcList := l.LrcList
	index := 0
	for it := lrcList.Iterator(); it.Has(); {
		c := it.Next().content
		if c != content[index] {
			t.Errorf("LRCNode Content=%s, want=%s", c, content[index])
		}
		index++

	}
}

func TestSplitLyric(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *LRCNode
	}{
		{"lrc:[00:49.88] 有一些话想 对你说", []string{"00", "49", ".88", " 有一些话想 对你说"},
			&LRCNode{time: time2Millisecond(0, 49, 880), content: "有一些话想 对你说"}},
		{"lrc:[00:49:88] 有一些话想 对你说", []string{"00", "49", ":88", " 有一些话想 对你说"},
			&LRCNode{time: time2Millisecond(0, 49, 880), content: "有一些话想 对你说"}},
		{"lrc:[00:49.880] 有一些话想 对你说", []string{"00", "49", ".880", " 有一些话想 对你说"},
			&LRCNode{time: time2Millisecond(0, 49, 880), content: "有一些话想 对你说"}},
		{"lrc:[00:49] 有一些话想 对你说", []string{"00", "49", " 有一些话想 对你说"},
			&LRCNode{time: time2Millisecond(0, 49, 0), content: "有一些话想 对你说"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitLyric(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitLyric() = %v, want %v", got, tt.want)
			}
		})
	}
}
