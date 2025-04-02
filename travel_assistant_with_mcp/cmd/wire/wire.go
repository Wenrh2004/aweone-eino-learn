//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/adapter"
	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/adapter/server"
	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/domain/agent"
	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/infrastructure/llm"
	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/infrastructure/tool"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/application"
	httpserver "github.com/Wenrh2004/travel_assistant_with_mcp/pkg/application/server/http"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
)

var infrastructureSet = wire.NewSet(
	tool.NewMCPToolsService,
	llm.NewChatModelService,
)

var domainSet = wire.NewSet(
	agent.NewDomain,
)

var adapterSet = wire.NewSet(
	adapter.NewLLMHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
)

// build App
func newApp(
	httpServer *httpserver.Server,
	conf *viper.Viper,
) *application.App {
	return application.NewApp(
		application.WithServer(httpServer),
		application.WithName(conf.GetString("app.name")),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*application.App, func(), error) {
	panic(wire.Build(
		infrastructureSet,
		domainSet,
		adapterSet,
		serverSet,
		newApp,
	))
}
