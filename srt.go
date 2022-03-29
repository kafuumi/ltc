package ltc

import (
	"bufio"
	"fmt"
	"github.com/Hami-Lemon/ltc/glist"
	"io"
	"os"
	"path/filepath"
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
	sIt, oIt := s.Content.Iterator(), other.Content.Iterator()
	//不对原来的链表做修改，合并的信息保存在一个新的链表中
	merge := glist.NewLinkedList[*SRTContent]()
	var sNode, oNode *SRTContent
	//分别获取两个链表的第一个元素
	if sIt.Has() && oIt.Has() {
		sNode, oNode = sIt.Next(), oIt.Next()
	}
	//开始迭代
	for sIt.Has() && oIt.Has() {
		//小于等于，当相等时，s中的元素添加进去
		if sNode.Start <= oNode.Start {
			merge.Append(sNode)
			sNode = sIt.Next()
		} else {
			merge.Append(oNode)
			oNode = oIt.Next()
		}
	}
	if sNode != nil && oNode != nil {
		//循环退出时，sNode和oNode指向的元素还没有进行比较，会导致缺少两条数据
		if sNode.Start <= oNode.Start {
			merge.Append(sNode)
			merge.Append(oNode)
		} else {
			merge.Append(oNode)
			merge.Append(sNode)
		}
	}
	//剩下的元素添加到链表中，最多只有一个链表有剩余元素
	for sIt.Has() {
		merge.Append(sIt.Next())
	}
	for oIt.Has() {
		merge.Append(oIt.Next())
	}
	s.Content = merge
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
		//不存在对应文件夹
		if os.IsNotExist(err) {
			panic("文件夹不存在:" + filepath.Dir(path))
		}
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
