package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/gisanglee/gicoin/utils"
)

const (
	privateKey    string = "3077020101042054f798c5aaacd7bcd19fe8173ed7d35316646abbae7e36d9dc6ad9679a16f3f6a00a06082a8648ce3d030107a14403420004325a4fd736a70c3e419eb2960bda3311812de445237c16755ed8ffd738e8adff21a28340b10ff62be835f042ef94df2a468c94588a82ce3d21a47800daae2571"
	signature     string = "82c61adf7232d8f5ffa763b45e0afe25508fb3344a1c99506f887331f85e382aedacb019c07209e77d46a0fc9927555892dd0ca76361d7c312f93a6a71436039"
	hashedMessage string = "5bcf85f8bd23016d2314ac1c727a771c847f83dd97e817abab751b330af4cfb7"
)

func Start() {
	privateByte, err := hex.DecodeString(privateKey)

	utils.HandleError(err)

	restoredKey, err := x509.ParseECPrivateKey(privateByte)

	utils.HandleError(err)
	fmt.Println(restoredKey)

	signatureBytes, err := hex.DecodeString(signature)

	utils.HandleError(err)

	rBytes := signatureBytes[:len(signatureBytes)/2]
	sBytes := signatureBytes[len(signatureBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}

	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	fmt.Println(bigR, bigS)
}
