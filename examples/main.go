package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/snakoner/go-tron-lib"
)

func main() {
	client := tron.New("https://nile.trongrid.io")

	trc20 := client.NewTRC20("TRPXG8YEMEaYE9dRs6fXvofFTiyMFE2mEg", true)
	pk := "428c1224b892b2b2a619075a19cc6d74b3a235b5ffcac30cbccf9643658f61fd"

	tx, err := trc20.BuildTransferTx(context.Background(), "TGi14cNEkHAMzkT8dnd42SD82dkSTvHguL", "TDxyML69uweBFRfoEBEbGYQUE3XTWzUPe8", big.NewInt(1000000), 100000000)
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := tron.SignTransaction(tx, pk)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.BroadcastTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.TxID)
}
