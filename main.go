package main

import "github.com/gisanglee/gicoin/blockchain"

func main() {
	//go explorer.Start(3000)
	//rest.Start(5000)
	//cli.Start()
	blockchain.Blockchain().AddBlock("First")
	blockchain.Blockchain().AddBlock("Second")
	blockchain.Blockchain().AddBlock("Third")
	blockchain.Blockchain().AddBlock("Fourth")
}
