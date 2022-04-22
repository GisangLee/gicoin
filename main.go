package main

import (
	"github.com/gisanglee/gicoin/blockchain"
	"github.com/gisanglee/gicoin/cli"
)

func main() {
	//go explorer.Start(3000)
	//rest.Start(5000)
	blockchain.Blockchain()
	cli.Start()
}
