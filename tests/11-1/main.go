package main

import (
	xconf "cipstests/chainmaker/common/config"
	"flag"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string

// 随机指定一个 seq
var seq int64

func init() {
	configFile = flag.String("c", "chainmaker/config/conf11-1.toml", "配置文件路径")

	seq = int64(rand.Intn(100000000))
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	clog.Info("-------- 11-1-1. 正常执行并完成事务提交 --------")
	test_11_1_1()
	clog.Info("-------- 11-1-2. 状态异常并完成事务回滚 --------")
	test_11_1_2()
}

// 正常执行并完成事务提交
func test_11_1_1() {
	command := "./txtools -c \"chainmaker/config/conf11-1.toml\" -app \"atomic\" " +
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

// 状态异常并完成事务回滚
func test_11_1_2() {
	command := "./txtools -c \"chainmaker/config/conf11-1.toml\" -app \"atomic\" " +
		"-op \"rollback\" -vf 300 -tp 401 -chain1 20007"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}
