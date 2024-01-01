package hash

import (
	"bytes"
	"crypto/rand"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
)

type params struct {
	saltLength  uint32
	time        uint32
	memory      uint32
	threadCount uint8
	keyLength   uint32
}

// mostly values from golang documentation
var p = &params{
	saltLength:  16,
	time:        3,
	memory:      64 * 1024,
	threadCount: 4,
	keyLength:   32,
}

func GenHash(c *gin.Context, password string) (hash []byte, salt []byte) {
	salt = genSalt(p.saltLength)

	hash = argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threadCount, p.keyLength)

	return hash, salt
}

func CompareHash(c *gin.Context, password string, hash []byte, salt []byte) (match bool) {
	newHash := argon2.IDKey([]byte(password), salt, p.time, p.memory, p.threadCount, p.keyLength)

	return bytes.Equal(newHash, hash)
}

func genSalt(n uint32) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
