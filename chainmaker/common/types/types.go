package types

import "math/big"

// CrosschainMessage 跨链消息结构体
type CrosschainMessage struct {
	//公用
	SrcChainId *big.Int
	DstChainId *big.Int
	Seq        *big.Int
	//应用层
	SrcAppId    *big.Int
	DstAppId    *big.Int
	PayloadReq  [][]byte
	PayloadResp [][]byte
	//事务层
	TransactionTypeId  *big.Int
	TransactionPayload [][]byte
	//验证层
	VerificationTypeId  *big.Int
	VerificationPayload [][]byte
	//传输层
	TransmissionTypeId  *big.Int
	TransmissionPayload [][]byte
	//转发层
	TransportTypeId  *big.Int
	TransportPayload [][]byte
	HashReq          []byte
	HashResp         []byte
	Ack              bool
}

// CmWithProof 跨链消息和证明结构体
type CmWithProof struct {
	Cm    *CrosschainMessage
	Proof []byte
}

// TrustRootInfo 信任根信息，用于 relayer 的 verification 监听并向目的链发送
type TrustRootInfo struct {
	// UpdateId 验证插件向 relayer 验证层发送 update 时生成
	UpdateId string
	// VerificationId 验证层 Id
	VerificationId *big.Int
	// TrustRoot 待更新的 TrustRoot
	TrustRoot [][]byte
}

// UpdateTxInfo 触发链上 update 操作的交易哈希
type UpdateTxInfo struct {
	UpdateId string
	// Update 交易对应的 TxId
	TxId string
}

// UpdateTxStatus 触发链上 update 操作的交易结果
type UpdateTxStatus struct {
	TxId string
	// 0 -> 未初始化，1 -> 成功，2 -> 失败，3 -> 查询超时
	Status uint8
}
