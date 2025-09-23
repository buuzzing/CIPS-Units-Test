package main

import (
	"fmt"
	"os/exec"

	clog "github.com/kpango/glg"
)

func main() {
	clog.Info("-------- 11-4-1. 正常执行并完成事务提交 --------")
	test_11_4_1()
	clog.Info("-------- 11-4-2. 状态异常并完成事务回滚 --------")
	test_11_4_2()
}

// 正常执行并完成事务提交
func test_11_4_1() {
	command := "./txtools -c \"chainmaker/config/conf11-4.toml\" -app \"atomic\" " +
		"-op \"send\" -vf 300 -tp 401 -chain1 10006"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))
}

// 状态异常并完成事务回滚
func test_11_4_2() {
	command := "./txtools -c \"chainmaker/config/conf11-4.toml\" -app \"atomic\" " +
		"-op \"rollback\" -vf 300 -tp 401 -chain1 10006"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))
}
