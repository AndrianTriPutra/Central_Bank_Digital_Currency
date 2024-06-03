package domain

type Balance struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type Wallet struct {
	Index   int    `json:"index"`
	Address string `json:"address"`
}

type DWallet struct {
	Address []string `json:"address"`
}
