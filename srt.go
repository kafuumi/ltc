package lrc2srt

import (
	"bufio"
	"fmt"
	"github.com/Hami-Lemon/lrc2srt/glist"
	"io"
	"os"
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
	//序号，从1开始,只在写入文件的时候设置这个属性
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
			Index: 0,
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

//可以类比为合并两个有序链表,
//算法参考:https://leetcode-cn.com/problems/merge-two-sorted-lists/solution/he-bing-liang-ge-you-xu-lian-biao-by-leetcode-solu/
func (s *SRT) mergeStack(other *SRT) {
	ls, _ := s.Content.(*glist.LinkedList[*SRTContent])
	lo, _ := other.Content.(*glist.LinkedList[*SRTContent])
	lhs, lho := ls.First, lo.First
	//临时的空结点
	preHead := &glist.Node[*SRTContent]{}
	prev := preHead

	for lhs != nil && lho != nil {
		//副本结点,从而不改变other对象中的内容
		oCopy := lho.Clone()

		if lhs.Element.Start <= oCopy.Element.Start {
			prev.Next = lhs
			lhs.Prev = prev
			lhs = lhs.Next
		} else {
			prev.Next = oCopy
			oCopy.Prev = prev
			lho = lho.Next
		}
		prev = prev.Next
	}
	if lhs == nil {
		//如果剩下的内容是other中的,则依次迭代复制到s中
		for n := lho; n != nil; n = n.Next {
			c := n.Clone()
			prev.Next = c
			c.Prev = prev
			prev = prev.Next
		}
	} else {
		prev.Next = lhs
	}
	ls.First = preHead.Next
	ls.First.Prev = nil
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
		for it := other.Content.ReverseIterator(); it.Has(); {
			s.Content.PushFront(it.Next())
		}
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
	index := 1
	for it := s.Content.Iterator(); it.Has(); {
		content := it.Next()
		content.Index = index
		index++
		_, err := writer.WriteString(content.String())
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
