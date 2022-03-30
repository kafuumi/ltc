package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/Hami-Lemon/ltc"
	"os"
	"strconv"
	"strings"
)

var (
	input    string     //输入，可以是歌词对应的歌曲id，也可以是文件名
	source   string     //歌词来源，默认163,可选163(网易云音乐)，QQ或qq(QQ音乐)，后续支持：kg(酷狗音乐)
	download boolFlag   //是否只下载歌词，当输入是歌曲id且设置该选项时，只下载歌词而不进行处理
	mode     modeFlag   //如果存在译文时的合并模式
	version  boolFlag   //当前程序版本信息，设置该选项时只输出版本信息
	format   formatFlag //字幕格式，可选: ass,srt，默认为ass
	output   string     //保存的文件名
)

//检查路径path是否有效
func checkPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func parseFlag() {
	flag.StringVar(&input, "i", "", "歌词来源，可以是歌词对应的歌曲id，也可以是歌词文件")
	flag.StringVar(&source, "s", "163", "选择从网易云还是QQ音乐上获取歌词，可选值：163(默认)，qq。")
	flag.Var(&download, "d", "设置该选项时，只下载歌词，而无需转换。")
	flag.Var(&mode, "m", "设置歌词原文和译文的合并模式，可选值：1(默认),2,3。")
	flag.Var(&version, "v", "获取当前程序版本信息。")
	flag.Var(&format, "f", "转换成的字幕文件格式，可选值：ass(默认),srt")
	flag.Usage = func() {
		fmt.Printf("LrcToCaptions(ltc) 将LRC歌词文件转换成字幕文件。\n")
		fmt.Printf("ltc version: %s\n\n", VERSION)
		fmt.Printf("用法：ltc [options] OutputFile\n\n")
		fmt.Printf("options:\n\n")
		flag.PrintDefaults()
		fmt.Println("")
	}
	flag.Parse()
	if other := flag.Args(); len(other) != 0 {
		output = other[0]
	}
	outputProcess()
}

func outputProcess() {
	//处理结果文件名
	if output == "" {
		//和输入源同名
		dot := strings.LastIndex(input, ".")
		if dot == -1 {
			output = input
		} else {
			output = input[:dot]
		}
	}
	//后缀名处理
	suffix := func(o, s string) string {
		if !strings.HasSuffix(o, s) {
			return o + s
		}
		return o
	}
	if download.IsSet() {
		output = suffix(output, ".lrc")
	} else {
		switch format.Value() {
		case FORMAT_SRT:
			output = suffix(output, ".srt")
		case FORMAT_ASS:
			output = suffix(output, ".ass")
		}
	}
}

// boolFlag bool值类型的参数
//实现flags包中的boolFlag接口，设置bool值时不要传具体的值
//即: -flag 等价与 -flag=true
type boolFlag bool

func (b *boolFlag) String() string {
	if b == nil {
		return "false"
	}
	return strconv.FormatBool(bool(*b))
}

func (b *boolFlag) Set(value string) error {
	if f, err := strconv.ParseBool(value); err != nil {
		return err
	} else {
		*b = boolFlag(f)
		return nil
	}
}

func (b *boolFlag) IsBoolFlag() bool {
	return true
}

func (b *boolFlag) IsSet() bool {
	return bool(*b)
}

//歌词合并模式的选项
type modeFlag ltc.SRTMergeMode

func (m *modeFlag) String() string {
	if m == nil {
		return "STACK_MODE"
	}
	switch ltc.SRTMergeMode(*m) {
	case ltc.SRT_MERGE_MODE_STACK:
		return "STACK_MODE"
	case ltc.SRT_MERGE_MODE_UP:
		return "UP_MODE"
	case ltc.SRT_MERGE_MODE_BOTTOM:
		return "BOTTOM_MODE"
	default:
		return "STACK_MODE"
	}
}

func (m *modeFlag) Set(value string) error {
	if value == "" {
		*m = modeFlag(ltc.SRT_MERGE_MODE_STACK)
	}
	v := strings.ToLower(value)
	switch v {
	case "1", "stack":
		*m = modeFlag(ltc.SRT_MERGE_MODE_STACK)
	case "2", "up":
		*m = modeFlag(ltc.SRT_MERGE_MODE_UP)
	case "3", "bottom":
		*m = modeFlag(ltc.SRT_MERGE_MODE_BOTTOM)
	default:
		return errors.New("invalid mode value:" + v + " only support 1, 2, 3")
	}
	return nil
}

func (m *modeFlag) Mode() ltc.SRTMergeMode {
	return ltc.SRTMergeMode(*m)
}

// Format 字幕文件的格式
type Format int

const (
	FORMAT_ASS Format = iota
	FORMAT_SRT
)

type formatFlag Format

func (f *formatFlag) String() string {
	if f == nil {
		return ""
	}
	ft := Format(*f)
	switch ft {
	case FORMAT_SRT:
		return "srt"
	case FORMAT_ASS:
		return "ass"
	}
	return ""
}

func (f *formatFlag) Set(value string) error {
	if value == "" {
		*f = formatFlag(FORMAT_ASS)
	}
	v := strings.ToLower(value)
	switch v {
	case "srt":
		*f = formatFlag(FORMAT_SRT)
	case "ass":
		*f = formatFlag(FORMAT_ASS)
	default:
		return errors.New("invalid format value:" + value)
	}
	return nil
}

func (f *formatFlag) Value() Format {
	return Format(*f)
}
