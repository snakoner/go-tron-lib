package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/snakoner/go-tron-lib"
)

func main() {

	client := tron.New("https://nile.trongrid.io")

	trc20 := client.NewTRC20("TRPXG8YEMEaYE9dRs6fXvofFTiyMFE2mEg", true)
	pk := "30aa9a4134118c36f4d458004697ae1c3f97680ac5fadfd560d84c6482ad04c6"

	tx, err := trc20.BuildTransferTx(context.Background(), "TZJ32TTQgjqcYWQf626xTWaZUT9iKLXxtS", "TDxyML69uweBFRfoEBEbGYQUE3XTWzUPe8", big.NewInt(1000000), 100000000)
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

	for {
		status, err := client.GetTransactionStatus(context.Background(), resp.TxID)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("tx status: %s", status)
		time.Sleep(3 * time.Second)
	}
	return

	log.Printf("trc20 transfer txid: %s", resp.TxID)

	rawTx, err := client.BuildTransferTRXTx(context.Background(), "TZJ32TTQgjqcYWQf626xTWaZUT9iKLXxtS", "TDxyML69uweBFRfoEBEbGYQUE3XTWzUPe8", big.NewInt(1000000))
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err = tron.SignTransaction(rawTx, pk)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.BroadcastTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trx transfer txid: %s", resp.TxID)
}
