package ltc

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hami-Lemon/ltc/glist"
)

type ASSNode struct {
	Start    int    //开始时间
	End      int    //结束时间
	Dialogue string //内容
}

func (a *ASSNode) String() string {
	builder := strings.Builder{}
	sh, sm, ss, sms := millisecond2Time(a.Start)
	eh, em, es, ems := millisecond2Time(a.End)
	builder.WriteString("Dialogue: 0,")
	sms /= 10
	ems /= 10
	builder.WriteString(fmt.Sprintf("%d:%02d:%02d.%02d,%d:%02d:%02d.%02d,Default,,",
		sh, sm, ss, sms, eh, em, es, ems))
	builder.WriteString("0000,0000,0000,,")
	builder.WriteString(a.Dialogue)
	return builder.String()
}

type ASS struct {
	Content glist.Queue[*ASSNode]
}

func LrcToAss(lrc *LRC) *ASS {
	return SrtToAss(LrcToSrt(lrc))
}

func SrtToAss(srt *SRT) *ASS {
	if srt == nil {
		return nil
	}
	ass := &ASS{
		Content: glist.NewLinkedList[*ASSNode](),
	}
	for it := srt.Content.Iterator(); it.Has(); {
		s := it.Next()
		node := &ASSNode{
			Start:    s.Start,
			End:      s.End,
			Dialogue: s.Text,
		}
		ass.Content.PushBack(node)
	}
	return ass
}

func (a *ASS) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		//不存在对应文件夹
		if os.IsNotExist(err) {
			panic("文件夹不存在:" + filepath.Dir(path))
		}
		return err
	}
	err = a.Write(f)
	err = f.Close()
	return err
}

func (a *ASS) Write(dst io.Writer) error {
	if err := writeScriptInfo(dst); err != nil {
		return err
	}
	if err := writeStyles(dst); err != nil {
		return err
	}
	if err := writeEventHeader(dst); err != nil {
		return err
	}
	for it := a.Content.Iterator(); it.Has(); {
		temp := it.Next()
		r := temp.String()
		_, err := fmt.Fprintf(dst, "%s\n", r)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeScriptInfo(dst io.Writer) error {
	text := `[Script Info]
Title: LRC ASS file
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080
Collisions: Reverse
WrapStyle: 2

`
	_, err := dst.Write([]byte(text))
	return err
}

func writeStyles(dst io.Writer) error {
	text := `[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,黑体,36,&H00FFFFFF,&H00FFFFFF,&H00000000,&00FFFFFF,-1,0,0,0,100,100,0,0,1,0,1,2,0,0,0,1

`
	_, err := dst.Write([]byte(text))
	return err
}

func writeEventHeader(dst io.Writer) error {
	text := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`
	_, err := dst.Write([]byte(text))
	return err
}
