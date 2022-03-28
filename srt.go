package lrc2srt

import (
	"bufio"
	"fmt"
	"github.com/Hami-Lemon/lrc2srt/glist"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SRTMergeMode int

const (
	SRT_MERGE_MODE_STACK SRTMergeMode = iota
	SRT_MERGE_MODE_UP
	SRT_MERGE_MODE_BOTTOM
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

/**
1
00:00:01,111 --> 00:00:10,111
字幕

*/
//返回SRT文件中，一句字幕的字符串表示形式
func (s *SRTContent) String() string {
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(s.Index))
	builder.WriteByte('\n')
	sh, sm, ss, sms := millisecond2Time(s.Start)
	eh, em, es, ems := millisecond2Time(s.End)
	builder.WriteString(fmt.Sprintf("%02d:%02d:%02d,%03d --> %02d:%02d:%02d,%03d\n",
		sh, sm, ss, sms, eh, em, es, ems))
	builder.WriteString(s.Text)
	builder.WriteString("\n\n")
	return builder.String()
}

type SRT struct {
	//歌曲名
	Title string
	//歌手名 未指定文件名是，文件名格式为：歌曲名-歌手名.srt
	Artist  string
	Content glist.Queue[*SRTContent]
}

// LrcToSrt LRC对象转换成SRT对象
func LrcToSrt(lrc *LRC) *SRT {
	if lrc == nil {
		return nil
	}
	srt := &SRT{
		Title:   lrc.Title,
		Artist:  lrc.Artist,
		Content: glist.NewLinkedList[*SRTContent](),
	}
	index := 1
	//上一条srt信息
	var prevSRT *SRTContent
	for it := lrc.LrcList.Iterator(); it.Has(); {
		lrcNode := it.Next()
		srtContent := &SRTContent{
			Index: index,
			Start: lrcNode.time,
			Text:  lrcNode.content,
		}
		if index != 1 {
			//上一条歌词的结束时间设置为当前歌词的开始时间
			prevSRT.End = srtContent.Start
		}
		srt.Content.PushBack(srtContent)
		index++
		prevSRT = srtContent
	}
	//最后一条歌词
	if prevSRT != nil {
		//结束时间是为其 开始时间+10 秒
		prevSRT.End = prevSRT.Start + 1000
	}
	return srt
}

// Merge 将另一个srt信息合并到当前srt中,有三种合并模式
//1. SRT_MERGE_MODE_STACK: 按照开始时间对两个srt信息进行排序,交错合并
//2. SRT_MERGE_MODE_UP: 当前srt信息排列在上,另一个排列在下,即 other 追加到后面
//3. SRT_MERGE_MODE_BOTTOM: 当前srt信息排列在下,另一个排列在上,即 other 添加到前面
func (s *SRT) Merge(other *SRT, mode SRTMergeMode) {
	switch mode {
	case SRT_MERGE_MODE_STACK:
		s.mergeStack(other)
	case SRT_MERGE_MODE_UP:
		s.mergeUp(other)
	case SRT_MERGE_MODE_BOTTOM:
		s.mergeBottom(other)
	}
}

//todo 改进算法,现在的算法太慢了
func (s *SRT) mergeStack(other *SRT) {
	size := s.Content.Size() + other.Content.Size()
	temp := make([]*SRTContent, size, size)
	index := 0
	for it := s.Content.Iterator(); it.Has(); {
		temp[index] = it.Next()
		index++
	}
	for it := other.Content.Iterator(); it.Has(); {
		temp[index] = it.Next()
		index++
	}
	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].Start < temp[j].Start
	})
	list := glist.NewLinkedList[*SRTContent]()
	for _, v := range temp {
		list.Append(v)
	}
	s.Content = list
}

func (s *SRT) mergeUp(other *SRT) {
	if other.Content.IsNotEmpty() {
		for it := other.Content.Iterator(); it.Has(); {
			s.Content.Append(it.Next())
		}
	}
}

func (s *SRT) mergeBottom(other *SRT) {
	oq := other.Content
	if oq.IsNotEmpty() {
		s.Content.PushFront(*(oq.PullBack()))
	}
}

// WriteFile 将SRT格式的数据写入指定的文件中
func (s *SRT) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = s.Write(f)
	err = f.Close()
	return err
}

// Write 将SRT格式的数据写入dst中
func (s *SRT) Write(dst io.Writer) error {
	//6KB的缓冲
	bufSize := 1024 * 6
	writer := bufio.NewWriterSize(dst, bufSize)
	for it := s.Content.Iterator(); it.Has(); {
		_, err := writer.WriteString(it.Next().String())
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
