package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/spf13/viper"
)

// ChatModelService llm chat client
type ChatModelService struct {
	ChatModel model.ToolCallingChatModel
}

// NewChatModelService DeepSeek Service client
func NewChatModelService(conf *viper.Viper) *ChatModelService {
	ctx := context.Background()
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  conf.GetString("app.llm.api_key"),
		BaseURL: conf.GetString("app.llm.base_url"),
		Region:  conf.GetString("app.llm.region"),
		Model:   conf.GetString("app.llm.model"),
	})
	if err != nil {
		panic(err)
	}
	return &ChatModelService{
		ChatModel: chatModel,
	}
}
