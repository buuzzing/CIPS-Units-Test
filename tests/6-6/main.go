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

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string

func init() {
	configFile = flag.String("c", "chainmaker/config/conf6-6.toml", "配置文件路径")
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 6-6-1. 长安链relayer向中继链转发仅包含部分签名的验证信息 --------")
	test_6_6_1(client)
	clog.Info("-------- 6-6-2. 长安链relayer向中继链转发被篡改的验证信息 --------")
	test_6_6_2(client)
}

// 长安链relayer向中继链转发仅包含部分签名的验证信息
func test_6_6_1(client *sdk.ChainClient) {
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
	ccMsgListBytes := utils.CmToLb(ccMsg)
	ccMsgListUint := utils.ConvertLbToUint(ccMsgListBytes)
	ccMsgBytes, _ := json.Marshal(ccMsgListUint)
	ccMsgHash := utils.CmToHash(ccMsg)
	proof := fmt.Sprintf("%x", ccMsgHash) +
		"1687f985433b446b85eb6d0a574fc152f681c032d27e6207569faca9c8329b961b4b60273ae700a7e2ffc04e19e316074a5977c8da56b75675927e2eee23772e24fb6baf4cf6d7ca7eaa668cda36d088502b3587667b6eb8f2b874622575e5861e7cf2fd8b4bc0d81e4719f009a5ecb7d925c970bc57889f3627d86629dc31d8" +
		"0221ed00f5e8ae81dbb296bb70900741f29915b82f917f496fd70aa0d3782dd9df"
	proofBytes, _ := hex.DecodeString(proof)

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
		{Key: "proof", Value: proofBytes},
	}
	resp, err := chaintools.InvokeContract(client, types.VerificationAddr2, "verify", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 长安链relayer向中继链转发被篡改的验证信息
func test_6_6_2(client *sdk.ChainClient) {
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
	ccMsgHash := utils.CmToHash(ccMsg)
	ccMsgError := types.CrosschainMessage{
		SrcChainId:          big.NewInt(999),
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
	ccMsgListBytes := utils.CmToLb(ccMsgError)
	ccMsgListUint := utils.ConvertLbToUint(ccMsgListBytes)
	ccMsgBytes, _ := json.Marshal(ccMsgListUint)
	proof := fmt.Sprintf("%x", ccMsgHash) +
		"1687f985433b446b85eb6d0a574fc152f681c032d27e6207569faca9c8329b961b4b60273ae700a7e2ffc04e19e316074a5977c8da56b75675927e2eee23772e24fb6baf4cf6d7ca7eaa668cda36d088502b3587667b6eb8f2b874622575e5861e7cf2fd8b4bc0d81e4719f009a5ecb7d925c970bc57889f3627d86629dc31d8" +
		"020b28e6b33555642dabdf2855d80b0955918bb52a7e5ef3159360c89a3239f264"
	proofBytes, _ := hex.DecodeString(proof)

	kvs := []*common.KeyValuePair{
		{Key: "ccMsg", Value: ccMsgBytes},
		{Key: "proof", Value: proofBytes},
	}
	resp, err := chaintools.InvokeContract(client, types.VerificationAddr2, "verify", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
