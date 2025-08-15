package models

// Hit represents a search result from vector database
type Hit struct {
	Score   float32
	Title   string
	Text    string
	DocID   string
	ChunkID string
	Path    string
}
