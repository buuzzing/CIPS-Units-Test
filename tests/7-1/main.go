package main

import (
	"bytes"
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"cipstests/chainmaker/common/types"
	"cipstests/chainmaker/common/utils"
	"encoding/binary"
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
	configFile = flag.String("c", "chainmaker/config/conf7-1.toml", "配置文件路径")

	seq = int64(rand.Intn(100000000))
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 7-1-1. 长安链relayer向中继链转发正确验证信息 --------")
	test_7_1_1()
	clog.Info("-------- 7-1-2. 长安链relayer向中继链转发被篡改的验证信息 --------")
	test_7_1_2(client)
}

// 长安链relayer向中继链转发正确验证信息
func test_7_1_1() {
	command := "./txtools -c \"chainmaker/config/conf7-1.toml\" -app \"atomic\" " +
		"-op \"send\" -vf 303 -tp 401"
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("执行命令失败: %v, 输出: %s", err, string(out)))
	}
	clog.Infof("命令输出: %s", string(out))

	// 等待一段时间，确保消息被处理
	time.Sleep(10 * time.Second)
}

// 长安链relayer向中继链转发被篡改的验证信息
func test_7_1_2(client *sdk.ChainClient) {

	ccMsg := types.CrosschainMessage{
		SrcChainId:          big.NewInt(123),
		DstChainId:          big.NewInt(456),
		Seq:                 big.NewInt(789),
		SrcAppId:            big.NewInt(111),
		DstAppId:            big.NewInt(222),
		PayloadReq:          [][]byte{[]byte("hello"), []byte("world")},
		PayloadResp:         [][]byte{[]byte("welcome"), []byte("home")},
		TransactionTypeId:   big.NewInt(333),
		TransactionPayload:  [][]byte{[]byte("transaction")},
		VerificationTypeId:  big.NewInt(444),
		VerificationPayload: [][]byte{[]byte("verification"), []byte("payload")},
		TransmissionTypeId:  big.NewInt(555),
		TransmissionPayload: [][]byte{[]byte("transmission")},
		TransportTypeId:     big.NewInt(666),
		TransportPayload:    [][]byte{[]byte("transport")},
		HashReq:             []byte{},
		HashResp:            []byte{},
		Ack:                 false,
	}
	//ccMsgHash := utils.CmToHash(ccMsg)
	ccMsgListBytes := utils.CmToLb(ccMsg)
	ccMsgListUint := utils.ConvertLbToUint(ccMsgListBytes)
	ccMsgBytes, _ := json.Marshal(ccMsgListUint)

	const hashLen = 32
	proof := new(bytes.Buffer)
	height := uint64(2)
	binary.Write(proof, binary.BigEndian, height)
	txHash := bytes.Repeat([]byte{0x11}, hashLen)
	proof.Write(txHash)
	node1 := bytes.Repeat([]byte{0xaa}, hashLen)
	proof.Write(node1)
	node2 := make([]byte, hashLen)
	proof.Write(node2)
	proofBytes := proof.Bytes()

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
		{Key: "proof", Value: proofBytes},
	}
	resp, err := chaintools.InvokeContract(client, types.VerificationAddr3, "verify", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
