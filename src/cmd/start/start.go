package start

import (
	"context"
	"fmt"

	"github.com/lujingwei002/keepalive"
	"github.com/lujingwei002/keepalive/config"
	internal_config "github.com/lujingwei002/keepalive/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	cfgFile string
)

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		cancelCtx, cancelFunc := context.WithCancel(ctx)
		errGroup, errCtx := errgroup.WithContext(cancelCtx)
		app := &Application{
			servers:    make(map[string]keepalive.Server),
			ctx:        cancelCtx,
			cancelFunc: cancelFunc,
			errGroup:   errGroup,
			errCtx:     errCtx,
		}
		return app.Run()
	},
}

func init() {
	Cmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.cobra.yaml)")
	Cmd.MarkPersistentFlagRequired("config")
}

type Application struct {
	servers    map[string]keepalive.Server
	ctx        context.Context
	cancelFunc context.CancelFunc
	errGroup   *errgroup.Group
	errCtx     context.Context
}

// implement keepalive.Application
func (app *Application) Run() error {
	if err := config.Parse(cfgFile); err != nil {
		return err
	}
	c := config.Get()
	if err := internal_config.InitApplication(app, c); err != nil {
		return err
	}
	for _, s := range app.servers {
		if err := s.Start(app); err != nil {
			return err
		}
	}
	return app.errGroup.Wait()

}

func (app *Application) AddServer(name string, server keepalive.Server) error {
	if _, ok := app.servers[name]; ok {
		return fmt.Errorf("server %s already exist", name)
	}
	app.servers[name] = server
	return nil
}

func (app *Application) Go(f func() error) {
	app.errGroup.Go(f)
}

func (app *Application) HandleClientMessage(message []byte) error {
	return nil
}
