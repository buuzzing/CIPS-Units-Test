package chaintools

import (
	"encoding/json"
	"fmt"
	"reflect"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

// 将交易响应转换为字符串
func respToString(resp *common.TxResponse) string {
	if resp == nil {
		return "nil transaction response"
	}

	info := "\nresp.Code->\n\t" + resp.Code.String()
	info += "\nresp.Message->\n\t" + resp.Message

	if resp.ContractResult != nil {
		info += "\nresp.ContractResult.Result->\n\t" + string(resp.ContractResult.Result) +
			"\nresp.ContractResult.Message->\n\t" + resp.ContractResult.Message
	}

	return info
}

// PrintTxResp 格式化输出长安链交易回执
// @param resp 交易回执
// @param resStruct 交易结果结构体，必须为指针，非结构体为 nil
func PrintTxResp(resp *common.TxResponse, resStruct interface{}) {
	const (
		txRespFormat = "Tx Basic Info:\n" +
			"\tTxId: %s\n" +
			"\tCode: %s\n" +
			"\tMessage: %s\n\n"

		txContractResultFormat = "Tx Contract Result:\n" +
			"\tCode: %d\n" +
			"\tResult: %+v\n" +
			"\tMessage: %s\n\n"

		txContractEventFormat = "Tx Contract Event #%d:\n" +
			"\tTopic: %s\n" +
			"\tTxId: %s\n" +
			"\tContractName: %s\n" +
			"\tEventData: %v\n\n"
	)

	if resp == nil {
		fmt.Println("Transaction response is nil")
		return
	}
	if resStruct != nil && reflect.ValueOf(resStruct).Kind() != reflect.Ptr {
		fmt.Println("Result structure must be a pointer")
		return
	}

	fmt.Printf(txRespFormat, resp.TxId, resp.Code.String(), resp.Message)

	if resStruct == nil {
		// 按字符串输出 ContractResult.Result
		fmt.Printf(txContractResultFormat,
			resp.ContractResult.Code, resp.ContractResult.Result, resp.ContractResult.Message)
	} else {
		// 将 ContractResult.Result 反序列化为 resStruct
		err := json.Unmarshal(resp.ContractResult.Result, resStruct)
		if err != nil {
			fmt.Printf(txContractResultFormat,
				resp.ContractResult.Code, err, resp.ContractResult.Message)
		} else {
			fmt.Printf(txContractResultFormat,
				resp.ContractResult.Code, resStruct, resp.ContractResult.Message)
		}
	}

	for index, event := range resp.ContractResult.ContractEvent {
		fmt.Printf(txContractEventFormat,
			index+1, event.Topic, event.TxId, event.ContractName, event.EventData)
	}
}
