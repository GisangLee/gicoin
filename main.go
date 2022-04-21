package main

import (
	"fmt"

	"github.com/gisanglee/gicoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockchain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")
	result := blockchain.AllBlocks()
	fmt.Println(result)
}
