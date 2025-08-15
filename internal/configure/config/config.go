package config

import (
	"context"
	"errors"
	"flag"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	QdrantGRPC   string `toml:"qdrant_grpc"`
	Collection   string `toml:"collection"`
	EmbeddingDim int    `toml:"embedding_dim"`
	ChunkSize    int    `toml:"chunk_size"`
	ChunkOverlap int    `toml:"chunk_overlap"`

	// Model selection
	Model        string `toml:"model"`
	DefaultModel string `toml:"default_model"`

	// CLI/runtime options
	Mode    string `toml:"mode"` // ingest | search
	HTMLDir string `toml:"dir"`
	TopK    uint64 `toml:"k"`
	Query   string `toml:"q"`
	Lang    string `toml:"lang"`

	// Not serialized; resolved config path
	ConfigPath string `toml:"-"`
}

func Defaults() Config {
	return Config{
		QdrantGRPC:   "localhost:6334",
		Collection:   "docs",
		EmbeddingDim: 1536,
		ChunkSize:    1200,
		ChunkOverlap: 250,
		DefaultModel: "text-embedding-3-small",
		Model:        "",
		Mode:         "ingest",
		HTMLDir:      "./html",
		TopK:         5,
		Query:        "",
		Lang:         "",
	}
}

func LoadFromFile(path string) (Config, error) {
	if path == "" {
		return Defaults(), errors.New("empty config path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// no file – just return defaults
			return Defaults(), nil
		}
		return Defaults(), err
	}

	cfg := Defaults()
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Defaults(), err
	}
	if cfg.Model == "" && cfg.DefaultModel != "" { // back-compat
		cfg.Model = cfg.DefaultModel
	}
	cfg.ConfigPath = path
	return cfg, nil
}

// ResolveConfigPath extracts -config from args. Defaults to "config.toml".
func ResolveConfigPath(args []string) string {
	path := "config.toml"
	for i := 1; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-config=") {
			path = strings.TrimPrefix(a, "-config=")
			break
		}
		if a == "-config" && i+1 < len(args) {
			path = args[i+1]
			break
		}
	}
	return path
}

// Parse merges defaults, config file and CLI flags into a single Config.
// Priority: defaults < config.toml < CLI flags.
func Parse(args []string) (Config, error) {
	path := ResolveConfigPath(args)

	// load file (merged with defaults)
	base, err := LoadFromFile(path)
	if err != nil {
		return Defaults(), err
	}

	// define flags using base values
	cfgPathFlag := flag.String("config", path, "path to config file")
	_ = cfgPathFlag
	mode := flag.String("mode", base.Mode, "ingest | search")
	dir := flag.String("dir", base.HTMLDir, "папка с HTML (для ingest)")
	qdr := flag.String("qdrant", base.QdrantGRPC, "Qdrant gRPC addr")
	topK := flag.Uint64("k", base.TopK, "top-k (для search)")
	query := flag.String("q", base.Query, "запрос (для search)")
	modelName := flag.String("model", base.Model, "OpenAI embedding model: text-embedding-3-small|large")
	lang := flag.String("lang", base.Lang, "фильтр языка payload.lang (опц.)")
	flag.Parse()

	merged := base
	merged.Mode = *mode
	merged.HTMLDir = *dir
	merged.QdrantGRPC = *qdr
	merged.TopK = *topK
	merged.Query = *query
	merged.Model = *modelName
	merged.Lang = *lang
	merged.ConfigPath = path

	return merged, nil
}

// context helpers
type ctxKey int

const ctxKeyConfig ctxKey = iota

func IntoContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, ctxKeyConfig, cfg)
}

func FromContext(ctx context.Context) (Config, bool) {
	v := ctx.Value(ctxKeyConfig)
	if v == nil {
		return Defaults(), false
	}
	c, ok := v.(Config)
	return c, ok
}
