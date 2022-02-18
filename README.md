# LrcToSrt
用于将LRT歌词文件转换成SRT字幕文件

## 功能
- [x] lrc文件转换成srt文件
- [x] 从网易云音乐或QQ音乐上获取歌词，并转换成srt文件
- [x] 从网易云音乐或QQ音乐上下载歌词

## 开始使用

```
Usage:
  D:\ProgrameStudy\lrc2srt\lts.exe [OPTIONS]

Application Options:
  -i, --id=                歌曲的id，网易云和QQ音乐均可。
  -I, --input=             需要转换的LRC文件路径。
  -s, --source=[163|qq|QQ] 当设置id时有效，指定从网易云（163）还是QQ音乐（qq）上获取歌词。
                           (default: 163)
  -d, --download           只下载歌词，而不进行解析。
  -m, --mode=[1|2|3]       原文和译文的排列模式,可选值有：[1] [2] [3] (default: 1)
  -v, --version            获取版本信息

Help Options:
  -h, --help               Show this help message
```
