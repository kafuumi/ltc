package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type CloudLyricBase struct {
	Version int    `json:"version"`
	Lyric   string `json:"lyric"`
}

type CloudLyric struct {
	Sgc       bool           `json:"sgc"`
	Sfy       bool           `json:"sfy"`
	Qfy       bool           `json:"qfy"`
	TransUser interface{}    `json:"transUser"`
	Lrc       CloudLyricBase `json:"lrc"`
	TLyric    CloudLyricBase `json:"tlyric"`
	Code      int            `json:"code"`
}

func Get163Lyric(id string) (lyric, tLyric string) {
	api := "https://music.163.com/api/song/lyric"
	params := url.Values{}
	params.Add("os", "pc")
	//歌曲的id号
	params.Add("id", id)
	//包含原始歌词
	params.Add("lv", "1")
	//包含翻译歌词
	params.Add("tv", "1")
	api = fmt.Sprintf("%s?%s", api, params.Encode())

	req, _ := http.NewRequest("GET", api, nil)
	//必须设置Referer,否则会请求失败
	req.Header.Add("Referer", "https://music.163.com")
	req.Header.Add("User-Agent", ChromeUA)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("网络错误:%v\n", err)
		os.Exit(1)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("网络请求失败，状态码为：%d\n", resp.StatusCode)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var cloudLyric CloudLyric
	err = json.NewDecoder(resp.Body).Decode(&cloudLyric)
	if cloudLyric.Sgc {
		fmt.Printf("获取歌词失败，返回的结果为：%+v，请检查id是否正确\n", cloudLyric)
		os.Exit(1)
	}
	return cloudLyric.Lrc.Lyric, cloudLyric.TLyric.Lyric
}
