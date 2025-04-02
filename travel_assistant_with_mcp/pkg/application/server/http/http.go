package http

import (
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
)

type Server struct {
	*server.Hertz
	logger *log.Logger
}

type Option func(s *Server)

func NewServer(server *server.Hertz, logger *log.Logger, opts ...Option) *Server {
	s := &Server{
		Hertz:  server,
		logger: logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (h *Server) Start() {
	if err := h.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		h.logger.Sugar().Fatalf("listen: %s\n", err)
	}
}

func (h *Server) Stop() {
	h.logger.Sugar().Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	h.Spin()

	h.logger.Sugar().Info("Server exiting")
}
