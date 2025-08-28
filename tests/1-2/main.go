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
	configFile = flag.String("c", "chainmaker/config/conf1-2.toml", "配置文件路径")

}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 1-2-1. 根据协议地址查询已注册协议的ID --------")
	test_1_2_1(client)
	clog.Info("-------- 1-2-2. 根据协议地址查询未注册协议的ID --------")
	test_1_2_2(client)
	clog.Info("-------- 1-2-3. 根据协议ID查询已注册协议的地址 --------")
	test_1_2_3(client)
	clog.Info("-------- 1-2-4. 根据不存在的ID查询协议的地址 --------")
	test_1_2_4(client)
	clog.Info("-------- 1-2-5. 根据非法协议ID查询协议的地址 --------")
	test_1_2_5(client)

}

// 根据协议地址查询已注册协议的ID
func test_1_2_1(client *sdk.ChainClient) {
	protocol_address := []byte("ctpinf") //空传输协议
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "getTransportId", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 根据协议地址查询未注册协议的ID
func test_1_2_2(client *sdk.ChainClient) {
	protocol_address := []byte("unregistered") //未注册协议
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "getTransportId", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 根据协议ID查询已注册协议的地址
func test_1_2_3(client *sdk.ChainClient) {
	protocol_id := big.NewInt(400).Bytes() //空传输协议
	kvs := []*common.KeyValuePair{
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "get", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 根据不存在的ID查询协议的地址
func test_1_2_4(client *sdk.ChainClient) {
	protocol_id := big.NewInt(500).Bytes() //不存在的协议
	kvs := []*common.KeyValuePair{
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "get", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 根据非法协议ID查询协议的地址
func test_1_2_5(client *sdk.ChainClient) {
	protocol_id := big.NewInt(0).Bytes() //0或负数
	kvs := []*common.KeyValuePair{
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "get", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
