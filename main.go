package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func main() {
	//defer db.Close()
	//go explorer.Start(3000)
	//rest.Start(5000)
	//blockchain.Blockchain()
	//cli.Start()
	difficulty := 4
	target := strings.Repeat("0", difficulty)

	nonce := 1

	for {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte("hello"+fmt.Sprint(nonce))))
		fmt.Printf("Hash : %s\nTarget: %s\nNonce: %d\n\n", hash, target, nonce)
		if strings.HasPrefix(hash, target) {
			return
		}

		nonce++
	}
}
