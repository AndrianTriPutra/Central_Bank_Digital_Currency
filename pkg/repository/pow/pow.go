package pow

import (
	"atp/cbdc/pkg/utils/domain"
	"crypto/sha256"
	"encoding/json"
	"log"
	"math"
	"math/big"
	"time"
)

func (r repository) NewProof(b *domain.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-r.setting.Difficult))
	//log.Printf("target:%v", target)
	pow := &ProofOfWork{b, target}
	return pow
}

func (r repository) Run(pOw *ProofOfWork) (uint, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0
	start := time.Now()
	elapsed := 1 * time.Nanosecond
	for nonce < math.MaxInt64 {
		data := r.initData(uint(nonce), pOw)
		hash = sha256.Sum256(data)

		//fmt.Printf("\rRun %d:%x", nonce, hash)
		elapsed = time.Since(start)
		//log.Printf("Run PoW:%v", elapsed)
		//time.Sleep(1 * time.Second)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pOw.Target) == -1 {
			break
		} else {
			nonce++
		}

	}
	log.Printf("PoW Process:%v", elapsed)

	return uint(nonce), hash[:]
}

func (r repository) Validate(pOw *ProofOfWork, nonce uint) bool {
	var intHash big.Int

	dataX := r.initData(nonce, pOw)
	hash := sha256.Sum256(dataX)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pOw.Target) == -1
}

func (r repository) initData(nonce uint, pOw *ProofOfWork) []byte {
	guesBlock := &domain.Block{
		Header: &domain.Header{
			Nonce:    nonce,
			Time:     pOw.Block.Header.Time,
			PrevHash: pOw.Block.Header.PrevHash,
		},
		Transaction: pOw.Block.Transaction,
	}

	js, _ := json.Marshal(guesBlock)
	//log.Printf("js:%s", string(js))

	return js
}
