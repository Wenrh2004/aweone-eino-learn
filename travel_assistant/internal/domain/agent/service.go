package agent

import (
	"context"

	"github.com/cloudwego/eino-ext/devops"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant/internal/infrastructure/llm"
	tool2 "github.com/Wenrh2004/travel_assistant/internal/infrastructure/third/tool"
)

type ChatService interface {
	Query(ctx context.Context, query string, history []*schema.Message) (*schema.StreamReader[*schema.Message], error)
}

type chatService struct {
	agent *react.Agent
}

func NewDomain(
	conf *viper.Viper,
	chatModel *llm.ChatModelService,
	poiService *tool2.POIService,
	routerService *tool2.RouterService,
	weatherService *tool2.WeatherService,
) ChatService {
	ctx := context.Background()
	if err := devops.Init(ctx); err != nil {
		panic(err)
	}
	poiSearchTool, err := poiService.GetPOISearchTool()
	if err != nil {
		panic(err)
	}
	drivingRouterTool, err := routerService.GetDrivingRouterTool()
	if err != nil {
		panic(err)
	}
	walkingRouterTool, err := routerService.GetWalkingRouterTool()
	if err != nil {
		panic(err)
	}
	weatherForecastSearchTool, err := weatherService.GetWeatherForecastTool()
	if err != nil {
		panic(err)
	}
	weatherNowSearchTool, err := weatherService.GetCurrentWeatherTool()
	if err != nil {
		panic(err)
	}
	agent, err := react.NewAgent(
		ctx,
		&react.AgentConfig{
			ToolCallingModel: chatModel.ChatModel,
			ToolsConfig: compose.ToolsNodeConfig{
				Tools: []einotool.BaseTool{
					poiSearchTool,
					drivingRouterTool,
					walkingRouterTool,
					weatherForecastSearchTool,
					weatherNowSearchTool,
				},
			},
			MaxStep: conf.GetInt("app.agent.max_step"),
		},
	)
	if err != nil {
		panic(err)
	}
	return &chatService{
		agent: agent,
	}
}

func (d *chatService) Query(ctx context.Context, query string, history []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	history = append(history, schema.UserMessage(query))
	stream, err := d.agent.Stream(ctx, history)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
