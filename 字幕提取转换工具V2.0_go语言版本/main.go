package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strconv"
	"strings"
)

func printBanner() {
	// 使用ANSI转义码定义颜色
	const (
		cyan    = "\033[36m"
		yellow  = "\033[33m"
		green   = "\033[32m"
		magenta = "\033[35m"
		reset   = "\033[0m"
		bold    = "\033[1m"
	)

	banner := fmt.Sprintf(`
%s███████╗██████╗ ████████╗     ██████╗ ██████╗ ███╗   ██╗██╗   ██╗███████╗██████╗ 
██╔════╝██╔══██╗╚══██╔══╝    ██╔════╝██╔═══██╗████╗  ██║██║   ██║██╔════╝╚════██╗
███████╗██████╔╝   ██║       ██║     ██║   ██║██╔██╗ ██║██║   ██║█████╗    █████╔╝
╚════██║██╔══██╗   ██║       ██║     ██║   ██║██║╚██╗██║╚██╗ ██╔╝██╔══╝    ╚═══██╗
███████║██║  ██║   ██║       ╚██████╗╚██████╔╝██║ ╚████║ ╚████╔╝ ███████╗██████╔╝
╚══════╝╚═╝  ╚═╝   ╚═╝        ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚═════╝ 
                                                                              
%s[ SRT Multilingual Subtitle Converter ]%s
`, cyan, yellow, reset)

	info := fmt.Sprintf(`
%s➤ 版本   : %s1.0
%s➤ 作者   : %sGK

%s▌ 功能概览 ▌───────────────────────────────────────────
  %s✓ 多语言转换   : 支持英/日/韩/繁 → 简中转换
  %s✓ 智能编码检测 : 自动识别UTF-8/GBK等常见编码
  %s✓ 字幕提取     : 从视频中提取字幕
  %sx 批量处理     : 支持目录批量转换
  %s✓ 格式保留     : 完整保持原有时序和样式

%s▌ 支持语言 ▌───────────────────────────────────────────
  %s• 输入语言:     auto：自动识别（识别为一种语言）    zh：简体中文    zh-TW：繁体中文    en：英语    ja：日语    ko：韩语    fr：法语    es：西班牙语    it：意大利语    de：德语    tr：土耳其语    ru：俄语    pt：葡萄牙语    vi：越南语    id：印尼语    th：泰语    ms：马来西亚语    ar：阿拉伯语	hi：印地语
  %s• 输出语言: 简体中文(zh-cn)

%s▌ 技术亮点 ▌───────────────────────────────────────────
  %s⦿ 基于腾讯云机器翻译API
  %s⦿ 专业级SRT格式解析引擎
  %s⦿ 自动错误恢复机制
  %s⦿ 多线程加速处理

%s⚠ 注意事项: 
  • 确保网络连接正常
  • 转换进度实时显示在进度条
%s`,
		green, magenta, green, magenta,
		yellow,
		green, green, green, green, green,
		yellow,
		green, green,
		yellow,
		green, green, green, green,
		yellow, reset,
	)

	fmt.Println(banner)
	fmt.Println(info)
}
func getFuncNumber() int {
	const (
		green = "\033[32m"
		reset = "\033[0m"
	)
	var funcs string
	funcnumber := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		funcsrt := fmt.Sprintf(`
%s1.字幕提取功能
%s2.字幕翻译功能
%s3.退出程序%s`, green, green, green, reset)
		color.Green("请选择需要的功能:")
		fmt.Println(funcsrt)

		funcs, _ = reader.ReadString('\n')
		funcs = strings.TrimSpace(funcs)
		funcnumber, _ = strconv.Atoi(funcs)
		if funcnumber != 1 && funcnumber != 2 && funcnumber != 3 {
			color.Red("你的输入有误，请重新输入")
			continue
		}
		break
	}
	return funcnumber
}

func main() {
	//输出工具版本、功能等信息
	printBanner()
	for {
		//获取需要执行的功能
		func_number := getFuncNumber()
		switch func_number {
		case 1:
			Srt_get()
		case 2:
			Srt_translate()
		case 3:
			os.Exit(0)
		}
	}

}
