package search

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"

	"test-ragger/internal/configure/config"
	"test-ragger/internal/models"
	"test-ragger/internal/utils"
)

// Usecase handles search operations
type Usecase struct {
	embeddingClient    EmbeddingClient
	qdrantPointsClient QdrantPointsClient
	promptBuilder      PromptBuilder
}

// New creates new search usecase
func New(
	embeddingClient EmbeddingClient,
	qdrantPointsClient QdrantPointsClient,
	promptBuilder PromptBuilder,
) *Usecase {
	return &Usecase{
		embeddingClient:    embeddingClient,
		qdrantPointsClient: qdrantPointsClient,
		promptBuilder:      promptBuilder,
	}
}

// Search executes search query and returns results
func (u *Usecase) Search(ctx context.Context, query string, topK uint64, model openai.EmbeddingModel, langFilter string) ([]models.Hit, error) {
	cfg, _ := config.FromContext(ctx)

	// create query embedding
	emb, err := u.embeddingClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{Model: model, Input: []string{query}})
	if err != nil {
		return nil, err
	}
	vec := emb.Data[0].Embedding

	// build filter if needed
	var filter *qdrant.Filter
	if langFilter != "" {
		filter = &qdrant.Filter{Must: []*qdrant.Condition{{ConditionOneOf: &qdrant.Condition_Field{Field: &qdrant.FieldCondition{Key: "lang", Match: &qdrant.Match{MatchValue: &qdrant.Match_Keyword{Keyword: langFilter}}}}}}}
	}

	// execute search
	resp, err := u.qdrantPointsClient.Search(ctx, &qdrant.SearchPoints{
		CollectionName: cfg.Collection,
		Vector:         vec,
		Limit:          topK,
		Params:         &qdrant.SearchParams{HnswEf: utils.Uint64Ptr(128)},
		Filter:         filter,
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, err
	}

	// convert to hits
	hits := make([]models.Hit, 0, len(resp.Result))
	for _, r := range resp.Result {
		pl := r.Payload
		hits = append(hits, models.Hit{
			Score:   r.GetScore(),
			Title:   fmt.Sprintf("%v", pl["title"]),
			Text:    fmt.Sprintf("%v", pl["text"]),
			Path:    fmt.Sprintf("%v", pl["path"]),
			DocID:   fmt.Sprintf("%v", pl["doc_id"]),
			ChunkID: fmt.Sprintf("%v", pl["chunk_id"]),
		})
	}
	return hits, nil
}

// BuildPrompt creates LLM prompt from search results
func (u *Usecase) BuildPrompt(query string, hits []models.Hit) string {
	return u.promptBuilder.Build(query, hits)
}
