package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"

	"github.com/gisanglee/gicoin/utils"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet

const (
	walletFileName string = "gicoin.wallet"
)

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)

	return os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	utils.HandleError(err)

	return privateKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)

	err = os.WriteFile(walletFileName, bytes, 0644)
	utils.HandleError(err)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
	}

	return w
}
