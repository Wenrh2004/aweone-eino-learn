package model

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/spf13/viper"
)

func NewChatModel(conf *viper.Viper) model.ChatModel {
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
	return chatModel
}
