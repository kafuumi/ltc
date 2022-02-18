package main

import (
	"fmt"
	"io"
	"os"
)

// Time2Millisecond 根据分，秒，毫秒 计算出对应的毫秒值
func Time2Millisecond(m, s, ms int) int {
	t := m*60 + s
	t *= 1000
	t += ms
	return t
}

// Millisecond2Time 根据毫秒值计算出对应的 时，分，秒，毫秒形式的时间值
func Millisecond2Time(millisecond int) (h, m, s, ms int) {
	ms = millisecond % 1000

	s = millisecond / 1000
	m = s / 60
	h = m / 60

	s %= 60
	m %= 60
	return
}

func ReadFile(name string) string {
	if name == "" {
		return ""
	}
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("打开文件失败：%v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("读取文件失败：%v\n", err)
		os.Exit(1)
	}
	return string(data)
}

func WriteFile(name string, data string) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Printf("保存文件失败：%v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	nw, err := file.WriteString(data)
	if err != nil || nw < len(data) {
		fmt.Printf("保存文件失败：%v\n", err)
		os.Exit(1)
	}
}
