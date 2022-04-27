package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/gisanglee/gicoin/utils"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

const (
	walletFileName string = "gicoin.wallet"
)

func hasWalletFile() bool {
	_, err := os.Stat(walletFileName)

	return os.IsExist(err)
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

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(walletFileName)

	utils.HandleError(err)

	key, err := x509.ParseECPrivateKey(keyAsBytes)

	utils.HandleError(err)

	return key

}

func aFromK(key *ecdsa.PrivateKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()

	z := append(x, y...)

	return fmt.Sprintf("%x", z)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
	}

	w.Address = aFromK(w.privateKey)

	return w
}
