package server

import (
	"fmt"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"sync/atomic"
	"time"
)

type promPushGatewayServer struct {
	cfg     *config.Config
	pusher  *push.Pusher
	running int32
	closeCh chan struct{}
}

func NewPromPushGatewayServer(cfg *config.Config, g prometheus.Gatherer) interfaces.MeterServer {
	pushServer := promPushGatewayServer{
		cfg:     cfg,
		running: 0,
		closeCh: make(chan struct{}),
	}
	pushServer.pusher = push.New(cfg.PushGateway.GatewayAddress, cfg.LocalIP).Gatherer(g)

	return &pushServer
}

func (s *promPushGatewayServer) Start() {
	if !(atomic.CompareAndSwapInt32(&s.running, 0, 1)) {
		return
	}
	go s.push()
}

func (s *promPushGatewayServer) Stop() {
	if !(atomic.CompareAndSwapInt32(&s.running, 1, 0)) {
		return
	}
	s.closeCh <- struct{}{}
}

func (s *promPushGatewayServer) push() {
	pushTicker := time.NewTicker(s.cfg.PushGateway.PushPeriod)
	defer pushTicker.Stop()

	now := time.Now()
	if err := s.pusher.Push(); err != nil {
		s.cfg.WriteErrorOrNot("failed to push to gateway: " + err.Error())
	} else {
		s.cfg.WriteInfoOrNot(fmt.Sprintf("successfully pushed to gateway, tick = %s, now = %s", time.Now().Sub(now), time.Now().Local().String()))
	}
	for {
		select {
		case <-pushTicker.C:
			now = time.Now()
			if err := s.pusher.Push(); err != nil {
				s.cfg.WriteErrorOrNot("failed to push to gateway: " + err.Error())
			} else {
				s.cfg.WriteInfoOrNot(fmt.Sprintf("successfully pushed to gateway, tick = %s, now = %s", time.Now().Sub(now), time.Now().Local().String()))
			}
		case <-s.closeCh:
			s.cfg.WriteInfoOrNot("push gateway server is closed")
			return
		}
	}
}
