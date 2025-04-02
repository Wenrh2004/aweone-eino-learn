package server

import (
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/adapter"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/application/server/http"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
)

func NewHTTPServer(logger *log.Logger, handler *adapter.LLMHandler) *http.Server {
	h := http.NewServer(hertzserver.Default(), logger)
	h.POST("/chat", handler.Chat)
	return h
}
