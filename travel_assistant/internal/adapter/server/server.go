package server

import (
	hertzserver "github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Wenrh2004/travel_assistant/internal/adapter"
	"github.com/Wenrh2004/travel_assistant/pkg/application/server/http"
	"github.com/Wenrh2004/travel_assistant/pkg/util/log"
)

func NewHTTPServer(logger *log.Logger, handler *adapter.LLMHandler) *http.Server {
	h := http.NewServer(hertzserver.Default(), logger)
	h.POST("/chat", handler.Chat)
	return h
}
