package config

// XConfig 定义 conf.toml 配置
type XConfig struct {
	Chain   string         `toml:"chain"`   // 启动的链配置文件列表
	Deploy  *DeployConfig  `toml:"deploy"`  // 协议合约部署配置
	Relayer *RelayerConfig `toml:"relayer"` // relayer 配置信息列表
	Log     *LogConfig     `toml:"log"`     // 日志配置
}

// DeployConfig 协议合约部署配置
type DeployConfig struct {
	Name            string `toml:"name"`              // 部署任务名称
	Enable          bool   `toml:"enable"`            // 是否启用该部署任务
	DeployTimeout   int    `toml:"deploy_timeout"`    // 部署合约超时时间（秒）
	DeployBatchSize int    `toml:"deploy_batch_size"` // 异步部署合约批量大小
}

// RelayerConfig relayer 配置信息
type RelayerConfig struct {
	Name             string              `toml:"name"`              // relayer 名称
	Enable           bool                `toml:"enable"`            // 是否启用 relayer
	ChainId          int                 `toml:"chain_id"`          // relayer 关联的长安链在协议栈中的 ID
	TransportIds     []int               `toml:"transport_ids"`     // relayer 关联的传输层 ID 列表，将监听这些 ID 的传输层事件
	TransportAddr    string              `toml:"transport_addr"`    // relayer 的传输层服务地址
	VerificationAddr string              `toml:"verification_addr"` // relayer 的验证层服务地址
	PluginServers    []PluginServerInfo  `toml:"plugin_servers"`    // relayer 关联的验证插件列表
	GatewayServers   []GatewayServerInfo `toml:"gateway_servers"`   // relayer 关联的其他网关列表
}

// PluginServerInfo 插件服务器信息
type PluginServerInfo struct {
	VerificationId int    `toml:"verification_id"` // plugin server 对应的验证层 ID
	PluginAddr     string `toml:"plugin_addr"`     // plugin server 服务地址
}

// GatewayServerInfo 网关服务器信息
type GatewayServerInfo struct {
	ChainId          int    `toml:"chain_id"`          // 该网关对应的链 ID
	VerificationAddr string `toml:"verification_addr"` // 该网关的验证层服务地址
	TransportAddr    string `toml:"transport_addr"`    // 该网关的传输层服务地址
}

// ChainConfig 长安链配置信息
type ChainConfig struct {
	ChainId      string          `toml:"chain_id"`      // 启动的链 ID
	AuthType     string          `toml:"auth_type"`     // 链认证类型
	HashType     string          `toml:"hash_type"`     // 哈希算法
	RpcAddress   string          `toml:"rpc_addr"`      // RPC 服务地址
	ConnNum      int             `toml:"conn_num"`      // 链节点连接数
	UserKeyPath  []string        `toml:"user_key_path"` // 用户密钥路径
	Registers    []*ContractInfo `toml:"registers"`     // 注册器信息
	Aggregatores []*ContractInfo `toml:"aggregators"`   // 聚合器信息
	Protocols    []*ContractInfo `toml:"protocols"`     // 协议信息
}

// ContractInfo 合约配置信息
type ContractInfo struct {
	Name       string                 `toml:"name"`        // 协议名称
	Address    string                 `toml:"addr"`        // 协议地址
	Type       string                 `toml:"type"`        // 协议类型: register, protocol, aggregator
	Layer      string                 `toml:"layer"`       // 协议所在层: app, transaction, transmission, verification, transport
	Runtime    string                 `toml:"runtime"`     // 协议运行时: go, wasm
	FilePath   string                 `toml:"file_path"`   // 协议文件路径
	ProtocolId int                    `toml:"protocol_id"` // 协议 ID
	Args       map[string]interface{} `toml:"args"`        // 其他参数
}

const (
	TypeRegister   = "register"   // 合约为注册器
	TypeProtocol   = "protocol"   // 合约为协议
	TypeAggregator = "aggregator" // 合约为聚合器
)

const (
	LayerApp          = "app"          // 应用层
	LayerTransaction  = "transaction"  // 事务层
	LayerTransmission = "transmission" // transmission 层
	LayerVerification = "verification" // 验证层
	LayerTransport    = "transport"    // 传输层
)

const (
	RuntimeGo   = "go"   // Go 运行时
	RuntimeWasm = "wasm" // WebAssembly 运行时
)

// LogConfig 日志配置
type LogConfig struct {
	Level     string        `toml:"level"` // 日志级别
	DebugConf *LogLevelConf `toml:"debug"`
	InfoConf  *LogLevelConf `toml:"info"`
	WarnConf  *LogLevelConf `toml:"warn"`
}

// LogLevelConf 具体日志级别配置
type LogLevelConf struct {
	Writer string `toml:"writer"` // 日志输出方式 file, console
	File   string `toml:"file"`   // 日志文件路径
}
