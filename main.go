package main

import (
	"github.com/gisanglee/gicoin/explorer"
	"github.com/gisanglee/gicoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(5000)
}
