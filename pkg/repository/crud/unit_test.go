package crud_test

import (
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/repository/crud"
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_LastKey(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

	db, err := levelDB.NewConnection(path)
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	repo := crud.NewRepository(db)

	ctx := context.Background()
	last, err := repo.GET(ctx, "lh")
	if err != nil {
		log.Fatalf("failed LastKey:" + err.Error())
	}

	log.Printf("last:%s", last)

}

func Test_GETChain(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

	db, err := levelDB.NewConnection(path)
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	repo := crud.NewRepository(db)

	chain, err := repo.GetChain(ctx)
	if err != nil {
		log.Fatalf("failed GetChain:" + err.Error())
	}

	js, err := json.MarshalIndent(chain, " ", " ")
	if err != nil {
		log.Fatalf("failed js:" + err.Error())
	}
	log.Printf("js:%s", string(js))
}

func Test_Wallet(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "pkg")]
	path := base + "database/"
	log.Printf("path:%s", path)

	db, err := levelDB.NewConnection(path)
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	repo := crud.NewRepository(db)

	// var w []string
	// w = append(w, "BANK A")
	// w = append(w, "BANK B")
	// w = append(w, "BANK C")
	// wallet := domain.DWallet{
	// 	Address: w,
	// }
	// js, err := json.Marshal(wallet)
	// if err != nil {
	// 	log.Fatalf("failed js:%s", err.Error())
	// }

	// err = repo.PUT(ctx, "wallet", string(js))
	// if err != nil {
	// 	log.Fatalf("failed PUT:%s", err.Error())
	// }
	// log.Println("succes create wallet")

	wallets, err := repo.GetWallet(ctx)
	if err != nil {
		log.Fatalf("failed GetWallet:" + err.Error())
	}
	for i, wallet := range wallets {
		log.Printf("wallet[%d]:%s", i+1, wallet)
	}

}
