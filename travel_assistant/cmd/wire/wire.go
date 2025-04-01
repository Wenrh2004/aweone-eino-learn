//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant/internal/adapter"
	"github.com/Wenrh2004/travel_assistant/internal/adapter/server"
	"github.com/Wenrh2004/travel_assistant/internal/domain/agent"
	"github.com/Wenrh2004/travel_assistant/internal/infrastructure/llm"
	"github.com/Wenrh2004/travel_assistant/internal/infrastructure/third"
	tool2 "github.com/Wenrh2004/travel_assistant/internal/infrastructure/third/tool"
	"github.com/Wenrh2004/travel_assistant/pkg/application"
	httpserver "github.com/Wenrh2004/travel_assistant/pkg/application/server/http"
	"github.com/Wenrh2004/travel_assistant/pkg/util/log"
)

var infrastructureSet = wire.NewSet(
	third.NewAmapClient,
	tool2.NewPOIService,
	tool2.NewRouterService,
	tool2.NewWeatherService,
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
