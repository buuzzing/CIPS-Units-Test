package main

import (
	"cipstests/chainmaker/common/chaintools"
	xconf "cipstests/chainmaker/common/config"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"slices"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	chainmaker_sdk_go "chainmaker.org/chainmaker/sdk-go/v2"
	clog "github.com/kpango/glg"
)

var TestTime int // 测试总时长，单位秒

// 配置文件路径
var configFile *string

func init() {
	configFile = flag.String("c", "chainmaker/config/conf14-1.toml", "配置文件路径")
	flag.IntVar(&TestTime, "t", 60, "测试总时长，单位秒")
}

func main() {
	flag.Parse()
	if err := xconf.Init(*configFile); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	client := chaintools.GetChainClient()

	var crossChainRecord, crossChainRecord2 []time.Duration
	testCrossChain(client, &crossChainRecord)
	testCrossChainWithoutStack(client, &crossChainRecord2)

	clog.Info("-------- 14-1 长安链-长安链间的跨链读时延 --------")
	fmt.Printf("跨链读时延（完整协议栈）测试统计结果:\n")
	printRecord(crossChainRecord)
	fmt.Printf("\n跨链读时延（未部署协议栈）测试统计结果:\n")
	printRecord(crossChainRecord2)
}

func testCrossChain(client *chainmaker_sdk_go.ChainClient, record *[]time.Duration) {
	appArgs := [][]byte{
		{0},            // OpType
		[]byte(""),     // OriginKey
		[]byte(""),     // OriginVal
		[]byte("key1"), // TargetKey
		[]byte(""),     // TargetVal
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

	ticker := time.NewTicker(time.Duration(TestTime) * time.Second)
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
			clog.Infof("跨链读（完整协议栈）时延测试停止")
			cancel()
			return
		case <-done:
			cost := time.Since(start)
			*record = append(*record, cost)
			clog.Infof("跨链读（完整协议栈）时延测试收到事件，时延: %dms", cost.Milliseconds())

			start = time.Now()
			_, err := client.InvokeContract("caggr", "sendMsg", "", kvs, -1, true)
			if err != nil {
				panic(fmt.Sprintf("调用合约失败: %v", err))
			}
		}
	}
}

func testCrossChainWithoutStack(client *chainmaker_sdk_go.ChainClient, record *[]time.Duration) {
	appArgs := [][]byte{
		{0},            // OpType
		[]byte(""),     // OriginKey
		[]byte(""),     // OriginVal
		[]byte("key1"), // TargetKey
		[]byte(""),     // TargetVal
		big.NewInt(300).Bytes(),
		big.NewInt(401).Bytes(),
	}
	appArgsBytes, _ := json.Marshal(&appArgs)

	kvs := []*common.KeyValuePair{
		{Key: "dstChainId", Value: big.NewInt(20007).Bytes()},
		{Key: "srcAppId", Value: big.NewInt(101).Bytes()},
		{Key: "dstAppId", Value: big.NewInt(101).Bytes()},
		{Key: "appArgs", Value: appArgsBytes},
	}

	ticker := time.NewTicker(time.Duration(TestTime) * time.Second)
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
			clog.Infof("跨链读（未部署协议栈）时延测试停止")
			cancel()
			return
		case <-done:
			cost := time.Since(start)
			*record = append(*record, cost)
			clog.Infof("跨链读（未部署协议栈）时延测试收到事件，时延: %dms", cost.Milliseconds())

			start = time.Now()
			_, err := client.InvokeContract("caggr", "sendMsg", "", kvs, -1, true)
			if err != nil {
				panic(fmt.Sprintf("调用合约失败: %v", err))
			}
		}
	}
}

func checkEvent(ctx context.Context, client *chainmaker_sdk_go.ChainClient, done chan<- struct{}) {
	ec, err := client.SubscribeContractEvent(ctx, -1, -1, "QueryApp", "tx_finalize")
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
			if len(contractEvent.EventData) < 2 ||
				contractEvent.EventData[1] != "OpTypeRead" {
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

	slices.Sort(record)

	var total time.Duration
	for _, r := range record {
		total += r
	}
	max := record[len(record)-1]
	min := record[0]
	avg := total / time.Duration(len(record))
	p90 := record[int(float64(len(record))*0.9)-1]
	p99 := record[int(float64(len(record))*0.99)-1]

	fmt.Printf("最小时延: %.3fs\n", float64(min.Milliseconds())/1000.0)
	fmt.Printf("平均时延: %.3fs\n", float64(avg.Milliseconds())/1000.0)
	fmt.Printf("P90 时延: %.3fs\n", float64(p90.Milliseconds())/1000.0)
	fmt.Printf("P99 时延: %.3fs\n", float64(p99.Milliseconds())/1000.0)
	fmt.Printf("最大时延: %.3fs\n", float64(max.Milliseconds())/1000.0)
}
