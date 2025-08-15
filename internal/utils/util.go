package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"os"
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

func Snippet(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

func BoolPtr(b bool) *bool { return &b }

func Uint64Ptr(u uint64) *uint64 { return &u }
