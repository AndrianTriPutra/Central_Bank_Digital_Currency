package blockchain

import (
	"atp/cbdc/pkg/utils/domain"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func (bc *blockchain) AddBlock(ctx context.Context, data []*domain.Transaction) error {
	b := new(domain.Block)
	b.Transaction = data

	prev, err := bc.repoCrud.GET(ctx, "lh")
	if err != nil {
		buff := err.Error()
		if !strings.Contains(buff, "not found") {
			//log.Fatalf("FAILED NewBlock GET:" + err.Error())
			return err
		}
	}

	b.Header = &domain.Header{
		PrevHash: prev,
		Time:     time.Now().UnixNano(),
	}

	pow := bc.repoPow.NewProof(b)
	nonce, hash := bc.repoPow.Run(pow)
	b.Header.Nonce = nonce
	b.Header.Hash = string(hash)
	valid := bc.repoPow.Validate(pow, nonce)

	nowHash := hex.EncodeToString(hash)
	log.Printf("valid ? %v", valid)

	b.Header.Hash = nowHash

	create := false
	if len(prev) == 0 {
		log.Println(" ==== Create GENESIS ====")

		// create wallet genesis
		pubkey := ""
		for _, tx := range b.Transaction {
			for _, txo := range tx.Outputs {
				pubkey = txo.PubKey
			}
		}
		x := pubkey
		w := domain.DWallet{}
		w.Address = append(w.Address, x)

		js, err := json.Marshal(w)
		if err != nil {
			return err
		}

		// put wallet genesis to db
		err = bc.repoCrud.PUT(ctx, "wallet", string(js))
		if err != nil {
			return err
		}

		b.Header.PrevHash = fmt.Sprintf("%x", [32]byte{})
		create = true
	} else {
		for _, tx := range b.Transaction {
			for _, txi := range tx.Inputs {
				if txi.ID != "0000000000000000000000000000000000000000000000000000000000000000" {
					create = true
					break
				}
			}
		}

	}

	if create {
		//put data
		js, _ := json.Marshal(b)
		err = bc.repoCrud.PUT(ctx, nowHash, string(js))
		if err != nil {
			return err
		}

		//put key
		err = bc.repoCrud.PUT(ctx, "lh", nowHash)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *blockchain) CreateBlock(ctx context.Context, a *domain.Block) error {
	err := bc.AddBlock(ctx, a.Transaction)
	if err != nil {
		return err
	}

	a.Transaction = []*domain.Transaction{}
	return err
}
