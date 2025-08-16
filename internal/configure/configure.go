package configure

import (
	"context"
	"fmt"
	"os"

	qdrant "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"test-ragger/internal/configure/config"
	"test-ragger/internal/models"
	"test-ragger/internal/usecase/ingest"
	"test-ragger/internal/usecase/search"
	"test-ragger/internal/utils"
	"test-ragger/internal/utils/chunker"
	"test-ragger/internal/utils/htmlx"
	"test-ragger/internal/utils/prompt"
)

// Container holds all application dependencies
type Container struct {
	Config config.Config

	// Ingest dependencies
	IngestEmbeddingClient        ingest.EmbeddingClient
	IngestQdrantCollectionClient ingest.QdrantCollectionClient
	IngestQdrantPointsClient     ingest.QdrantPointsClient
	IngestHTMLParser             ingest.HTMLParser
	IngestTextChunker            ingest.TextChunker

	// Search dependencies
	SearchEmbeddingClient    search.EmbeddingClient
	SearchQdrantPointsClient search.QdrantPointsClient
	SearchPromptBuilder      search.PromptBuilder

	// Internal connections (for cleanup)
	grpcConn *grpc.ClientConn
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

// NewContainer creates and configures all dependencies
func NewContainer(ctx context.Context, args []string) (*Container, error) {
	cfg, err := config.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// OpenAI client
	embeddingClient := openai.NewClient(utils.MustEnv("OPENAI_API_KEY"))

	// Qdrant gRPC connection
	conn, err := grpc.NewClient(cfg.QdrantGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect to qdrant: %w", err)
	}

	// Qdrant clients
	collectionsClient := qdrant.NewCollectionsClient(conn)
	pointsClient := qdrant.NewPointsClient(conn)

	// Services
	htmlParser := &htmlParserImpl{}
	textChunker := chunker.New()
	promptBuilder := &promptBuilderImpl{}

	return &Container{
		Config: cfg,

		// Ingest dependencies
		IngestEmbeddingClient:        &openaiClientAdapter{client: embeddingClient},
		IngestQdrantCollectionClient: &qdrantCollectionClientAdapter{client: collectionsClient},
		IngestQdrantPointsClient:     &qdrantPointsClientAdapter{client: pointsClient},
		IngestHTMLParser:             htmlParser,
		IngestTextChunker:            textChunker,

		// Search dependencies
		SearchEmbeddingClient:    &openaiClientAdapter{client: embeddingClient},
		SearchQdrantPointsClient: &qdrantPointsClientAdapter{client: pointsClient},
		SearchPromptBuilder:      promptBuilder,

		grpcConn: conn,
	}, nil
}

// Implementation adapters

type htmlParserImpl struct{}

func (h *htmlParserImpl) ToText(ctx context.Context, path string) (string, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer f.Close()
	return htmlx.ToText(f, path)
}

type promptBuilderImpl struct{}

func (p *promptBuilderImpl) Build(query string, hits []models.Hit) string {
	return prompt.Build(query, hits)
}

// Client adapters

type openaiClientAdapter struct {
	client *openai.Client
}

func (a *openaiClientAdapter) CreateEmbeddings(ctx context.Context, req openai.EmbeddingRequestConverter) (openai.EmbeddingResponse, error) {
	return a.client.CreateEmbeddings(ctx, req)
}

type qdrantCollectionClientAdapter struct {
	client qdrant.CollectionsClient
}

func (a *qdrantCollectionClientAdapter) Get(ctx context.Context, req *qdrant.GetCollectionInfoRequest) (*qdrant.GetCollectionInfoResponse, error) {
	return a.client.Get(ctx, req)
}

func (a *qdrantCollectionClientAdapter) Create(ctx context.Context, req *qdrant.CreateCollection) (*qdrant.CollectionOperationResponse, error) {
	return a.client.Create(ctx, req)
}

type qdrantPointsClientAdapter struct {
	client qdrant.PointsClient
}

func (a *qdrantPointsClientAdapter) Upsert(ctx context.Context, req *qdrant.UpsertPoints) (*qdrant.PointsOperationResponse, error) {
	return a.client.Upsert(ctx, req)
}

func (a *qdrantPointsClientAdapter) Search(ctx context.Context, req *qdrant.SearchPoints) (*qdrant.SearchResponse, error) {
	return a.client.Search(ctx, req)
}
