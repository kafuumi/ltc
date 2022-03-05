package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

type SRTContent struct {
	//序号，从1开始
	Index int
	//开始时间，单位毫秒
	Start int
	//结束时间，单位毫秒
	End int
	//歌词内容
	Text string
}
type SRT struct {
	//歌曲名
	Title string
	//歌手名 未指定文件名是，文件名格式为：歌曲名-歌手名.srt
	Artist  string
	Content []*SRTContent
}

// Option 运行时传入的选项
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
	VERSION = `"0.1.1" (build 2022.03.05)`
)

var (
	ChromeUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.82 Safari/537.36"
	client   = http.Client{}
	opt      Option
)

func main() {
	args, err := flags.Parse(&opt)
	if err != nil {
		os.Exit(0)
	}
	//显示版本信息
	if opt.Version {
		fmt.Printf("LrcToSrt(lts) version %s\n", VERSION)
		os.Exit(0)
	}
	//获取保存的文件名
	if len(args) != 0 {
		opt.Output = args[0]
	}

	//获取歌词，lyric为原文歌词，tranLyric为译文歌词
	var lyric, tranLyric string
	if opt.Id != "" {
		if opt.Source != "163" {
			lyric, tranLyric = GetQQLyric(opt.Id)
		} else {
			lyric, tranLyric = Get163Lyric(opt.Id)
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
			WriteFile(o, lyric)
			if tranLyric != "" {
				WriteFile("tran_"+o, tranLyric)
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
		lyric = ReadFile(opt.Input)
		if lyric == "" {
			fmt.Println("获取歌词失败，文件内容为空。")
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: 请指定需要转换的歌词。")
		os.Exit(1)
	}
	//原文和译文作为两条歌词流信息分开保存，但最终生成的srt文件会同时包含两个信息
	lyricSRT, tranLyricSRT := Lrc2Srt(lyric), Lrc2Srt(tranLyric)
	SaveSRT(lyricSRT, tranLyricSRT, opt.Output)
}

// SaveSRT 保存数据为SRT文件
func SaveSRT(srt *SRT, tranSrt *SRT, name string) {
	if tranSrt == nil {
		//没有译文时，用一个空的对象的代替，减少nil判断
		tranSrt = &SRT{Content: make([]*SRTContent, 0)}
		//因为没有译文，所以mode选项无效，设为2之后，后面不用做多余判断
		opt.Mode = 2
	}
	//处理结果文件的文件名
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

	file, err := os.Create(name)
	if err != nil {
		fmt.Printf("创建结果文件失败：%v\n", err)
		os.Exit(1)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	//6KB的缓存区，大部分歌词生成的SRT文件均在4-6kb左右
	writer := bufio.NewWriterSize(file, 1024*6)

	/*
		原文和译文歌词的排列方式，因为原文歌词中可能包含一些非歌词信息，
		例如作词者，作曲者等，而在译文歌词中却可能不包含这些
	*/
	//srt的序号
	index := 1
	switch opt.Mode {
	case 1:
		//将两个歌词合并成一个新的数组
		size := len(srt.Content) + len(tranSrt.Content)
		temp := make([]*SRTContent, size, size)
		i := 0
		for _, v := range srt.Content {
			temp[i] = v
			i++
		}
		for _, v := range tranSrt.Content {
			temp[i] = v
			i++
		}
		//按开始时间进行排序，使用SliceStable确保一句歌词的原文在译文之前
		sort.SliceStable(temp, func(i, j int) bool {
			return temp[i].Start < temp[j].Start
		})
		//写入文件
		for i, v := range temp {
			v.Index = i + 1
			_, _ = writer.WriteString(v.String())
		}
	case 2:
		//原文在上，译文在下
		for _, item := range srt.Content {
			item.Index = index
			_, _ = writer.WriteString(item.String())
			index++
		}

		for _, item := range tranSrt.Content {
			item.Index = index
			_, _ = writer.WriteString(item.String())
			index++
		}
	case 3:
		//译文在上，原文在下
		for _, item := range tranSrt.Content {
			item.Index = index
			_, _ = writer.WriteString(item.String())
			index++
		}
		for _, item := range srt.Content {
			item.Index = index
			_, _ = writer.WriteString(item.String())
			index++
		}
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("保存结果失败：%v\n", err)
	} else {
		fmt.Printf("转换文件完成，保存结果为：%s\n", name)
	}
}

// Lrc2Srt 将原始个LRC字符串歌词解析SRT对象
func Lrc2Srt(src string) *SRT {
	if src == "" {
		return nil
	}
	//标准的LRC文件为一行一句歌词
	lyrics := strings.Split(src, "\n")
	//标识标签的正则 [ar:A-SOUL]形式
	infoRegx := regexp.MustCompile(`^\[([a-z]+):([\s\S]*)]`)

	srt := &SRT{Content: make([]*SRTContent, 0, len(lyrics))}
	//解析标识信息
	for {
		if len(lyrics) == 0 {
			break
		}
		l := lyrics[0]
		//根据正则表达式进行匹配
		info := infoRegx.FindStringSubmatch(l)
		//标识信息位于歌词信息前面，当出现未匹配成功时，即可退出循环
		if info != nil {
			//info 中为匹配成功的字符串和 子组合（正则表达式中括号包裹的部分）
			//例如，对于标识信息：[ar:A-SOUL]，info中的数据为[[ar:A-SOUL] ar A-SOUL]
			key := info[1]
			switch key {
			case "ar":
				//歌手名
				if len(info) == 3 {
					srt.Artist = info[2]
				}
			case "ti":
				//歌曲名
				if len(info) == 3 {
					srt.Title = info[2]
				}
			}
			lyrics = lyrics[1:]
		} else {
			break
		}
	}
	//歌词信息的正则，"[00:10.222]超级的敏感"或“[00:10:222]超级的敏感”或“[00:10]超级的敏感”或“[00:10.22]超级的敏感”或“[00:10:22]超级的敏感”
	lyricRegx := regexp.MustCompile(`\[(\d\d):(\d\d)([.:]\d{2,3})?]([\s\S]+)`)
	index := 0
	for _, l := range lyrics {
		content := lyricRegx.FindStringSubmatch(l)
		if content != nil {
			c := SplitLyric(content[1:])
			if c != nil {
				if index != 0 {
					//前一条字幕的结束时间为当前字幕开始的时间
					srt.Content[index-1].End = c.Start
				}
				srt.Content = append(srt.Content, c)
				index++
			}
		}
	}
	//最后一条字幕
	last := srt.Content[index-1]
	//最后一条字幕的结束时间为其开始时间 + 10秒
	last.End = last.Start + 10000
	return srt
}

// SplitLyric 对分割出来的歌词信息进行解析
func SplitLyric(src []string) *SRTContent {
	minute, err := strconv.Atoi(src[0])
	second, err := strconv.Atoi(src[1])
	if err != nil {
		fmt.Printf("错误的时间格式：%s\n", src)
		return nil
	}
	millisecond, content := 0, ""

	_len := len(src)
	if _len == 3 {
		//歌词信息没有毫秒值
		content = src[2]
	} else if _len == 4 {
		content = src[3]
		//字符串的第一个字符是 "." 或 ":"
		ms := src[2][1:]
		millisecond, err = strconv.Atoi(ms)
		//QQ音乐歌词文件中，毫秒值只有两位，需要特殊处理一下
		if len(ms) == 2 {
			millisecond *= 10
		}
		if err != nil {
			fmt.Printf("错误的时间格式：%s\n", src)
			return nil
		}
	}

	srtContent := &SRTContent{}
	srtContent.Start = Time2Millisecond(minute, second, millisecond)
	srtContent.Text = content
	return srtContent
}

//返回SRT文件中，一句字幕的字符串表示形式
/**
1
00:00:01,111 --> 00:00:10,111
字幕

*/
func (s *SRTContent) String() string {
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(s.Index))
	builder.WriteByte('\n')
	sh, sm, ss, sms := Millisecond2Time(s.Start)
	eh, em, es, ems := Millisecond2Time(s.End)
	builder.WriteString(fmt.Sprintf("%02d:%02d:%02d,%03d --> %02d:%02d:%02d,%03d\n",
		sh, sm, ss, sms, eh, em, es, ems))
	builder.WriteString(s.Text)
	builder.WriteString("\n\n")

	return builder.String()
}
