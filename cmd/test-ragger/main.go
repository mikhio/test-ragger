package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
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

	// Setup structured logging
	logLevel := slog.LevelInfo
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		switch level {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "INFO":
			logLevel = slog.LevelInfo
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting test-ragger application")
	slog.Debug("Log level set", "level", logLevel.String())

	// Initialize dependencies container
	slog.Info("Loading configuration from config.toml")
	container, err := configure.NewContainer(ctx, os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	cfg := container.Config
	ctx = config.IntoContext(ctx, cfg)

	slog.Info("Configuration loaded successfully", "mode", cfg.Mode, "collection", cfg.Collection)
	slog.Info("Initialized OpenAI client", "model", cfg.Model)
	slog.Info("Connected to Qdrant", "endpoint", cfg.QdrantGRPC)

	var model openai.EmbeddingModel
	switch cfg.Model {
	case "text-embedding-3-small", "text-embedding-3-large":
		model = openai.EmbeddingModel(cfg.Model)
	default:
		log.Fatalf("unknown model: %s", cfg.Model)
	}

	switch cfg.Mode {
	case "ingest":
		slog.Info("Starting ingest mode", "html_dir", cfg.HTMLDir)
		ingestUC := ingest.New(
			container.IngestEmbeddingClient,
			container.IngestQdrantCollectionClient,
			container.IngestQdrantPointsClient,
			container.IngestHTMLParser,
			container.IngestTextChunker,
		)
		if err := ingestUC.Run(ctx, cfg.HTMLDir, model); err != nil {
			slog.Error("Ingest stopped with error", "error", err)
			os.Exit(1)
		}

		slog.Info("Ingest process completed successfully")
	case "search":
		if cfg.Query == "" {
			log.Fatal("-q is required in search mode")
		}
		slog.Info("Starting search mode", "query", cfg.Query, "top_k", cfg.TopK)
		searchUC := search.New(
			container.SearchEmbeddingClient,
			container.SearchQdrantPointsClient,
			container.SearchPromptBuilder,
		)
		slog.Info("Creating embedding for search query")
		hits, err := searchUC.Search(ctx, cfg.Query, cfg.TopK, model, cfg.Lang)
		if err != nil {
			log.Fatal(err)
		}
		slog.Info("Search completed", "results_count", len(hits))

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
