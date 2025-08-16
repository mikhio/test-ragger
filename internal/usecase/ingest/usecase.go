package ingest

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	qdrant "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"

	"test-ragger/internal/configure/config"
	"test-ragger/internal/utils"
)

// Usecase handles HTML ingestion into Qdrant
type Usecase struct {
	embeddingClient        EmbeddingClient
	qdrantCollectionClient QdrantCollectionClient
	qdrantPointsClient     QdrantPointsClient
	htmlParser             HTMLParser
	textChunker            TextChunker
}

// New creates new ingest usecase
func New(
	embeddingClient EmbeddingClient,
	qdrantCollectionClient QdrantCollectionClient,
	qdrantPointsClient QdrantPointsClient,
	htmlParser HTMLParser,
	textChunker TextChunker,
) *Usecase {
	return &Usecase{
		embeddingClient:        embeddingClient,
		qdrantCollectionClient: qdrantCollectionClient,
		qdrantPointsClient:     qdrantPointsClient,
		htmlParser:             htmlParser,
		textChunker:            textChunker,
	}
}

// Run executes HTML ingestion process
func (u *Usecase) Run(ctx context.Context, htmlDir string, model openai.EmbeddingModel) error {
	cfg, _ := config.FromContext(ctx)

	slog.Info("Ensuring collection exists", "collection", cfg.Collection, "dimension", cfg.EmbeddingDim)
	if err := u.ensureCollection(ctx, cfg.Collection, cfg.EmbeddingDim); err != nil {
		return fmt.Errorf("ensureCollection: %w", err)
	}
	slog.Info("Collection ready for ingestion")

	return filepath.Walk(htmlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		slog.Info("Parsing HTML file", "path", path)
		text, title, err := u.htmlParser.ToText(ctx, path)
		if err != nil {
			return err
		}
		if len(text) == 0 {
			slog.Info("Skipping empty file", "path", path)
			return nil
		}
		slog.Info("Parsed HTML to text", "title", title, "characters", len(text))

		docID := "doc_" + utils.Sha1Hex(path)
		slog.Info("Chunking text", "chunk_size", cfg.ChunkSize, "overlap", cfg.ChunkOverlap)
		chunks := u.textChunker.ChunkText(text, cfg.ChunkSize, cfg.ChunkOverlap)
		slog.Info("Created chunks", "count", len(chunks))
		batch := make([]*qdrant.PointStruct, 0, len(chunks))

		slog.Info("Creating embeddings", "chunks_count", len(chunks))
		for i, c := range chunks {
			slog.Debug("Processing chunk", "chunk_index", i, "chunk_length", len(c.Text))
			// create embedding
			res, err := u.embeddingClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
				Model: model,
				Input: []string{c.Text},
			})
			if err != nil {
				return fmt.Errorf("embedding: %w", err)
			}
			if i%10 == 0 && i > 0 {
				slog.Info("Embeddings progress", "completed", i, "total", len(chunks))
			}

			vec := res.Data[0].Embedding
			if len(vec) != cfg.EmbeddingDim {
				return fmt.Errorf("dim mismatch: got %d want %d", len(vec), cfg.EmbeddingDim)
			}

			// create payload
			payload := map[string]*qdrant.Value{
				"doc_id":      {Kind: &qdrant.Value_StringValue{StringValue: docID}},
				"chunk_id":    {Kind: &qdrant.Value_StringValue{StringValue: c.ChunkID}},
				"title":       {Kind: &qdrant.Value_StringValue{StringValue: title}},
				"path":        {Kind: &qdrant.Value_StringValue{StringValue: path}},
				"start":       {Kind: &qdrant.Value_DoubleValue{DoubleValue: float64(c.Start)}},
				"end":         {Kind: &qdrant.Value_DoubleValue{DoubleValue: float64(c.End)}},
				"text":        {Kind: &qdrant.Value_StringValue{StringValue: c.Text}},
				"ingested_at": {Kind: &qdrant.Value_StringValue{StringValue: time.Now().Format(time.RFC3339)}},
				"lang":        {Kind: &qdrant.Value_StringValue{StringValue: "ru"}},
				"type":        {Kind: &qdrant.Value_StringValue{StringValue: "html"}},
			}

			batch = append(batch, &qdrant.PointStruct{
				Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: docID + "_" + c.ChunkID}},
				Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: vec}}},
				Payload: payload,
			})
		}

		slog.Info("Upserting points to Qdrant", "points_count", len(batch), "collection", cfg.Collection)
		_, err = u.qdrantPointsClient.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: cfg.Collection,
			Points:         batch,
			Wait:           pointer.To(true),
		})
		if err != nil {
			return err
		}

		slog.Info("Successfully ingested file", "path", path, "chunks", len(chunks))

		return nil
	})
}

// Убеждаемся что создана коллекция, если нет - то создаем
func (u *Usecase) ensureCollection(ctx context.Context, collection string, dim int) error {
	_, err := u.qdrantCollectionClient.Get(ctx, &qdrant.GetCollectionInfoRequest{CollectionName: collection})
	if err == nil {
		slog.Info("Collection already exists", "collection", collection)
		return nil
	}
	slog.Info("Creating new collection", "collection", collection, "dimension", dim)
	_, err = u.qdrantCollectionClient.Create(ctx, &qdrant.CreateCollection{
		CollectionName: collection,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     uint64(dim),
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})
	if err == nil {
		slog.Info("Collection created successfully", "collection", collection)
	}
	return err
}
