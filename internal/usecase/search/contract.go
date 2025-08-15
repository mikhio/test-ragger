package search

import (
	"context"

	qdrant "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"

	"test-ragger/internal/models"
)

// EmbeddingClient represents OpenAI API client for embeddings
type EmbeddingClient interface {
	CreateEmbeddings(ctx context.Context, req openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error)
}

// QdrantPointsClient handles point operations
type QdrantPointsClient interface {
	Search(ctx context.Context, req *qdrant.SearchPoints) (*qdrant.SearchResponse, error)
}

// PromptBuilder creates LLM prompts from search results
type PromptBuilder interface {
	Build(query string, hits []models.Hit) string
}
