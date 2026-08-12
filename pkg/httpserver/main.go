package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(port string, handler http.Handler) error {

	// Создаем безопасную конфигурацию TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13, // Минимальная версия TLS 1.3 (самая безопасная)
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1Mb
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		TLSConfig:      tlsConfig, // Явно задаем безопасную TLS конфигурацию
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
