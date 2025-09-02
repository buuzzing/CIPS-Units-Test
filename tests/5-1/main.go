package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"cipstests/chainmaker/common/types"
	"cipstests/chainmaker/common/utils"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string

// 随机指定一个 seq
var seq int64

func init() {
	configFile = flag.String("c", "chainmaker/config/conf5-1.toml", "配置文件路径")

	seq = int64(rand.Intn(100000000))
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 5-1-1. 长安链relayer向中继链转发正确验证信息 --------")
	test_5_1_1()
	clog.Info("-------- 5-1-2. 长安链relayer向中继链转发错误验证信息 --------")
	test_5_1_2(client)
}

// 长安链relayer向中继链转发正确验证信息
func test_5_1_1() {
	command := "go run ~/goproject/CIPS-Gemini-ChainMaker/scripts/tx/* -c \"chainmaker/config/conf5-1.toml\" -app \"sendMsg\" " +
		"-op \"send\" -vf 301 -tp 401 -chain1 20007"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}

// 长安链relayer向中继链转发错误验证信息
func test_5_1_2(client *sdk.ChainClient) {

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
	ccMsgListBytes := utils.CmToLb(ccMsg)
	ccMsgListUint := utils.ConvertLbToUint(ccMsgListBytes)
	ccMsgBytes, _ := json.Marshal(ccMsgListUint)
	ccMsgHash := utils.CmToHash(ccMsg)
	proof := fmt.Sprintf("%x", ccMsgHash) +
		"1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c212c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b" +
		"020b28e6b33555642dabdf2855d80b0955918bb52a7e5ef3159360c89a3239f264"
	proofBytes, _ := hex.DecodeString(proof)

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
		{Key: "proof", Value: proofBytes},
	}
	resp, err := chaintools.InvokeContract(client, types.VerificationAddr1, "verify", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
