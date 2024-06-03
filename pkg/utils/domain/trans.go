package domain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	ID      string     `json:"id"`
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

type Trans struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type TxOutput struct {
	PubKey string `json:"pubkey"`
	Value  int    `json:"value"`
}

type TxInput struct {
	ID  string `json:"id"`
	Out int    `json:"out"`
	Sig string `json:"sig"`
}

func NewData(data Transaction) *Transaction {
	t := new(Transaction)

	t.ID = data.ID
	t.Inputs = data.Inputs
	t.Outputs = data.Outputs

	return t
}

func (b *Block) GiveData(data Transaction) {
	m := NewData(data)
	b.Transaction = append(b.Transaction, m)
}

func (t *Transaction) Hash() string {
	m, _ := json.Marshal(t)
	return fmt.Sprintf("%x", sha256.Sum256(m))
}
