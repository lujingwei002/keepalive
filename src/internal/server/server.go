package server

import (
	"github.com/lujingwei002/keepalive"
	"github.com/lujingwei002/keepalive/internal/server/ws"
)

var (
	servers map[string]func() keepalive.Server
)

func init() {
	servers = make(map[string]func() keepalive.Server)
	Register("ws", ws.New)
}

func Register(name string, f func() keepalive.Server) {
	servers[name] = f
}

func Get(name string) (f func() keepalive.Server, ok bool) {
	f, ok = servers[name]
	return
}
