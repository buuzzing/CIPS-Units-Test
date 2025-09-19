package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	chainmaker_sdk_go "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

const TestTime = 1800 // 测试总时长，单位秒

// 配置文件路径
var configFile *string

func init() {
	configFile = flag.String("c", "chainmaker/config/conf14-1.toml", "配置文件路径")
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	var crossChainRecord, singleChainRecord []time.Duration
	testCrossChain(client, &crossChainRecord)
	testSingleChain(client, &singleChainRecord)

	clog.Info("-------- 14-1-1 长安链-长安链间的跨链读时延 --------")
	fmt.Printf("跨链读时延测试统计结果:\n")
	printRecord(crossChainRecord)
	fmt.Printf("\n单链跨合约读时延测试统计结果:\n")
	printRecord(singleChainRecord)
}

func testCrossChain(client *chainmaker_sdk_go.ChainClient, record *[]time.Duration) {
	appArgs := [][]byte{
		[]byte("key1"),
		big.NewInt(301).Bytes(),
		big.NewInt(401).Bytes(),
	}
	appArgsBytes, _ := json.Marshal(&appArgs)

	kvs := []*common.KeyValuePair{
		{Key: "dstChainId", Value: big.NewInt(20007).Bytes()},
		{Key: "srcAppId", Value: big.NewInt(101).Bytes()},
		{Key: "dstAppId", Value: big.NewInt(101).Bytes()},
		{Key: "appArgs", Value: appArgsBytes},
	}

	ticker := time.NewTicker(TestTime * time.Second)
	defer ticker.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go checkEvent(ctx, client, done)

	start := time.Now()
	_, err := client.InvokeContract("caggr", "sendMsg", "", kvs, -1, true)
	if err != nil {
		panic(fmt.Sprintf("调用合约失败: %v", err))
	}

	for {
		select {
		case <-ticker.C:
			clog.Infof("跨链时延测试停止")
			cancel()
			return
		case <-done:
			cost := time.Since(start)
			*record = append(*record, cost)
			clog.Infof("跨链时延测试收到事件，时延: %dms", cost.Milliseconds())

			start = time.Now()
			_, err := client.InvokeContract("caggr", "sendMsg", "", kvs, -1, true)
			if err != nil {
				panic(fmt.Sprintf("调用合约失败: %v", err))
			}
		}
	}
}

func testSingleChain(client *chainmaker_sdk_go.ChainClient, record *[]time.Duration) {
	kvs := []*common.KeyValuePair{
		{Key: "key", Value: []byte("key1")},
	}

	ticker := time.NewTicker(TestTime * time.Second)
	defer ticker.Stop()
	tickerT := time.NewTicker(50 * time.Millisecond)
	defer tickerT.Stop()

	for {
		select {
		case <-ticker.C:
			clog.Infof("单链时延测试停止")
			return
		case <-tickerT.C:
			start := time.Now()
			_, err := client.InvokeContract("QueryApp", "getRecordPublic", "", kvs, -1, true)
			if err != nil {
				panic(fmt.Sprintf("调用合约失败: %v", err))
			}
			cost := time.Since(start)
			*record = append(*record, cost)
			clog.Infof("单链时延测试收到事件，时延: %dms", cost.Milliseconds())
		}
	}
}

func checkEvent(ctx context.Context, client *chainmaker_sdk_go.ChainClient, done chan<- struct{}) {
	ec, err := client.SubscribeContractEvent(ctx, -1, -1, "QueryApp", "queryResult")
	if err != nil {
		panic(fmt.Sprintf("订阅合约事件失败: %v", err))
	}

	for {
		select {
		case event, ok := <-ec:
			if !ok || event == nil {
				clog.Errorf("事件为空")
				return
			}
			contractEvent, ok := event.(*common.ContractEventInfo)
			if !ok {
				clog.Errorf("事件类型错误")
				return
			}
			if len(contractEvent.EventData) < 2 || string(contractEvent.EventData[0]) != "key1" {
				continue
			}
			clog.Debugf("收到合约事件: %v", contractEvent.EventData)
			done <- struct{}{}
		case <-ctx.Done():
			clog.Infof("context done")
			return
		}
	}
}

func printRecord(record []time.Duration) {
	fmt.Printf("总发送次数: %d\n", len(record))

	var total time.Duration
	var max time.Duration
	var min time.Duration
	for i, r := range record {
		total += r
		if i == 0 || r > max {
			max = r
		}
		if i == 0 || r < min {
			min = r
		}
	}
	avg := total / time.Duration(len(record))

	fmt.Printf("平均时延: %dms\n", avg.Milliseconds())
	fmt.Printf("最大时延: %dms\n", max.Milliseconds())
	fmt.Printf("最小时延: %dms\n", min.Milliseconds())
}
