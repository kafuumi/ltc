package lrc2srt

import "testing"

func TestSRTContent_String(t *testing.T) {
	type fields struct {
		Index int
		Start int
		End   int
		Text  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"srtContent String()", fields{1, 10, 20, "test"},
			"1\n00:00:00,010 --> 00:00:00,020\ntest\n\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SRTContent{
				Index: tt.fields.Index,
				Start: tt.fields.Start,
				End:   tt.fields.End,
				Text:  tt.fields.Text,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLrcToSrt(t *testing.T) {
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
	srt := LrcToSrt(l)
	if srt.Title != "title" {
		t.Errorf("SRT Title=%s, want=%s", srt.Title, "title")
	}
	if srt.Artist != "artist" {
		t.Errorf("SRT Artist=%s, want=%s", srt.Artist, "altist")
	}
	index := 0
	for it := srt.Content.Iterator(); it.Has(); {
		c := it.Next().Text
		if c != content[index] {
			t.Errorf("srt Text=%s, want=%s", c, content[index])
		}
		index++
	}
}

func TestSRT_MergeStack(t *testing.T) {

}

func TestSRT_MergeUp(t *testing.T) {

}
func TestSRT_MergeBottom(t *testing.T) {

}
