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

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty

	} else if b.Height%difficultyInterval == 0 {

		return b.recalculateDifficulty()

	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) UTXOByAddress(address string) []*Utxo {
	var uTxOuts []*Utxo

	sTxOuts := make(map[string]bool)

	for _, block := range b.Blocks() {
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
						uTxOuts = append(uTxOuts, &Utxo{TxId: tx.Id, Amount: output.Amount, Index: index})
					}
				}
			}
		}
	}

	return uTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.UTXOByAddress(address)

	var amount int

	for _, txOut := range txOuts {
		amount += txOut.Amount
	}

	return amount
}

func Blockchain() *blockchain {
	if b == nil {
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
	}
	fmt.Println(b.NewestHash)
	return b
}
