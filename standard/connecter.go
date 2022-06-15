package standard

// IConnect 连接标准
type IConnect interface {
	// Connect 连接
	Connect(config *Config, sos ...ServerOption) (IBuilder, error)
}
