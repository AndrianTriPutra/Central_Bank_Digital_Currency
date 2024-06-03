package blockchain_test

import (
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

func Test_FindUnspentTransactions(t *testing.T) {
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

	_, err = repoBC.FindUnspentTransactions(ctx, "Bank Central")
	if err != nil {
		log.Fatalf("failed FindUnspentTransactions:" + err.Error())
	}
	fmt.Println()

	// for i, tx := range unspentTxs {
	// 	js, err := json.MarshalIndent(tx, " ", " ")
	// 	if err != nil {
	// 		log.Fatalf("failed MarshalIndent:" + err.Error())
	// 	}
	// 	log.Printf("js[%d]:\n%s", i, string(js))
	// }

}

func Test_FindUTXO(t *testing.T) {
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

	unspentTxs, err := repoBC.FindUTXO(ctx, "Bank Central")
	if err != nil {
		log.Fatalf("failed Test_FindUTXO:" + err.Error())
	}
	fmt.Println()

	js, err := json.MarshalIndent(unspentTxs, " ", " ")
	if err != nil {
		log.Fatalf("failed MarshalIndent:" + err.Error())
	}
	log.Printf("js:\n%s", string(js))

}
