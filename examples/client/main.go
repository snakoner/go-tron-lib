package main

import (
	"context"
	"log"
	"math/big"

	"github.com/snakoner/go-tron-lib"
)

func main() {
	client := tron.New("https://nile.trongrid.io")

	nowBlock, err := client.GetNowBlock(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("nowBlock: %s", nowBlock)

	trc20 := client.NewTRC20("TRPXG8YEMEaYE9dRs6fXvofFTiyMFE2mEg")
	tx, err := trc20.BuildTransferTx(context.Background(), "TZJ32TTQgjqcYWQf626xTWaZUT9iKLXxtS", "TDxyML69uweBFRfoEBEbGYQUE3XTWzUPe8", big.NewInt(1000000), 100000000)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tx: %s", tx)

	signedTx, err := tron.SignTransaction(tx, "30aa9a4134118c36f4d458004697ae1c3f97680ac5fadfd560d84c6482ad04c6")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("signedTx: %s", signedTx)

	resp, err := client.BroadcastTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("resp: %s", resp)

	status, err := client.WaitForStatusSuccess(context.Background(), resp.TxID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status: %s", status)
}
