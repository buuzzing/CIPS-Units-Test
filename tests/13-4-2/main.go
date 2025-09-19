package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

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
		clog.Info("-------- 13-4-2. 长安链向 LoongChain 发送跨链请求， 且选择 <NUL, 基于公证人的跨链验证, 基础跨链传输协议>--------")
		test_13_4_2()
	}
}

func test_13_4_2() {
	command := "./txtools -c \"chainmaker/config/conf13-4.toml\" -app \"autoResp\" " +
		"-op \"send\" -vf 301 -tp 401 -chain1 1360"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))
}

func check_target() {
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

func check_origin() {
	command := "./txtools -c \"chainmaker/config/conf13-4.toml\" -app \"autoResp\" " +
		"-op \"check\""
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))
}
