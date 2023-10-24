package keepalive

import "github.com/lujingwei002/keepalive/config"

// interface

// message direction
// server: ws|tcp => transport => |middleware| => backend.sink
// backend.source => |middleware| => transport => server: ws|tcp

type Application interface {
	Run() error

	Go(f func() error)

	AddServer(name string, s Server) error

	HandleClientMessage(message []byte) error
}

type Server interface {
	// 初始化服务器
	Init(app Application, config *config.Server) error
	// 启动服务器
	Start(app Application) error
}

type Backend interface {
}

type Source interface {
}

type Sink interface {
}

type Transport interface {
}
