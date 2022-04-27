package blockchain

import (
	"fmt"
	"sync"

	"github.com/gisanglee/gicoin/db"
	"github.com/gisanglee/gicoin/utils"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchain(b)
}

func persistBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

func Blocks(b *blockchain) []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)

		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}

	return blocks
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]

	actualTime := (newestBlock.TimeStamp / 60) - (lastRecalculatedBlock.TimeStamp / 60)
	expectedTime := difficultyInterval * blockInterval

	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1

	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1

	}

	return b.CurrentDifficulty

}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty

	} else if b.Height%difficultyInterval == 0 {

		return recalculateDifficulty(b)

	} else {
		return b.CurrentDifficulty
	}
}

func UTXOByAddress(address string, b *blockchain) []*Utxo {
	var uTxOuts []*Utxo

	sTxOuts := make(map[string]bool)

	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					sTxOuts[input.TxId] = true
				}
			}

			for index, output := range tx.TxOuts {
				_, ok := sTxOuts[tx.Id]

				if output.Owner == address {
					if !ok {
						uTxOut := &Utxo{TxId: tx.Id, Amount: output.Amount, Index: index}

						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTXOByAddress(address, b)

	var amount int

	for _, txOut := range txOuts {
		amount += txOut.Amount
	}

	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}

		// search checkpoin on the db
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			// restore b from bytes
			fmt.Println("Restoring....")
			b.restore(checkpoint)
		}
	})
	return b
}
