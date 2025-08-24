package chaintools

import (
	xconf "cipstests/chainmaker/common/config"
	"fmt"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// NewChainClient 创建一个新的长安链客户端实例
func NewChainClient(nodeAddr string, nodeConnNum int, hashType string,
	authType string, chainId string, userKeyPath string) (*sdk.ChainClient, error) {
	// 节点配置（节点地址和连接数）
	nodeConfig := sdk.NewNodeConfig(
		sdk.WithNodeAddr(nodeAddr),
		sdk.WithNodeConnCnt(nodeConnNum),
	)

	// 加密配置（公钥模式下可指定哈希算法）
	cryptoConfig := sdk.NewCryptoConfig(
		sdk.WithHashAlgo(hashType),
	)

	client, err := sdk.NewChainClient(
		// 添加节点配置
		sdk.AddChainClientNodeConfig(nodeConfig),
		// 添加哈希算法配置
		sdk.WithCryptoConfig(cryptoConfig),
		// 认证模式
		sdk.WithAuthType(authType),
		// 链 ID
		sdk.WithChainClientChainId(chainId),
		// 用户密钥路径
		sdk.WithUserSignKeyFilePath(userKeyPath),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetChainClient 获取一个默认配置的长安链客户端实例
func GetChainClient() *sdk.ChainClient {
	// 获取链配置
	chainConf := xconf.GetChainConfig()
	if chainConf == nil {
		panic("获取链配置失败")
	}

	// 创建链客户端
	client, err := NewChainClient(
		chainConf.RpcAddress,
		chainConf.ConnNum,
		chainConf.HashType,
		chainConf.AuthType,
		chainConf.ChainId,
		chainConf.UserKeyPath[0],
	)
	if err != nil {
		panic(fmt.Sprintf("创建链客户端失败: %v", err))
	}

	return client
}
