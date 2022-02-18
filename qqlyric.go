package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type QQLyric struct {
	RetCode int    `json:"retcode"`
	Code    int    `json:"code"`
	SubCode int    `json:"subcode"`
	Lyric   string `json:"lyric"`
	Trans   string `json:"trans"`
}

func GetQQLyric(id string) (lyric, tLyric string) {
	api := "https://c.y.qq.com/lyric/fcgi-bin/fcg_query_lyric_new.fcg"
	params := url.Values{}
	//返回格式
	params.Add("format", "json")
	params.Add("inCharset", "utf-8")
	params.Add("outCharset", "utf-8")
	params.Add("platform", "yqq.json")
	params.Add("g_tk", "5381")
	//歌曲的id号
	params.Add("songmid", id)
	//返回结果为原始结果，而不是base64编码的结果（base64编码后数据量会增大）
	params.Add("nobase64", "1")

	api = fmt.Sprintf("%s?%s", api, params.Encode())

	req, _ := http.NewRequest("GET", api, nil)
	//必须设置Referer,否则会请求失败
	req.Header.Add("Referer", "https://y.qq.com")
	req.Header.Add("User-Agent", ChromeUA)
	req.Header.Add("accept-encoding", "gzip")
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

	//返回的数据是gzip压缩，需要解压
	reader, _ := gzip.NewReader(resp.Body)
	var qqLyric QQLyric
	err = json.NewDecoder(reader).Decode(&qqLyric)
	if qqLyric.RetCode != 0 {
		fmt.Printf("获取歌词失败，返回的结果为：%+v，请检查id是否正确\n", qqLyric)
		os.Exit(1)
	}
	return qqLyric.Lyric, qqLyric.Trans
}
