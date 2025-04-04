package vector

import (
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/spf13/viper"
)

func NewRetriever(conf *viper.Viper, emb embedding.Embedder) retriever.Retriever {
	ctx := context.Background()
	cli, err := client.NewClient(ctx, client.Config{
		Address:  conf.GetString("app.vector.addr"),
		Username: conf.GetString("app.vector.username"),
		Password: conf.GetString("app.vector.password"),
	})
	if err != nil {
		panic(err)
	}
	r, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:         cli,
		Collection:     conf.GetString("app.vector.collection"),
		VectorField:    conf.GetString("app.vector.field"),
		TopK:           conf.GetInt("app.vector.top_k"),
		ScoreThreshold: conf.GetFloat64("app.vector.score_threshold"),
		Embedding:      emb,
	})
	if err != nil {
		panic(err)
	}
	return r
}
