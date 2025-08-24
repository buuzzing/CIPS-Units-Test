package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// configer 全局配置对象
var configer xConfiger

// xConfiger 配置管理对象
type xConfiger struct {
	MainConf  *XConfig
	ChainConf *ChainConfig
}

// Init 配置初始化函数
func Init(file string) error {
	configer = xConfiger{}

	configer.MainConf = &XConfig{}
	if err := loadAndParseConfig(file, configer.MainConf); err != nil {
		return fmt.Errorf("failed to load main config: %w", err)
	}

	chainConfig := &ChainConfig{}
	if err := loadAndParseConfig(configer.MainConf.Chain, chainConfig); err != nil {
		return fmt.Errorf("failed to load chain config: %w", err)
	}
	configer.ChainConf = chainConfig

	return nil
}

// loadAndParseConfig 加载及解析配置文件
func loadAndParseConfig(file string, obj interface{}) error {
	if _, err := toml.DecodeFile(file, obj); err != nil {
		return fmt.Errorf("failed to decode config file %s: %w", file, err)
	}

	return nil
}

// GetConfiger 获取全局配置对象
func GetConfiger() *xConfiger {
	return &configer
}

// GetChainConfig 获取链配置
func GetChainConfig() *ChainConfig {
	return configer.ChainConf
}

// GetDeployConfig 获取部署配置
func GetDeployConfig() *DeployConfig {
	return configer.MainConf.Deploy
}

// GetRelayerConfig 获取 relayer 配置
func GetRelayerConfig() *RelayerConfig {
	return configer.MainConf.Relayer
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	return configer.MainConf.Log
}

// GetTransportAddrById 根据传输层 ID 获取传输层地址
func GetTransportAddrById(id int) string {
	for _, p := range configer.ChainConf.Protocols {
		if p.Layer == LayerTransport && p.Type == TypeProtocol && p.ProtocolId == id {
			return p.Address
		}
	}
	return ""
}

// GetAggregatorAddr 获取聚合器地址
func GetAggregatorAddr() string {
	// 目前一条链只有一个聚合器
	if len(configer.ChainConf.Aggregatores) != 1 {
		return ""
	}
	return configer.ChainConf.Aggregatores[0].Address
}

// GetRegAddrByLayer 根据层级获取注册器地址
func GetRegAddrByLayer(layer string) string {
	for _, r := range configer.ChainConf.Registers {
		if r.Layer == layer && r.Type == TypeRegister {
			return r.Address
		}
	}
	return ""
}
