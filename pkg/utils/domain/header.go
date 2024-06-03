package domain

type Header struct {
	PrevHash string `json:"prev_hash"`
	Hash     string `json:"hash"`
	Nonce    uint   `json:"nonce"`
	Time     int64  `json:"time"`
}
