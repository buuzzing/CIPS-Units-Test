package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"cipstests/chainmaker/common/types"
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
	configFile = flag.String("c", "chainmaker/config/conf1-1.toml", "配置文件路径")

}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("--------1-1-1. 正常注册新部署的传输协议 --------")
	test_1_1_1(client)
	clog.Info("-------- 1-1-2. 重复注册已部署的传输协议 --------")
	test_1_1_2(client)

}

// 正常注册新部署的传输协议
func test_1_1_1(client *sdk.ChainClient) {
	protocol_address := []byte("NewAddress9")
	protocol_id := big.NewInt(105).Bytes()
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "set", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 重复注册已部署的传输协议
func test_1_1_2(client *sdk.ChainClient) {
	protocol_address := []byte("NewAddress9") //与1-1-1中的协议地址一致
	protocol_id := big.NewInt(105).Bytes()
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "set", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
