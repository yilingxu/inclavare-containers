package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alibaba/inclavare-containers/shim/runtime/signature/server/conf"

	"github.com/alibaba/inclavare-containers/shim/runtime/signature/server/api"
	"github.com/golang/glog"
)

type Server struct {
	config    *conf.Config
	sigCh     chan os.Signal
	apiServer *api.ApiServer
}

func NewServer(conf *conf.Config) (*Server, error) {
	apiSvr, err := api.NewApiServer(":9080", conf)
	if err != nil {
		glog.Errorf("Failed to create ApiServer.err:%s", err.Error())
	}
	svr := &Server{
		config:    conf,
		sigCh:     make(chan os.Signal, 1),
		apiServer: apiSvr,
	}
	// signal trap
	signal.Notify(svr.sigCh, syscall.SIGINT)
	return svr, nil
}

func (svr *Server) Start(stopChan <-chan struct{}) {

	glog.Info("Starting HttpServer ...")
	go func() {
		if err := svr.apiServer.RunForeground(); err != nil {
			panic(err)
		}
	}()
	<-stopChan
}
