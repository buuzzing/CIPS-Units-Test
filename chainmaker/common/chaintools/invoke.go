package chaintools

import (
	clog "github.com/kpango/glg"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

func InvokeContract(
	c *sdk.ChainClient,
	contractName, method string,
	kvs []*common.KeyValuePair,
	withSync bool,
) (*common.TxResponse, error) {
	resp, err := c.InvokeContract(
		contractName, method, "", kvs, -1, withSync)
	if err != nil || resp.Code != common.TxStatusCode_SUCCESS {
		_ = clog.Errorf("InvokeContract failed: %v, Name: %s, Method: %s, resp: %s",
			err, contractName, method, respToString(resp))
		return nil, err
	}

	return resp, nil
}
