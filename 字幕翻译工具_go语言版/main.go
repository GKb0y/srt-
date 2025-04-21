package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	endpoint   = "tmt.tencentcloudapi.com"
	service    = "tmt"
	region     = "ap-beijing"
	apiVersion = "2018-03-21"
	maxRetries = 3                       // 最大重试次数
	reqPerSec  = 5                       // 每秒最大请求数
	interval   = time.Second / reqPerSec // 每个请求间隔
)

type Subtitle struct {
	Index   int
	Start   string
	End     string
	Content string
}

type TencentTranslator struct {
	SecretID  string
	SecretKey string
	Client    *http.Client
}

func NewTencentTranslator(secretID, secretKey string) *TencentTranslator {
	return &TencentTranslator{
		SecretID:  secretID,
		SecretKey: secretKey,
		Client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (t *TencentTranslator) generateSignature(timestamp int64, payload string) string {
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	keyDate := hmacSHA256("TC3"+t.SecretKey, date)
	keyService := hmacSHA256(string(keyDate), service)
	keySigning := hmacSHA256(string(keyService), "tc3_request")

	canonicalRequest := fmt.Sprintf(
		"POST\n/\n\ncontent-type:application/json\nhost:%s\n\ncontent-type;host\n%s",
		endpoint,
		sha256Hex(payload),
	)

	stringToSign := fmt.Sprintf(
		"TC3-HMAC-SHA256\n%d\n%s/%s/tc3_request\n%s",
		timestamp,
		date,
		service,
		sha256Hex(canonicalRequest),
	)

	signature := hmacSHA256(string(keySigning), stringToSign)
	return fmt.Sprintf(
		"TC3-HMAC-SHA256 Credential=%s/%s/%s/tc3_request, SignedHeaders=content-type;host, Signature=%s",
		t.SecretID,
		date,
		service,
		hex.EncodeToString(signature),
	)
}

func (t *TencentTranslator) Translate(text string, source, target string) (string, error) {
	credential := common.NewCredential(t.SecretID, t.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqTimeout = 30
	client, _ := tmt.NewClient(credential, region, cpf)

	request := tmt.NewTextTranslateRequest()
	request.SourceText = common.StringPtr(text)
	request.Source = common.StringPtr(source)
	request.Target = common.StringPtr(target)
	request.ProjectId = common.Int64Ptr(0)

	response, err := client.TextTranslate(request)
	if err != nil {
		return "", fmt.Errorf("翻译失败: %w", err)
	}
	return *response.Response.TargetText, nil
}

func parseSRT(content string) []Subtitle {
	blocks := regexp.MustCompile(`\n\n+`).Split(content, -1)
	var subs []Subtitle

	for _, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) < 3 {
			continue
		}

		timeParts := strings.Split(lines[1], " --> ")
		if len(timeParts) != 2 {
			continue
		}

		subs = append(subs, Subtitle{
			Index:   len(subs) + 1,
			Start:   strings.TrimSpace(timeParts[0]),
			End:     strings.TrimSpace(timeParts[1]),
			Content: strings.Join(lines[2:], "\n"),
		})
	}
	return subs
}

func composeSRT(subs []Subtitle) string {
	var buf strings.Builder
	for _, sub := range subs {
		buf.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n",
			sub.Index, sub.Start, sub.End, sub.Content))
	}
	return buf.String()
}

func getCredentials() (string, string, error) {
	file, err := os.Open("key.txt")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	keys := make(map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "secret_id=") {
			keys["secret_id"] = strings.TrimPrefix(line, "secret_id=")
		} else if strings.HasPrefix(line, "secret_key=") {
			keys["secret_key"] = strings.TrimPrefix(line, "secret_key=")
		}
	}

	if keys["secret_id"] == "" || keys["secret_key"] == "" {
		return "", "", fmt.Errorf("missing keys in file")
	}

	return keys["secret_id"], keys["secret_key"], scanner.Err()
}

func hmacSHA256(key, data string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
func createProgressBar(total int) *pb.ProgressBar {
	return pb.New(total).
		SetTemplateString(`{{ cyan "🚀 翻译进度:" }} {{percent . }} {{ bar . "█" "░" "▒" " " "█"}} {{counters . }} {{rtime . "%s"}}`).
		Set(pb.Bytes, false).
		SetRefreshRate(200*time.Millisecond).
		SetWidth(100).
		SetWidth(120).
		Set(pb.Terminal, true).
		Start()
}
func waitForEnter() {
	done := make(chan bool)

	// 倒计时协程
	go func() {
		for i := 5; i > 0; i-- {
			fmt.Printf("\r程序将自动在 %-2d 秒后退出（按回车立即退出）...", i)
			time.Sleep(1 * time.Second)
		}
		done <- true
	}()

	// 等待回车
	go func() {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		done <- true
	}()

	// 主线程阻塞，任一条件触发退出
	<-done
	fmt.Println("\r程序已退出。                                                                ") // 清理行尾
}
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
		green, green, green, green,
		yellow,
		green, green,
		yellow,
		green, green, green, green,
		yellow, reset,
	)

	fmt.Println(banner)
	fmt.Println(info)
}
func main() {
	printBanner()
	var language string
	var inputFile string
	var outputFile string
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("请输入字幕文件的语言(默认为auto)：")
	language, _ = reader.ReadString('\n')
	language = strings.TrimSpace(language)
	if language == "" {
		language = "auto"
	}
	for {
		fmt.Println("请输入需要转换的字幕文件：")
		reader.Reset(os.Stdin)
		inputFile, _ = reader.ReadString('\n')
		inputFile = strings.TrimSpace(inputFile)
		if inputFile == "" {
			fmt.Println("未输入需要转换的文件")
			continue
		} else {
			break
		}
	}

	fmt.Println("请输入输出文件(默认为out.srt)：")
	reader.Reset(os.Stdin)
	outputFile, _ = reader.ReadString('\n')
	outputFile = strings.TrimSpace(outputFile)
	if outputFile == "" {
		outputFile = "out.srt"
	}

	secretID, secretKey, err := getCredentials()
	if err != nil {
		fmt.Printf("读取密钥失败: %v\n", err)
		os.Exit(1)
	}

	translator := NewTencentTranslator(secretID, secretKey)

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("读取文件失败: %v\n", err)
		os.Exit(1)
	}

	subs := parseSRT(string(content))

	bar := createProgressBar(len(subs))

	rateLimiter := time.Tick(interval) // 速率限制器
	batchSize := 3                     // 减少批处理大小
	var translatedSubs []Subtitle

	for i := 0; i < len(subs); {
		select {
		case <-rateLimiter:
			end := i + batchSize
			if end > len(subs) {
				end = len(subs)
			}

			var texts []string
			for _, sub := range subs[i:end] {
				texts = append(texts, sub.Content)
			}

			combined := strings.Join(texts, "\n")
			var result string
			var err error

			// 重试逻辑
			for retry := 0; retry < maxRetries; retry++ {
				result, err = translator.Translate(combined, language, "zh")
				if err == nil {
					break
				}

				if strings.Contains(err.Error(), "RequestLimitExceeded") {
					waitTime := time.Duration(retry+1) * time.Second
					//fmt.Printf("\n触发频率限制，等待1秒后重试...\n")
					time.Sleep(waitTime)
					continue
				}
				break
			}

			if err != nil {
				fmt.Printf("\n最终翻译失败: %v\n", err)
				result = combined // 保留原文
			}

			results := strings.Split(result, "\n")
			for j, sub := range subs[i:end] {
				newContent := sub.Content
				if j < len(results) {
					newContent = results[j]
				}
				translatedSubs = append(translatedSubs, Subtitle{
					Index:   sub.Index,
					Start:   sub.Start,
					End:     sub.End,
					Content: newContent,
				})
			}
			bar.Add(end - i)
			i = end
		}
	}

	bar.Finish()
	fmt.Print("\033[2K\r") // 清除最后一行

	if err := os.WriteFile(outputFile, []byte(composeSRT(translatedSubs)), 0644); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n转换完成! 结果已保存到 %s\n", outputFile)
	// 使用示例（30秒超时）
	waitForEnter()
}
