package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	clog "github.com/kpango/glg"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "check-origin":
			check_origin()
		case "check-target":
			check_target()
		default:
			fmt.Println("无效的参数。使用 'check-origin' 或 'check-target'")
			return
		}
	} else {
		clog.Info("-------- 13-4-1. LoongChain 向长安链发送跨链请求, 且选择 <NUL, 基于公证人的跨链验证, 基础跨链传输协议>--------")
		test_13_4_1()
	}
}

func test_13_4_1() {
	msg := "Hello" + time.Now().Format("20060102150405")
	clog.Infof("发送跨链消息: %s", msg)
	// 构造请求体
	payload := map[string]interface{}{
		"sourceChainID":  "1360",
		"targetChainID":  "20006",
		"sourceAppID":    "104",
		"targetAppID":    "104",
		"transactionId":  "0",
		"verificationId": "301",
		"transportId":    "401",
		"ack":            true,
		"appArgs":        msg,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("JSON 序列化失败: %v", err))
	}

	req, err := http.NewRequest("POST", "http://220.189.210.171:28090/api/send_cc_msg", bytes.NewBuffer(body))
	if err != nil {
		panic(fmt.Sprintf("创建请求失败: %v", err))
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "172.30.10.55:8080")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("发送请求失败: %v", err))
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	clog.Infof("HTTP 响应: %s", respBody.String())

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}

func check_origin() {
	url := "http://220.189.210.171:28090/api/app_check/1360"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("创建请求失败: %v", err))
	}

	req.Header.Set("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "172.30.10.55:8080")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("发送请求失败: %v", err))
	}
	defer resp.Body.Close()

	respBody := new(bytes.Buffer)
	respBody.ReadFrom(resp.Body)
	clog.Infof("HTTP 响应: %s", respBody.String())
}

func check_target() {
	command := "./txtools -c \"chainmaker/config/conf13-4.toml\" -app \"autoResp\" " +
		"-op \"check\""
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))
}
