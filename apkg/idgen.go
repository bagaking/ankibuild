package apkg

import (
	"crypto/rand"
	"crypto/sha1"
	"io"
	"strings"
	"sync"
	"time"
)

// Counter for appending to the Unix timestamp to create unique IDs. This is not thread-safe.
var (
	counter    int64 = 0
	mutexIDGen sync.Mutex
)

// genID generates a unique ID based on current time and an incrementing counter.
func genID() int {
	mutexIDGen.Lock()
	defer mutexIDGen.Unlock()

	// Reset the counter if it reaches 999
	if counter >= 999 {
		counter = 0
	}

	counter++ // Increment the counter
	return int(time.Now().Unix()*1000 + counter)
}

const (
	alphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	alphaLength  = len(alphaNumeric)
)

// generateGUID generates a globally unique identifier for Anki records.
// todo: 只是个猜测的实现
func genGUID() (string, error) {
	b := make([]byte, 10) // Anki uses a 10-char base91 sequence
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, byteVal := range b {
		sb.WriteByte(alphaNumeric[int(byteVal)%alphaLength])
	}
	return sb.String(), nil
}

// generateSortField generates an integer value to be used as the sort field (sfld).
// todo: 只是个猜测的实现
func generateSortField(front string) int {
	if front == "" {
		return 0
	}
	// This is a simplified and not precise way to create an integer value for the sort field:
	// taking the hash of front text and converting the first 4 bytes into an integer.
	h := sha1.New()
	h.Write([]byte(front))
	bs := h.Sum(nil)
	return int(bs[0])<<24 | int(bs[1])<<16 | int(bs[2])<<8 | int(bs[3])
}

// calculateChecksum generates a checksum for the provided text.
// todo: 只是个猜测的实现
// calculateChecksum generates a 64-bit checksum for the provided text.
func calculateChecksum(text string) int64 {
	h := sha1.New()
	h.Write([]byte(text))
	bs := h.Sum(nil)
	// Anki's checksum is the first three bytes, we'll use a 64-bit integer to store it.
	return int64(bs[0])<<16 | int64(bs[1])<<8 | int64(bs[2])
}
