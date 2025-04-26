package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Wenrh2004/travel_assistant_with_mcp/pkg/util/log"
)

func NewConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}

type MCPConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

func GetServerConfig(conf *viper.Viper, logger *log.Logger) (*MCPConfig, error) {
	confPath := conf.GetString("app.mcp.server.config")
	if confPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		confPath = filepath.Join(homeDir, "mcp.json")
	}
	// Check if config file exists
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := MCPConfig{
			MCPServers: make(map[string]ServerConfig),
		}

		// Create the file with default config
		configData, err := sonic.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error creating default config: %w", err)
		}

		if err := os.WriteFile(confPath, configData, 0644); err != nil {
			return nil, fmt.Errorf("error writing default config file: %w", err)
		}

		logger.Info("Created default config file", zap.String("path", confPath))
		return &defaultConfig, nil
	}

	// Read existing config
	configData, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading config file %s: %w",
			confPath,
			err,
		)
	}

	var config MCPConfig
	if err := sonic.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}
