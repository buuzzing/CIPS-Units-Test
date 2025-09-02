package utils

import (
	"cipstests/chainmaker/common/types"
	"encoding/binary"
	"math/big"

	"crypto/sha256"
)

// CmToHash 计算跨链消息的哈希值
func CmToHash(ccMsg types.CrosschainMessage) []byte {
	elems := CmToLb(ccMsg)
	hasher := sha256.New()
	for _, elem := range elems {
		hasher.Write(elem)
	}
	result := hasher.Sum(nil)

	return result
}

// CmToLb 将跨链消息转换为字节数组
func CmToLb(ccMsg types.CrosschainMessage) [][]byte {
	var ret [][]byte
	var i int
	// 转发层
	ret = append(ret, BigIntToBytes32(ccMsg.TransportTypeId))
	ret = append(ret, IntToBytes32(len(ccMsg.TransportPayload)))
	for i = 0; i < len(ccMsg.TransportPayload); i += 1 {
		ret = append(ret, ccMsg.TransportPayload[i])
	}
	// 验证层
	ret = append(ret, BigIntToBytes32(ccMsg.VerificationTypeId))
	ret = append(ret, IntToBytes32(len(ccMsg.VerificationPayload)))
	for i = 0; i < len(ccMsg.VerificationPayload); i += 1 {
		ret = append(ret, ccMsg.VerificationPayload[i])
	}
	// 传输层
	ret = append(ret, BigIntToBytes32(ccMsg.TransmissionTypeId))
	ret = append(ret, IntToBytes32(len(ccMsg.TransmissionPayload)))
	for i = 0; i < len(ccMsg.TransmissionPayload); i += 1 {
		ret = append(ret, ccMsg.TransmissionPayload[i])
	}
	// 事务层
	ret = append(ret, BigIntToBytes32(ccMsg.TransactionTypeId))
	ret = append(ret, IntToBytes32(len(ccMsg.TransactionPayload)))
	for i = 0; i < len(ccMsg.TransactionPayload); i += 1 {
		ret = append(ret, ccMsg.TransactionPayload[i])
	}
	// 应用层
	ret = append(ret, BigIntToBytes32(ccMsg.SrcAppId))
	ret = append(ret, BigIntToBytes32(ccMsg.DstAppId))
	ret = append(ret, IntToBytes32(len(ccMsg.PayloadReq)))
	for i = 0; i < len(ccMsg.PayloadReq); i += 1 {
		ret = append(ret, ccMsg.PayloadReq[i])
	}
	ret = append(ret, IntToBytes32(len(ccMsg.PayloadResp)))
	for i = 0; i < len(ccMsg.PayloadResp); i += 1 {
		ret = append(ret, ccMsg.PayloadResp[i])
	}
	// 公共参数
	ret = append(ret, BigIntToBytes32(ccMsg.SrcChainId))
	ret = append(ret, BigIntToBytes32(ccMsg.DstChainId))
	ret = append(ret, BigIntToBytes32(ccMsg.Seq))

	ret = append(ret, BoolToByte(ccMsg.Ack))

	return ret

}

// CmFromLb 将字节数组转换为跨链消息
func CmFromLb(data [][]byte) types.CrosschainMessage {
	var ccMsg types.CrosschainMessage
	var endIndex int
	convertLb := func(start int) ([][]byte, int) {
		length := Bytes32ToInt(data[start])
		end := start + length
		var result [][]byte
		for i := start + 1; i < end+1; i += 1 {
			result = append(result, data[i])
		}
		return result, end
	}
	//转发层
	ccMsg.TransportTypeId = Bytes32ToBigInt(data[0])
	ccMsg.TransportPayload, endIndex = convertLb(1)
	//验证层
	ccMsg.VerificationTypeId = Bytes32ToBigInt(data[endIndex+1])
	ccMsg.VerificationPayload, endIndex = convertLb(endIndex + 2)
	//传输层
	ccMsg.TransmissionTypeId = Bytes32ToBigInt(data[endIndex+1])
	ccMsg.TransmissionPayload, endIndex = convertLb(endIndex + 2)
	//事务层
	ccMsg.TransactionTypeId = Bytes32ToBigInt(data[endIndex+1])
	ccMsg.TransactionPayload, endIndex = convertLb(endIndex + 2)
	//应用层
	ccMsg.SrcAppId = Bytes32ToBigInt(data[endIndex+1])
	ccMsg.DstAppId = Bytes32ToBigInt(data[endIndex+2])
	ccMsg.PayloadReq, endIndex = convertLb(endIndex + 3)
	ccMsg.PayloadResp, endIndex = convertLb(endIndex + 1)
	//公共
	ccMsg.SrcChainId = Bytes32ToBigInt(data[endIndex+1])
	ccMsg.DstChainId = Bytes32ToBigInt(data[endIndex+2])
	ccMsg.Seq = Bytes32ToBigInt(data[endIndex+3])
	ccMsg.Ack = ByteToBool(data[endIndex+4])

	return ccMsg
}

// Bytes32ToBigInt 将字节数组转换为 big.Int
func Bytes32ToBigInt(dataBytes []byte) *big.Int {
	return new(big.Int).SetBytes(dataBytes)
}

// BigIntToBytes32 将 big.Int 转换为字节数组
func BigIntToBytes32(dataBigInt *big.Int) []byte {
	bigBytes := dataBigInt.Bytes()
	var bigBytes32 [32]byte
	copy(bigBytes32[32-len(bigBytes):], bigBytes)
	return bigBytes32[:]
}

// ByteToBool 将单字节转换为布尔值
func ByteToBool(dataBytes []byte) bool {
	dataByte := dataBytes[0]
	return dataByte != 0
}

// BoolToByte 将布尔值转换为单字节
func BoolToByte(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// Bytes32ToInt 将字节数组转换为整数
func Bytes32ToInt(dataBytes []byte) int {
	// 直接从最后8字节提取
	dataInt := int(binary.BigEndian.Uint64(dataBytes[24:]))
	return dataInt
}

// IntToBytes32 将整数转换为字节数组
func IntToBytes32(dataInt int) []byte {
	dataUint64 := uint64(dataInt)
	var dataBytes [32]byte
	// 将 uint64 值转换为大端序，并填充到 [32]byte 数组的末尾8个字节
	binary.BigEndian.PutUint64(dataBytes[24:], dataUint64)
	return dataBytes[:]
}

func ConvertLbToUint(data [][]byte) [][]uint {
	// 将 [][]byte 转换为 [][]uint
	// 这里假设 data 中的每个 []byte 都是一个 uint 的二进制表示
	result := make([][]uint, len(data))
	for i, b := range data {
		result[i] = make([]uint, len(b))
		for j, v := range b {
			result[i][j] = uint(v)
		}
	}
	return result
}
