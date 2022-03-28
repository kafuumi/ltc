package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	l2s "github.com/Hami-Lemon/lrc2srt"
	"github.com/jessevdk/go-flags"
)

type Option struct {
	Id       string `short:"i" long:"id" description:"歌曲的id，网易云和QQ音乐均可。"`
	Input    string `short:"I" long:"input" description:"需要转换的LRC文件路径。"`
	Source   string `short:"s" long:"source" description:"当设置id时有效，指定从网易云（163）还是QQ音乐（qq）上获取歌词。" default:"163" choice:"163" choice:"qq" choice:"QQ"`
	Download bool   `short:"d" long:"download" description:"只下载歌词，而不进行解析。"`
	Mode     int    `short:"m" long:"mode" default:"1" description:"原文和译文的排列模式,可选值有：[1] [2] [3]" choice:"1" choice:"2" choice:"3"`
	Version  bool   `short:"v" long:"version" description:"获取版本信息"`
	Output   string `no-flag:""`
}

const (
	// VERSION 当前版本
	VERSION = `"0.2.3" (build 2022.03.28)`
)

var (
	opt Option
)

func main() {
	//TODO 支持转ass文件
	args, err := flags.Parse(&opt)
	if err != nil {
		os.Exit(0)
	}
	//显示版本信息
	if opt.Version {
		fmt.Printf("LrcToSrt(lts) version %s\n", VERSION)
		return
	}
	//获取保存的文件名
	if len(args) != 0 {
		opt.Output = args[0]
	}

	//获取歌词，lyric为原文歌词，tranLyric为译文歌词
	var lyric, tranLyric string
	if opt.Id != "" {
		if opt.Source != "163" {
			lyric, tranLyric = l2s.GetQQLyric(opt.Id)
		} else {
			lyric, tranLyric = l2s.Get163Lyric(opt.Id)
		}
		//下载歌词
		if opt.Download {
			//对文件名进行处理
			o := opt.Output
			if o == "" {
				o = opt.Id + ".lrc"
			} else if !strings.HasSuffix(o, ".lrc") {
				o += ".lrc"
			}
			l2s.WriteFile(o, lyric)
			if tranLyric != "" {
				l2s.WriteFile("tran_"+o, tranLyric)
			}
			fmt.Println("下载歌词完成！")
			return
		}
	} else if opt.Input != "" {
		//从文件中获取歌词
		if !strings.HasSuffix(opt.Input, ".lrc") {
			fmt.Println("Error: 不支持的格式，目前只支持lrc歌词文件。")
			os.Exit(1)
		}
		lyric = l2s.ReadFile(opt.Input)
		if lyric == "" {
			fmt.Println("获取歌词失败，文件内容为空。")
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: 请指定需要转换的歌词。")
		os.Exit(1)
	}
	lrc, lrcT := l2s.ParseLRC(lyric), l2s.ParseLRC(tranLyric)
	srt, srtT := l2s.LrcToSrt(lrc), l2s.LrcToSrt(lrcT)
	if srtT != nil {
		var mode l2s.SRTMergeMode
		switch opt.Mode {
		case 1:
			mode = l2s.SRT_MERGE_MODE_STACK
		case 2:
			mode = l2s.SRT_MERGE_MODE_UP
		case 3:
			mode = l2s.SRT_MERGE_MODE_BOTTOM
		}
		srt.Merge(srtT, mode)
	}
	name := opt.Output
	if name == "" {
		title := srt.Title
		if title != "" {
			//以歌曲名命名
			name = fmt.Sprintf("%s.srt", title)
		} else if opt.Id != "" {
			//以歌曲的id命名
			name = fmt.Sprintf("%s.srt", opt.Id)
		} else if opt.Input != "" {
			//以LRC文件的文件名命名
			name = fmt.Sprintf("%s.srt", opt.Input[:len(opt.Input)-4])
		} else {
			//以当前时间的毫秒值命名
			name = fmt.Sprintf("%d.srt", time.Now().Unix())
		}
	} else if !strings.HasSuffix(name, ".srt") {
		name += ".srt"
	}
	if err = srt.WriteFile(name); err != nil {
		fmt.Println("出现错误,保存失败")
		panic(err.Error())
	}
}
