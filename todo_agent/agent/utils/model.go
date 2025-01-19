package utils

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/spf13/viper"
)

func NewChatModel(ctx context.Context, conf *viper.Viper) *openai.ChatModel {
	// 创建并配置 ChatModel
	temp := float32(conf.GetFloat64("openai.temperature"))
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:       conf.GetString("openai.model"),
		APIKey:      conf.GetString("openai.api_key"),
		BaseURL:     conf.GetString("openai.base_url"),
		Temperature: &temp,
	})
	if err != nil {
		panic(err)
	}
	return chatModel
}
