package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SubtitleInfo struct {
	Index    string // 流索引
	Language string // 语言代码
	Title    string // 字幕标题
}

func getSrtFile(srtIndexs string, language string, filePath string) (string, error) {
	tmp := filePath + language + ".srt"
	outpath := strings.Replace(tmp, "mkv", "", -1)
	cmd := exec.Command("ffmpeg", "-i", filePath, "-map", srtIndexs, outpath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return outpath, nil
}

// 执行FFmpeg命令获取视频信息
func getFFmpegInfo(filePath string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", filePath)
	output, err := cmd.CombinedOutput()
	//fmt.Println(string(output))
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			fmt.Println("未找到 ffmpeg，请先安装并确保它在 PATH 中！")
			fmt.Println("安装指南：https://ffmpeg.org/download.html")
			return "", nil
		}
	}
	return string(output), nil
}

func parseSubtitles(output string) []SubtitleInfo {
	var subtitles []SubtitleInfo
	lines := strings.Split(output, "\n")

	streamRe := regexp.MustCompile(`Stream #(\d+:\d+).*?\(([a-z]{3})\): Subtitle`)
	titleRe := regexp.MustCompile(`\s+title\s+:\s+(.+)$`)

	var currentSub *SubtitleInfo

	for _, line := range lines {
		// 检测新流开始
		if matches := streamRe.FindStringSubmatch(line); matches != nil {
			if currentSub != nil {
				subtitles = append(subtitles, *currentSub)
			}
			currentSub = &SubtitleInfo{
				Index:    matches[1],
				Language: matches[2],
			}
			continue
		}

		// 检测标题信息
		if currentSub != nil {
			if titleMatch := titleRe.FindStringSubmatch(line); titleMatch != nil {
				currentSub.Title = strings.TrimSpace(titleMatch[1])
			}
		}

		// 检测空行表示流信息结束
		if strings.TrimSpace(line) == "" && currentSub != nil {
			subtitles = append(subtitles, *currentSub)
			currentSub = nil
		}
	}

	// 添加最后一个流
	if currentSub != nil {
		subtitles = append(subtitles, *currentSub)
	}

	return subtitles
}
func Srt_get() {
	var videoPath string
	var srtIndexs string
	var language string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入需要提取字幕的视频文件名：")
	videoPath, _ = reader.ReadString('\n')
	videoPath = strings.TrimSpace(videoPath)
	//videoPath := "What.If.S03E03.mkv" // 替换为你的视频文件路径

	output, err := getFFmpegInfo(videoPath)
	if err != nil {
		log.Fatal("FFmpeg执行错误:", err)
	}

	subtitles := parseSubtitles(output)

	if len(subtitles) == 0 {
		fmt.Println("未找到字幕流")
		return
	}

	fmt.Printf("视频文件: %s\n", videoPath)
	fmt.Println("找到字幕流:")
	index := 1
	for _, sub := range subtitles {
		output2 := fmt.Sprintf("  %-4d: 语言: %-4s", index, sub.Language)
		if sub.Title != "" {
			output2 += fmt.Sprintf(" 标题: %s", sub.Title)
		}
		fmt.Println(output2)
		index++
	}
	for {
		fmt.Println("请输入需要提取的字幕索引(如：1)：")
		srtIndexs, _ = reader.ReadString('\n')
		srtIndexs = strings.TrimSpace(srtIndexs)
		i, _ := strconv.Atoi(srtIndexs)
		if i > len(subtitles) {
			color.Red("你输入的索引有误")
			continue
		}
		subindex := i - 1

		language = subtitles[subindex].Language
		srtIndexs = subtitles[subindex].Index
		outputpath, errout := getSrtFile(srtIndexs, language, videoPath)
		if errout != nil {
			color.Red("字幕提取出错！请检查并重新输入。err:%s", errout)
			continue
		}
		color.Green("字幕提取成功，已保存到:%s\n", outputpath)
		break
	}

}
