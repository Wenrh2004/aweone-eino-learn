package tool

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/spf13/viper"

	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/config"
	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
	mcpclient "github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/mcp"
)

type MCPToolsService struct {
	cli map[string]*client.StdioMCPClient
}

func NewMCPToolsService(conf *viper.Viper, logger *log.Logger) *MCPToolsService {
	serverConfig, err := config.GetServerConfig(conf, logger)
	if err != nil {
		panic(err)
	}
	clients, err := mcpclient.CreateMCPClients(serverConfig)
	if err != nil {
		panic(err)
	}
	return &MCPToolsService{
		cli: clients,
	}
}

func (t *MCPToolsService) GetMCPTool() ([]tool.BaseTool, error) {
	var tools []tool.BaseTool
	for _, mcpClient := range t.cli {
		baseTool, err := mcp.GetTools(context.Background(), &mcp.Config{
			Cli: mcpClient,
		})
		if err != nil {
			return nil, err
		}
		tools = append(tools, baseTool...)
	}

	return tools, nil
}
