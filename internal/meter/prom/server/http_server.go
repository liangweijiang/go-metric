package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"net/http"
	"net/http/pprof"
	"strconv"
	"sync/atomic"
)

// promHttpServer encapsulates the necessary components to run an HTTP server for exposing Prometheus metrics.
// It includes the handler for metrics export, the underlying HTTP server instance, configuration settings,
// a channel for triggering a shutdown, and an atomic flag indicating the server's running state.
type promHttpServer struct {
	exporterHandler http.Handler
	server          *http.Server
	cfg             *config.Config
	closeCh         chan struct{}
	running         int32
}

// NewPromHttpServer initializes a new Prometheus HTTP server based on the provided configuration and exporter handler.
// It sets up the necessary structures to start and stop the server, including configurations and channels for control.
// Returns a MeterServer interface which can be used to manage the lifecycle of the HTTP server for metrics exposure.
func NewPromHttpServer(cfg *config.Config, exporterHandler http.Handler) interfaces.MeterServer {

	server := promHttpServer{
		cfg:             cfg,
		exporterHandler: exporterHandler,
		running:         0,
		closeCh:         make(chan struct{}),
	}

	return &server
}

// Start initializes and begins listening for HTTP requests on the configured Prometheus port.
// It sets up various endpoints like health check, metrics retrieval, and profiling routes.
// If the server is already running, the method will not restart it.
// A shutdown hook is also set up to gracefully stop the server when requested.
func (s *promHttpServer) Start() {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		s.cfg.WriteInfoOrNot("prom http server is already running")
		return
	}
	s.cfg.WriteInfoOrNot(fmt.Sprintf("starting prom http server, port:%d", s.cfg.PrometheusPort))
	mux := http.NewServeMux()
	logRoute := func(route string) string {
		s.cfg.WriteInfoOrNot(fmt.Sprintf("http handler, method:Get, uri:%s", route))
		return route
	}
	mux.HandleFunc(logRoute("/actuator/health"), s.healthCheck)
	mux.HandleFunc(logRoute("/metrics"), func(w http.ResponseWriter, r *http.Request) {
		if s.exporterHandler != nil {
			s.exporterHandler.ServeHTTP(w, r)
		}
	})
	mux.HandleFunc(logRoute("/debug/pprof/"), pprof.Index)
	mux.HandleFunc(logRoute("/debug/pprof/cmdline"), pprof.Cmdline)
	mux.HandleFunc(logRoute("/debug/pprof/profile"), pprof.Profile)
	mux.HandleFunc(logRoute("/debug/pprof/symbol"), pprof.Symbol)
	mux.HandleFunc(logRoute("/debug/pprof/trace"), pprof.Trace)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.PrometheusPort),
		Handler: mux,
	}
	go s.startHTTPServer()
	go func() {
		select {
		case <-s.closeCh:
			s.cfg.WriteInfoOrNot("prom http server is shutting down")
			err := s.server.Shutdown(context.Background())
			if err != nil {
				s.cfg.WriteErrorOrNot(fmt.Sprintf("failed to shutdown prom http server with error: %s", err.Error()))
				return
			}
		}
	}()
}

// Stop halts the promHTTP server operation by setting its running state to stopped, logging the action, and signaling the close channel to initiate a shutdown sequence.
func (s *promHttpServer) Stop() {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		s.cfg.WriteInfoOrNot("prom http server is already stopped")
		return
	}
	s.cfg.WriteInfoOrNot("stopping prom http server")
	s.closeCh <- struct{}{}
}

// startHTTPServer initiates the HTTP server to serve Prometheus metrics and other endpoints.
// It listens on the configured PrometheusPort and handles errors during startup, logging them accordingly.
func (s *promHttpServer) startHTTPServer() {
	s.cfg.WriteInfoOrNot("prom http server listen and server on: " + strconv.Itoa(s.cfg.PrometheusPort))
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.cfg.WriteErrorOrNot(fmt.Sprintf("faield to start prom http server on : %d with error: %s ",
			s.cfg.PrometheusPort, err.Error()))
	}
}

// healthCheck responds to HTTP requests with a JSON message indicating the service status is "UP".
// It sets the "Content-Type" header to "application/json" and marshals a simple JSON object with a "status" field.
// This endpoint is typically used to check the availability of the service.
func (s *promHttpServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/json")
	msg, _ := json.Marshal(map[string]interface{}{"status": "UP"})
	_, _ = w.Write(msg)
}
