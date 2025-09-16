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
	configFile = flag.String("c", "chainmaker/config/conf1-3.toml", "配置文件路径")
	protocolId = time.Now().Unix()
	protocolAddr = "NewAddress" + fmt.Sprintf("%d", protocolId)
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	clog.Info("-------- 1-3-1. 更新已注册协议 --------")
	test_1_3_1(client)
	clog.Info("-------- 1-3-2. 更新未注册协议 --------")
	test_1_3_2(client)

}

// 更新已注册协议
func test_1_3_1(client *sdk.ChainClient) {
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
	//更新协议
	protocol_address = []byte("ctpinf")   //空传输协议
	kvs = []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
		{Key: "id", Value: protocol_id},
	}
	resp, err = chaintools.InvokeContract(client, types.TransportRegAddr, "update", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}

// 更新未注册协议
func test_1_3_2(client *sdk.ChainClient) {
	protocol_address := []byte("ctpinf")
	protocol_id := big.NewInt(999).Bytes() //未注册id
	kvs := []*common.KeyValuePair{
		{Key: "address", Value: protocol_address},
		{Key: "id", Value: protocol_id},
	}
	resp, err := chaintools.InvokeContract(client, types.TransportRegAddr, "update", kvs, true)
	if err != nil {
		panic(err)
	}

	chaintools.PrintTxResp(resp, nil)
}
