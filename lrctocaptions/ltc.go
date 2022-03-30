package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/Hami-Lemon/ltc"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// VERSION 当前版本
	VERSION      = `"0.3.4" (build 2022.03.30)`
	VERSION_INFO = "LrcToCaptions(ltc) version: %s\n"
)

func main() {
	parseFlag()
	//TODO 酷狗的krc精准到字，更利于打轴 https://shansing.com/read/392/
	//显示版本信息
	if version.IsSet() {
		fmt.Printf(VERSION_INFO, VERSION)
		return
	}
	//未指定来源
	if input == "" {
		fmt.Printf("未指定歌词来源\n")
		flag.Usage()
		os.Exit(0)
	}
	//获取歌词，lyric为原文歌词，tranLyric为译文歌词
	var lyric, tranLyric string
	//从文件中获取
	if checkPath(input) {
		if !strings.HasSuffix(input, ".lrc") {
			fmt.Println("Error: 不支持的格式，目前只支持lrc歌词文件。")
			panic("")
		}
		if data, err := ioutil.ReadFile(input); err == nil {
			if len(data) == 0 {
				fmt.Println("获取歌词失败，文件内容为空。")
				panic("")
			}
			lyric = string(data)
		} else {
			panic("读取文件失败:" + input + err.Error())
		}
	} else {
		//从网络上获取
		if source != "163" {
			lyric, tranLyric = ltc.GetQQLyric(input)
		} else {
			lyric, tranLyric = ltc.Get163Lyric(input)
		}
	}

	//下载歌词
	if download.IsSet() {
		//对文件名进行处理
		o := output
		if o == "" {
			o = input + ".lrc"
		} else if !strings.HasSuffix(o, ".lrc") {
			o += ".lrc"
		}
		writeFile(o, lyric)
		if tranLyric != "" {
			writeFile("tran_"+o, tranLyric)
		}
		fmt.Println("下载歌词完成！")
		return
	}

	lrc, lrcT := ltc.ParseLRC(lyric), ltc.ParseLRC(tranLyric)
	//先转换成srt
	srt, srtT := ltc.LrcToSrt(lrc), ltc.LrcToSrt(lrcT)
	if srtT != nil {
		//原文和译文合并
		srt.Merge(srtT, mode.Mode())
	}
	switch format.Value() {
	case FORMAT_SRT:
		if err := srt.WriteFile(output); err != nil {
			fmt.Println("出现错误,保存失败")
			panic(err.Error())
		}
	case FORMAT_ASS:
		ass := ltc.SrtToAss(srt)
		if err := ass.WriteFile(output); err != nil {
			fmt.Println("出现错误,保存失败")
			panic(err.Error())
		}
	}
	//如果是相对路径，则获取其对应的绝对路径
	if !filepath.IsAbs(output) {
		//如果是相对路径，父目录即是当前运行路径
		dir, er := os.Getwd()
		if er == nil {
			output = dir + string(os.PathSeparator) + output
		}
	}
	fmt.Printf("保存结果为：%s\n", output)
}

func writeFile(file string, content string) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("创建结果文件失败：%v\n", err)
		panic("")
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(content)
	err = writer.Flush()
	if err != nil {
		fmt.Printf("保存文件失败：%v\n", err)
		panic("")
	}
}
