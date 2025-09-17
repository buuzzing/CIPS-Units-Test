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

// 随机指定一个 seq
var seq int64

func init() {
	configFile = flag.String("c", "chainmaker/config/conf2-1.toml", "配置文件路径")
	seq = time.Now().Unix()
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	clog.Info("-------- 2-1-1. 长安链向中继链（长安链）发送跨链消息 --------")
	test_2_1_1()
}

// 长安链向中继链（长安链）发送正确的跨链消息
func test_2_1_1() {
	command := "./txtools -c \"chainmaker/config/conf2-1.toml\" -app \"sendMsg\" " +
		"-op \"send\" -vf 300 -tp 401 -chain1 20007"

	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}
