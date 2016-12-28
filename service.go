package kate

import (
	"net"
	"net/http"
	"time"

	"github.com/k81/kate/context"
	"github.com/k81/kate/log"
	"github.com/k81/kate/utils"
)

type ServiceOption struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

var (
	DefaultOption = &ServiceOption{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
)

type KateService struct {
	Addr     string
	ctx      context.Context
	listener net.Listener
	router   *Router
	server   *http.Server
	option   *ServiceOption
}

func Service(ctx context.Context, router *Router, option *ServiceOption) (srv *KateService) {
	if option == nil {
		option = DefaultOption
	}

	srv = &KateService{
		ctx:    ctx,
		router: router,
		option: option,
		server: &http.Server{
			Handler:        router,
			ReadTimeout:    option.ReadTimeout,
			WriteTimeout:   option.WriteTimeout,
			MaxHeaderBytes: option.MaxHeaderBytes,
		},
	}

	return
}

func (s *KateService) Listen(addr string) (err error) {
	s.Addr = addr
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		log.Fatal(s.ctx, "listen", "addr", addr, "error", err)
	}

	log.Info(s.ctx, "service started", "addr", s.Addr)
	go s.run()
	return
}

func (s *KateService) run() {
	if err := s.server.Serve(s.listener); err != nil && !utils.IsErrClosing(err) {
		log.Fatal(s.ctx, "start", "error", err)
	}
}

func (s *KateService) Stop() (err error) {
	if err = s.listener.Close(); err != nil {
		log.Error(s.ctx, "stopping service", "addr", s.Addr, "error", err)
		return
	}
	s.router.Wait()
	log.Info(s.ctx, "service stopped", "addr", s.Addr)
	return
}
