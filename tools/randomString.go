package tools

import (
	"math/rand"
	"time"
)

// charset is default charset to create random string
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// StringWithCharset returns random string provided length based on provided charset
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// String generates random length string based on alphanumeric chars
func String(length int) string {
	return StringWithCharset(length, charset)
}
