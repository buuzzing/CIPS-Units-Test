package chaintools

import (
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// 查询一个部署合约交易
func QueryTx(c *sdk.ChainClient, txId string) (info string, ok bool) {
	resp, err := c.GetTxByTxId(txId)
	if err != nil {
		if strings.Contains(err.Error(), "no such transaction") {
			// 交易不存在，可能是还未打包
			return "no such transaction", false
		} else {
			return err.Error(), false
		}
	}

	if resp.Transaction.Result.Code != common.TxStatusCode_SUCCESS {
		contractResult := resp.Transaction.Result.ContractResult
		if strings.Contains(contractResult.Message, "contract exist") {
			// 合约已经存在
			return fmt.Sprintf("contract exist @ Query Tx Block #%d", resp.BlockHeight), true
		} else {
			return fmt.Sprintf("%s @ Query Tx Block #%d", contractResult.Message, resp.BlockHeight), false
		}
	} else {
		// 交易存在且无错误，说明合约已经部署成功
		return fmt.Sprintf("contract deploy ok @ Query Tx Block #%d", resp.BlockHeight), true
	}
}
