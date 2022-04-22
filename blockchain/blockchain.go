package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"previous_hash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)

	if totalBlocks == 0 {
		return ""
	}

	return GetBlockchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{Data: data, Hash: "", PrevHash: getLastHash(), Height: len(GetBlockchain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})
	}

	return b
}

func AllBlocks() []*Block {
	return GetBlockchain().blocks
}

var ErrNotFound = errors.New("Block Not Found")

func (b *blockchain) GetBlock(height int) (*Block, error) {

	if len(b.blocks) < height {
		return nil, ErrNotFound
	}

	return b.blocks[height-1], nil
}
