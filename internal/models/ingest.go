package models

// ChunkInfo represents a text chunk with metadata
type ChunkInfo struct {
	Text    string
	Start   int
	End     int
	ChunkID string
}
