package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Block struct {
	Header      *Header        `json:"header"`
	Transaction []*Transaction `json:"transaction"`
}

func (b *Block) Hash() string {
	m, _ := json.Marshal(b)
	return fmt.Sprintf("%x", sha256.Sum256(m))
}
