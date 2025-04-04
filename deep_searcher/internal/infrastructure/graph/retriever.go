package graph

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"github.com/Wenrh2004/deep_searcher/pkg/util"
)

// chat template
var (
	subQueryChatTemplate = prompt.FromMessages(schema.FString, &schema.Message{
		Role:    schema.User,
		Content: FollowupQueryPrompt,
	})
	intermediateAnswerChatTemplate = prompt.FromMessages(schema.FString, &schema.Message{
		Role:    schema.User,
		Content: IntermediateAnswerPrompt,
	})
	getSupportedDocsChatTemplate = prompt.FromMessages(schema.FString, &schema.Message{
		Role:    schema.User,
		Content: GetSupportedDocsPrompt,
	})
	reflectionChatTemplate = prompt.FromMessages(schema.FString, &schema.Message{
		Role:    schema.User,
		Content: ReflectionPrompt,
	})
)

// handler
var (
	supportedDocsPostHandler = func(ctx context.Context, output []*schema.Document, state *state) ([]*schema.Document, error) {
		state.allRetrievedResults = append(state.allRetrievedResults, output...)
		return output, nil
	}
	subQueryPreHandler = func(ctx context.Context, input map[string]any, state *state) (map[string]any, error) {
		input["intermediate_contexts"] = state.intermediateContexts
		return input, nil
	}
)

// checker
var (
	enoughInfoChecker = func(info *schema.Message) (bool, error) {
		if strings.ToLower(info.Content) == "yes" {
			return true, nil
		} else if strings.ToLower(info.Content) == "no" {
			return false, nil
		} else {
			return false, errors.New("[ChainOfRag.Retriever.enoughInfoChecker] info is not match")
		}
	}
	retrieverBranchCondition = func(_ context.Context, input *schema.Message) (endNode string, err error) {
		if isEnoughInfo, err := enoughInfoChecker(input); err != nil {
			return "", err
		} else if !isEnoughInfo {
			return "get_sub_query", nil
		}
		return "collate_output", nil
	}
)

type RetrieverGraph struct {
	chatModel model.ChatModel
	embedder  embedding.Embedder
	vectorDB  retriever.Retriever
}

type state struct {
	intermediateContexts []string
	allRetrievedResults  []*schema.Document
}

func NewRetrieverGraph(chatModel model.ChatModel, emb embedding.Embedder, vdb retriever.Retriever) *RetrieverGraph {
	return &RetrieverGraph{
		chatModel: chatModel,
		embedder:  emb,
		vectorDB:  vdb,
	}
}

func (r *RetrieverGraph) getReflectGetSubQueryGraph() (*compose.Graph[map[string]any, *schema.Message], error) {
	graph := compose.NewGraph[map[string]any, *schema.Message]()
	if err := graph.AddChatModelNode("chat_model", r.chatModel, compose.WithOutputKey("sub_query")); err != nil {
		return nil, err
	}
	if err := graph.AddChatTemplateNode("prompt", subQueryChatTemplate); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "prompt"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("prompt", "chat_model"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("chat_model", compose.END); err != nil {
		return nil, err
	}

	return graph, nil
}

func (r *RetrieverGraph) getRetrieverGraph() (*compose.Graph[*schema.Message, []*schema.Document], error) {
	graph := compose.NewGraph[*schema.Message, []*schema.Document]()
	if err := graph.AddRetrieverNode("retriever", r.vectorDB); err != nil {
		return nil, err
	}
	if err := graph.AddLambdaNode("deduplicate", compose.InvokableLambda[[]*schema.Document, []*schema.Document](
		func(ctx context.Context, input []*schema.Document) ([]*schema.Document, error) {
			return util.DeduplicateResults(input), nil
		}), compose.WithOutputKey("retrieved_documents")); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "retriever"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("retriever", "deduplicate"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("deduplicate", compose.END); err != nil {
		return nil, err
	}

	return graph, nil
}

func (r *RetrieverGraph) getIntermediateAnswerGraph() (*compose.Graph[map[string]any, *schema.Message], error) {
	graph := compose.NewGraph[map[string]any, *schema.Message]()
	if err := graph.AddChatModelNode("chat_model", r.chatModel, compose.WithOutputKey("answer")); err != nil {
		return nil, err
	}
	if err := graph.AddChatTemplateNode("prompt", intermediateAnswerChatTemplate); err != nil {
		return nil, err
	}
	if err := graph.AddLambdaNode("collate_output", compose.InvokableLambda(util.IntermediateAnswereModifier)); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "prompt"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("prompt", "chat_model"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "collate_output"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("chat_model", "collate_output"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("collate_output", compose.END); err != nil {
		return nil, err
	}

	return graph, nil
}

func (r *RetrieverGraph) getSupportedDocsGraph() (*compose.Graph[map[string]any, []*schema.Document], error) {
	graph := compose.NewGraph[map[string]any, []*schema.Document]()
	if err := graph.AddChatModelNode("chat_model", r.chatModel, compose.WithOutputKey("answer")); err != nil {
		return nil, err
	}
	if err := graph.AddChatTemplateNode("prompt", getSupportedDocsChatTemplate); err != nil {
		return nil, err
	}
	if err := graph.AddLambdaNode("select_the_supported_docs", compose.InvokableLambda(func(ctx context.Context, input map[string]any) (output []*schema.Document, err error) {
		supportedDocIndices, err := util.LiteralEval(input["answer"].(string))
		if err != nil {
			return nil, err
		}
		retrievedDocs, ok := input["retrieved_documents"]
		if !ok {
			return nil, errors.New("retrieved_documents is not found")
		}
		retrievedResults := retrievedDocs.([]*schema.Document)
		supportedRetrievedResults := util.GetSupportedDocuments(retrievedResults, supportedDocIndices)
		return supportedRetrievedResults, nil
	})); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "prompt"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("prompt", "chat_model"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("chat_model", "select_the_supported_docs"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "select_the_supported_docs"); err != nil {
		return nil, err
	}
	return graph, nil
}

func (r *RetrieverGraph) getCheckEnoughInfoGraph() (*compose.Graph[map[string]any, *schema.Message], error) {
	graph := compose.NewGraph[map[string]any, *schema.Message]()
	if err := graph.AddChatModelNode("chat_model", r.chatModel); err != nil {
		return nil, err
	}
	if err := graph.AddChatTemplateNode("prompt", reflectionChatTemplate); err != nil {
		return nil, err
	}
	if err := graph.AddEdge(compose.START, "prompt"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("prompt", "chat_model"); err != nil {
		return nil, err
	}
	if err := graph.AddEdge("chat_model", compose.END); err != nil {
		return nil, err
	}
	return graph, nil
}

func (r *RetrieverGraph) GewRetrieverGraph() (*compose.Graph[map[string]any, map[string]any], error) {
	subQueryGraph, err := r.getReflectGetSubQueryGraph()
	if err != nil {
		return nil, err
	}
	retrieveGraph, err := r.getRetrieverGraph()
	if err != nil {
		return nil, err
	}
	answerGraph, err := r.getIntermediateAnswerGraph()
	if err != nil {
		return nil, err
	}
	supportedDocsGraph, err := r.getSupportedDocsGraph()
	if err != nil {
		return nil, err
	}
	checkEnoughInfoGraph, err := r.getCheckEnoughInfoGraph()
	if err != nil {
		return nil, err
	}
	// build the retriever graph
	g := compose.NewGraph[map[string]any, map[string]any](compose.WithGenLocalState(func(ctx context.Context) *state {
		return &state{
			intermediateContexts: []string{},
			allRetrievedResults:  []*schema.Document{},
		}
	}))
	if err := g.AddGraphNode("get_sub_query", subQueryGraph, compose.WithStatePreHandler(subQueryPreHandler)); err != nil {
		return nil, err
	}
	if err := g.AddGraphNode("retrieve", retrieveGraph); err != nil {
		return nil, err
	}
	if err := g.AddGraphNode("get_intermediate_answer", answerGraph, compose.WithStatePostHandler(func(ctx context.Context, output []string, state *state) ([]string, error) {
		state.intermediateContexts = append(state.intermediateContexts, output...)
		return state.intermediateContexts, nil
	})); err != nil {
		return nil, err
	}
	if err := g.AddGraphNode("get_supported_docs", supportedDocsGraph, compose.WithStatePostHandler(supportedDocsPostHandler)); err != nil {
		return nil, err
	}
	if err := g.AddGraphNode("check_enough_info", checkEnoughInfoGraph); err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode("collate_output", compose.InvokableLambda(util.GetCollateOutput)); err != nil {
		return nil, err
	}
	// 1. get sub query
	if err := g.AddEdge(compose.START, "get_sub_query"); err != nil {
		return nil, err
	}
	// 2. retrieve answer
	if err := g.AddEdge("get_sub_query", "retrieve"); err != nil {
		return nil, err
	}
	// 3. get intermediate answer
	if err := g.AddEdge("get_sub_query", "get_intermediate_answer"); err != nil {
	}
	if err := g.AddEdge("retrieve", "get_intermediate_answer"); err != nil {
		return nil, err
	}
	// 4. get supported documents
	if err := g.AddEdge("get_sub_query", "get_supported_docs"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("retrieve", "get_supported_docs"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("get_intermediate_answer", "get_supported_docs"); err != nil {
		return nil, err
	}
	// 5. check the info is enough
	if err := g.AddEdge(compose.START, "check_enough_info"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("get_intermediate_answer", "check_enough_info"); err != nil {
		return nil, err
	}
	// 6. determine whether you need to continue iterating
	if err := g.AddBranch("get_sub_query", compose.NewGraphBranch(retrieverBranchCondition, map[string]bool{compose.END: true})); err != nil {
		return nil, err
	}
	// 7. collate output
	if err := g.AddEdge("collate_output", compose.END); err != nil {
		return nil, err
	}
	return g, nil
}
