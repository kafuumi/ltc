package main

import "testing"

func BenchmarkSRTContent_String(b *testing.B) {
	srt := SRTContent{
		Index: 10,
		Start: 100,
		End:   200,
		Text:  "言语 不起作用,想看到 具体行动",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = srt.String()
	}
}

func BenchmarkLrc2Srt(b *testing.B) {
	id := "28891491"
	lyric, lyricT := Get163Lyric(id)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Lrc2Srt(lyric), Lrc2Srt(lyricT)
	}
}
