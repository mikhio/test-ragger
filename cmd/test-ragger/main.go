package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"

	"test-ragger/internal/configure"
	"test-ragger/internal/configure/config"
	"test-ragger/internal/usecase/ingest"
	"test-ragger/internal/usecase/search"
	"test-ragger/internal/utils"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	// Initialize dependencies container
	container, err := configure.NewContainer(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	cfg := container.Config
	ctx = config.IntoContext(ctx, cfg)

	var model openai.EmbeddingModel
	switch cfg.Model {
	case "text-embedding-3-small", "text-embedding-3-large":
		model = openai.EmbeddingModel(cfg.Model)
	default:
		log.Fatalf("unknown model: %s", cfg.Model)
	}

	switch cfg.Mode {
	case "ingest":
		ingestUC := ingest.New(
			container.IngestEmbeddingClient,
			container.IngestQdrantCollectionClient,
			container.IngestQdrantPointsClient,
			container.IngestHTMLParser,
			container.IngestTextChunker,
		)
		if err := ingestUC.Run(ctx, cfg.HTMLDir, model); err != nil {
			log.Fatal(err)
		}
	case "search":
		if cfg.Query == "" {
			log.Fatal("-q is required in search mode")
		}
		searchUC := search.New(
			container.SearchEmbeddingClient,
			container.SearchQdrantPointsClient,
			container.SearchPromptBuilder,
		)
		hits, err := searchUC.Search(ctx, cfg.Query, cfg.TopK, model, cfg.Lang)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Query: %s\nTop-%d results:\n", cfg.Query, cfg.TopK)
		for i, h := range hits {
			fmt.Printf("#%d score=%.4f %s\n%s\npath=%s\n---\n", i+1, h.Score, h.Title, utils.Snippet(h.Text, 280), h.Path)
		}

		fmt.Println("\n--- PROMPT ---")
		fmt.Println(searchUC.BuildPrompt(cfg.Query, hits))
	default:
		log.Fatalf("unknown mode: %s", cfg.Mode)
	}
}
