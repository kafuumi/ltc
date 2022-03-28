package lrc2srt

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Hami-Lemon/lrc2srt/glist"
)

type LRCNode struct {
	//歌词出现的时间,单位毫秒
	time int
	//歌词内容
	content string
}

type LRC struct {
	//歌曲名
	Title string
	//歌手名
	Artist string
	//专辑名
	Album string
	//歌词作者
	Author string
	//歌词列表
	LrcList glist.Queue[*LRCNode]
}

func ParseLRC(src string) *LRC {
	if src == "" {
		return nil
	}
	//标准的LRC文件为一行一句歌词
	lyrics := strings.Split(src, "\n")
	//标识标签的正则 [ar:A-SOUL]形式
	infoRegx := regexp.MustCompile(`^\[([a-z]+):([\s\S]*)]`)

	lrc := &LRC{LrcList: glist.NewLinkedList[*LRCNode]()}
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
			//例如，对于标识信息：[ar:A-SOUL]，info中的数据为[[ar:A-SOUL], ar, A-SOUL]
			key := info[1]
			switch key {
			case "ar":
				//歌手名
				if len(info) == 3 {
					lrc.Artist = info[2]
				}
			case "ti":
				//歌曲名
				if len(info) == 3 {
					lrc.Title = info[2]
				}
			case "al":
				//专辑名
				if len(info) == 3 {
					lrc.Album = info[2]
				}
			case "by":
				//歌词作者
				if len(info) == 3 {
					lrc.Author = info[2]
				}
			}
			lyrics = lyrics[1:]
		} else {
			break
		}
	}
	//歌词信息的正则，"[00:10.222]超级的敏感"或“[00:10:222]超级的敏感”或“[00:10]超级的敏感”或“[00:10.22]超级的敏感”或“[00:10:22]超级的敏感”
	lyricRegx := regexp.MustCompile(`\[(\d\d):(\d\d)([.:]\d{2,3})?]([\s\S]+)`)
	for _, l := range lyrics {
		content := lyricRegx.FindStringSubmatch(l)
		if content != nil {
			node := SplitLyric(content[1:])
			if node != nil {
				lrc.LrcList.PushBack(node)
			}
		}
	}
	return lrc
}

// SplitLyric 对分割出来的歌词信息进行解析
func SplitLyric(src []string) *LRCNode {
	minute, err := strconv.Atoi(src[0])
	second, err := strconv.Atoi(src[1])
	if err != nil {
		panic("错误的时间格式:" + strings.Join(src, " "))
		return nil
	}
	millisecond, content := 0, ""

	_len := len(src)
	if _len == 3 {
		//歌词信息没有毫秒值
		content = src[2]
	} else if _len == 4 {
		content = src[3]
		//毫秒字符串的第一个字符是 "." 或 ":",需要去掉
		ms := src[2][1:]
		millisecond, err = strconv.Atoi(ms)
		//QQ音乐歌词文件中，毫秒值只有两位，需要特殊处理一下
		if len(ms) == 2 {
			millisecond *= 10
		}
		if err != nil {
			panic("错误的时间格式:" + strings.Join(src, " "))
			return nil
		}
	}
	lrcNode := &LRCNode{
		time:    time2Millisecond(minute, second, millisecond),
		content: content,
	}
	return lrcNode
}
