package chunker

import (
	"fmt"
	"strings"

	"test-ragger/internal/models"
)

// TextChunker implements text chunking logic
type TextChunker struct{}

// New creates a new TextChunker instance
func New() *TextChunker {
	return &TextChunker{}
}

// ChunkText splits text into overlapping chunks with word boundary preservation
func (t *TextChunker) ChunkText(text string, size, overlap int) []models.ChunkInfo {
	var chunks []models.ChunkInfo

	position := 0
	chunkIndex := 0

	for position < len(text) {
		end := position + size
		if end > len(text) {
			end = len(text)
		}

		fragment := text[position:end]

		// Try not to cut words at the end of chunk
		// Only adjust if we're not at the end of text and we can find a good word boundary
		if end < len(text) {
			if lastSpace := strings.LastIndex(fragment, " "); lastSpace > int(float64(size)*0.6) {
				fragment = fragment[:lastSpace]
				end = position + lastSpace
			}
		}

		chunks = append(chunks, models.ChunkInfo{
			Text:    fragment,
			Start:   position,
			End:     end,
			ChunkID: fmt.Sprintf("ch_%d", chunkIndex),
		})

		chunkIndex++

		// Prevent infinite loop if end doesn't advance
		if end <= position {
			break
		}

		// Move position with overlap
		newPosition := end - overlap

		// Size of current chunk is less than overlap - exit
		if newPosition <= position {
			break
		}

		position = newPosition
	}

	return chunks
}

// ChunkTextWithCustomID splits text into chunks with custom chunk ID format
func (t *TextChunker) ChunkTextWithCustomID(text string, size, overlap int, idPrefix string) []models.ChunkInfo {
	chunks := t.ChunkText(text, size, overlap)

	// Update chunk IDs with custom prefix
	for i := range chunks {
		chunks[i].ChunkID = fmt.Sprintf("%s_%d", idPrefix, i)
	}

	return chunks
}

// EstimateChunkCount estimates how many chunks will be created
func (t *TextChunker) EstimateChunkCount(textLength, chunkSize, overlap int) int {
	if textLength <= chunkSize {
		return 1
	}

	effectiveStep := chunkSize - overlap
	if effectiveStep <= 0 {
		return textLength // fallback for invalid parameters
	}

	return ((textLength - chunkSize) / effectiveStep) + 1
}
