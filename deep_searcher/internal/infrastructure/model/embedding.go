package model

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/spf13/viper"
)

func NewEmbedding(conf *viper.Viper) embedding.Embedder {
	ctx := context.Background()
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		BaseURL: conf.GetString("app.embedding.base_url"),
		Region:  conf.GetString("app.embedding.region"),
		APIKey:  conf.GetString("app.embedding.api_key"),
		Model:   conf.GetString("app.embedding.model"),
	})
	if err != nil {
		panic(err)
	}
	return emb
}
