package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/gisanglee/gicoin/utils"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	utils.HandleError(err)

	message := utils.Hash("I Love You")

	byteMessage, err := hex.DecodeString(message)

	utils.HandleError(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, byteMessage)

	utils.HandleError(err)

	fmt.Printf("R:%d\nS:%d\n\n", r, s)

}
