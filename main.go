package main

import (
	"github.com/gisanglee/gicoin/cli"
	"github.com/gisanglee/gicoin/db"
)

func main() {
	defer db.Close()
	//go explorer.Start(3000)
	//rest.Start(5000)
	//blockchain.Blockchain()
	cli.Start()
}
