package agent

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/infrastructure/llm"
	"github.com/Wenrh2004/travel_assistant_with_mcp/internal/infrastructure/tool"
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
	service *tool.MCPToolsService,
) ChatService {
	ctx := context.Background()
	mcpTools, err := service.GetMCPTool()
	if err != nil {
		panic(err)
	}
	agent, err := react.NewAgent(
		ctx,
		&react.AgentConfig{
			Model: chatModel.ChatModel,
			ToolsConfig: compose.ToolsNodeConfig{
				Tools: mcpTools,
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
