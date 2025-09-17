package main

import (
	xconf "cipstests/chainmaker/common/config"
	"flag"
	"fmt"
	"os/exec"
	"time"

	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string

func init() {
	configFile = flag.String("c", "chainmaker/config/conf6-4.toml", "配置文件路径")
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	clog.Info("-------- 6-4-1. 长安链relayer向布比链转发正确验证信息 --------")
	test_6_4_1()
}

// 长安链relayer向布比链转发正确验证信息
func test_6_4_1() {
	command := "./txtools -c \"chainmaker/config/conf6-1.toml\" -app \"autoResp\" " +
		"-op \"send\" -vf 302 -tp 401 -chain1 30006"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}
