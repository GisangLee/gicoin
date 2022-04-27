package blockchain

import (
	"time"

	"github.com/gisanglee/gicoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct {
	TxId  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type Utxo struct {
	TxId   string
	Amount int
	Index  int
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{TxId: "", Index: -1, Owner: "COINBASE"},
	}

	txOuts := []*TxOut{
		{Owner: address, Amount: minerReward},
	}

	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}

	tx.getId()

	return &tx
}

func makeTx(from string, to string, amount int) (*Tx, error) {

}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("gi", to, amount)

	if err != nil {
		return err
	}

	m.Txs = append(m.Txs, tx)

	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("gi")

	txs := m.Txs

	txs = append(txs, coinbase)

	m.Txs = nil

	return txs
}
