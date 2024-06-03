package utransaction_test

import (
	"atp/cbdc/app/usecase/utransaction"
	"atp/cbdc/pkg/adapter/levelDB"
	"atp/cbdc/pkg/repository/blockchain"
	"atp/cbdc/pkg/repository/crud"
	"atp/cbdc/pkg/repository/pow"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_Chain(t *testing.T) {
	ctx := context.Background()

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	log.Printf("basepath:%s", basepath)
	base := basepath[0:strings.Index(basepath, "app")]
	path := base + "database/"
	log.Printf("path:%s", path)

	db, err := levelDB.NewConnection(path)
	if err != nil {
		log.Fatalf("failed connect to database:" + err.Error())
	}

	setting := pow.Setting{
		Difficult: 13,
	}
	genesis := blockchain.Genesis{
		Address: "Bank Central",
		Balance: 100,
	}
	repoPoW := pow.NewRepository(setting)
	repoCrud := crud.NewRepository(db)
	repoBC := blockchain.NewRepository(genesis, repoCrud, repoPoW)

	ucase := utransaction.NewUsecase(repoBC, repoCrud)

	chain, err := ucase.GetChain(ctx)
	if err != nil {
		log.Fatalf("failed Chain:" + err.Error())
	}
	js, err := json.MarshalIndent(chain, " ", " ")
	if err != nil {
		log.Fatalf("failed js:" + err.Error())
	}
	log.Printf("js:%s", string(js))
	fmt.Println()

}
