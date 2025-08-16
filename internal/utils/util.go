package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func MustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("env %s is required", key)
	}
	return value
}

func Sha1Hex(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func Sha1Hash(s string) uint32 {
	h := sha1.Sum([]byte(s))
	// Convert first 4 bytes to uint32
	return uint32(h[0])<<24 | uint32(h[1])<<16 | uint32(h[2])<<8 | uint32(h[3])
}

func Snippet(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func BoolPtr(b bool) *bool { return &b }

func Uint64Ptr(u uint64) *uint64 { return &u }

// CleanUTF8 removes invalid UTF-8 characters from a string
func CleanUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Replace invalid UTF-8 sequences with replacement character
	var cleaned strings.Builder
	cleaned.Grow(len(s))

	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			// Skip invalid byte
			s = s[1:]
		} else {
			cleaned.WriteRune(r)
			s = s[size:]
		}
	}

	return cleaned.String()
}
