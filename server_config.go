package thorn

type ServerConfig struct {
	// http地址
	ManageAddr  string

	// 客户端连接地址
	ControlAddr string
}

func NewServerConfig(mAddr string, cAddr string) ServerConfig {
	return ServerConfig{
		ManageAddr:  mAddr,
		ControlAddr: cAddr,
	}
}