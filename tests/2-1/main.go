package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"cipstests/chainmaker/common/types"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os/exec"
	"time"

	"math/rand"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string

// 随机指定一个 seq
var seq int64

func init() {
	configFile = flag.String("c", "chainmaker/config/conf2-1.toml", "配置文件路径")

	seq = int64(rand.Intn(100000000))
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 2-1-1. 长安链向中继链（长安链）发送跨链消息 --------")
	test_2_1_1()
	clog.Info("-------- 2-1-2. 长安链向中继链（长安链）发送错误编码的跨链消息 --------")
	test_2_1_2(client)
	clog.Info("-------- 2-1-3. 长安链向中继链（长安链）发送重复的跨链消息 --------")
	test_2_1_3(client)
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

// 长安链向中继链（长安链）发送错误编码的跨链消息
// 直接调用中继链验证层的接收方法
func test_2_1_2(client *sdk.ChainClient) {
	ccMsgBytes := []byte("error encoding")

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportAddr1, "receiveIn", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 长安链向中继链（长安链）发送重复的跨链消息
func test_2_1_3(client *sdk.ChainClient) {
	ccMsg := types.CrosschainMessage{
		SrcChainId:          big.NewInt(1),
		DstChainId:          big.NewInt(1),
		Seq:                 big.NewInt(seq),
		SrcAppId:            big.NewInt(1),
		DstAppId:            big.NewInt(1),
		PayloadReq:          [][]byte{},
		PayloadResp:         [][]byte{},
		TransactionTypeId:   big.NewInt(1),
		TransactionPayload:  [][]byte{},
		VerificationTypeId:  big.NewInt(1),
		VerificationPayload: [][]byte{},
		TransmissionTypeId:  big.NewInt(1),
		TransmissionPayload: [][]byte{},
		TransportTypeId:     big.NewInt(1),
		TransportPayload:    [][]byte{},
		HashReq:             []byte{},
		HashResp:            []byte{},
		Ack:                 false,
	}
	ccMsgBytes, _ := json.Marshal(ccMsg)

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
	}
	// 正常发送一次
	_, err := chaintools.InvokeContract(client, types.TransportAddr1, "sendOut", kvs, true)
	if err != nil {
		panic(err)
	}

	// 重复发送一次
	resp, err := chaintools.InvokeContract(client, types.TransportAddr1, "sendOut", kvs, true)
	if err != nil {
		panic(err)
	}
	// 预期这里会报错，提示重复消息
	chaintools.PrintTxResp(resp, nil)
}
