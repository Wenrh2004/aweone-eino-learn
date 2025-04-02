package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/config"
)

func CreateMCPClients(
	config *config.MCPConfig,
) (map[string]*client.StdioMCPClient, error) {
	clients := make(map[string]*client.StdioMCPClient)

	for name, server := range config.MCPServers {
		var env []string
		for k, v := range server.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cli, err := client.NewStdioMCPClient(
			server.Command,
			env,
			server.Args...)
		if err != nil {
			for _, c := range clients {
				c.Close()
			}
			return nil, fmt.Errorf(
				"failed to create MCP client for %s: %w",
				name,
				err,
			)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logger.Info("Initializing server", zap.String("server", name))
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "eino-mcp-client",
			Version: "0.1.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = cli.Initialize(ctx, initRequest)
		if err != nil {
			cli.Close()
			for _, c := range clients {
				c.Close()
			}
			return nil, fmt.Errorf(
				"failed to initialize MCP client for %s: %w",
				name,
				err,
			)
		}

		clients[name] = cli
	}

	return clients, nil
}
