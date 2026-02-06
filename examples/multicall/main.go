package main

import (
	"context"
	"fmt"
	"log"

	"github.com/snakoner/go-tron-lib"
)

const (
	multicallAddress = "TPP1ToFfmVXVTeWfJHAjmXsWnUGV8EkmnW"
	rpc              = "https://nile.trongrid.io"
	tokenAddress     = "TRPXG8YEMEaYE9dRs6fXvofFTiyMFE2mEg"
)

var addresses = []string{
	"TFbBApWL6TfyBhB8Tr322NeSUPePjnH4qe",
	"TZJ32TTQgjqcYWQf626xTWaZUT9iKLXxtS",
}

func main() {
	client := tron.New(rpc)
	multicall := client.NewMulticall(multicallAddress)

	balances, err := multicall.BalanceOf(context.Background(), tokenAddress, addresses)
	if err != nil {
		log.Fatal(err)
	}

	for i, balance := range balances {
		fmt.Printf("balance[%s]: %s\n", addresses[i], balance)
	}
}
