package ingest

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

// QdrantCollectionClient handles collection operations
type QdrantCollectionClient interface {
	Get(ctx context.Context, req *qdrant.GetCollectionInfoRequest) (*qdrant.GetCollectionInfoResponse, error)
	Create(ctx context.Context, req *qdrant.CreateCollection) (*qdrant.CollectionOperationResponse, error)
}

// QdrantPointsClient handles point operations
type QdrantPointsClient interface {
	Upsert(ctx context.Context, req *qdrant.UpsertPoints) (*qdrant.PointsOperationResponse, error)
}

// HTMLParser extracts text from HTML content
type HTMLParser interface {
	ToText(ctx context.Context, path string) (text, title string, err error)
}

// TextChunker splits text into chunks
type TextChunker interface {
	ChunkText(text string, size, overlap int) []models.ChunkInfo
}
