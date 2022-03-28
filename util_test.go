package lrc2srt

import "testing"

func TestMillisecond2Time(t *testing.T) {
	type args struct {
		millisecond int
	}
	tests := []struct {
		name   string
		args   args
		wantH  int
		wantM  int
		wantS  int
		wantMs int
	}{
		{"ms_0", args{0}, 0, 0, 0, 0},
		{"ms_100", args{100}, 0, 0, 0, 100},
		{"ms_1000", args{1000}, 0, 0, 1, 0},
		{"ms_1100", args{1100}, 0, 0, 1, 100},
		{"ms_60000", args{60000}, 0, 1, 0, 0},
		{"ms_3600000", args{3600000}, 1, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotH, gotM, gotS, gotMs := millisecond2Time(tt.args.millisecond)
			if gotH != tt.wantH {
				t.Errorf("Millisecond2Time() gotH = %v, want %v", gotH, tt.wantH)
			}
			if gotM != tt.wantM {
				t.Errorf("Millisecond2Time() gotM = %v, want %v", gotM, tt.wantM)
			}
			if gotS != tt.wantS {
				t.Errorf("Millisecond2Time() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotMs != tt.wantMs {
				t.Errorf("Millisecond2Time() gotMs = %v, want %v", gotMs, tt.wantMs)
			}
		})
	}
}

func TestTime2Millisecond(t *testing.T) {
	type args struct {
		m  int
		s  int
		ms int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"0:0.0", args{0, 0, 0}, 0},
		{"0:0.1", args{0, 0, 1}, 1},
		{"0:0.999", args{0, 0, 999}, 999},
		{"0:1.0", args{0, 1, 0}, 1000},
		{"0:1.999", args{0, 1, 999}, 1999},
		{"1:0.0", args{1, 0, 0}, 60000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := time2Millisecond(tt.args.m, tt.args.s, tt.args.ms); got != tt.want {
				t.Errorf("Time2Millisecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTime2Millisecond(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time2Millisecond(999, 999, 999)
	}
}

func BenchmarkMillisecond2Time(b *testing.B) {
	for i := 0; i < b.N; i++ {
		millisecond2Time(9999999999)
	}
}
