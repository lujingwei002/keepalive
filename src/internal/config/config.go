package config

import (
	"fmt"

	"github.com/lujingwei002/keepalive"
	"github.com/lujingwei002/keepalive/config"
	"github.com/lujingwei002/keepalive/internal/server"
)

// 根据配置初始化应用程序
func InitApplication(app keepalive.Application, config *config.Config) error {
	for name, c := range config.Servers {
		if build, ok := server.Get(c.Type); !ok {
			return fmt.Errorf("server type %s not support", c.Type)
		} else if s := build(); s == nil {
			return fmt.Errorf("create server %s fail", c.Type)
		} else if err := s.Init(app, &c); err != nil {
			return fmt.Errorf("server %s init fail: %s", c.Type, err)
		} else if err := app.AddServer(name, s); err != nil {
			return err
		}
		fmt.Println(name)
	}
	return nil
}
