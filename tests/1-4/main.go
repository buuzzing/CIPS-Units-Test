package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"cipstests/chainmaker/common/types"
	"flag"
	"fmt"
	"math/big"
	"time"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

// 配置文件路径
var configFile *string
var protocolId int64
var protocolAddr string
func init() {
	configFile = flag.String("c", "chainmaker/config/conf1-4.toml", "配置文件路径")
	protocolId = time.Now().Unix()
	protocolAddr = "NewAddress" + fmt.Sprintf("%d", protocolId)
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 1-4-1. 退出已注册协议 --------")
	test_1_4_1(client)
	clog.Info("-------- 1-4-2. 退出未注册协议 --------")
	test_1_4_2(client)

}

// 退出已注册协议
func test_1_4_1(client *sdk.ChainClient) {
	//注册协议
	protocol_address := []byte(protocolAddr)
	protocol_id := big.NewInt(protocolId).Bytes()
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "set", kvs, true)
	if err != nil {
		panic(err)
	}
	//退出协议
	kvs = []*common.KeyValuePair{
		{Key: "id", Value: protocol_id},
	}
	resp, err = chaintools.InvokeContract(client, types.TransportRegAddr, "leave", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)

	resp, err = chaintools.InvokeContract(client, types.TransportRegAddr, "get", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 退出未注册协议
func test_1_4_2(client *sdk.ChainClient) {
	protocol_id := big.NewInt(999).Bytes()
	kvs := []*common.KeyValuePair{
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "leave", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
