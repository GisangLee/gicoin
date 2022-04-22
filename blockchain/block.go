package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/gisanglee/gicoin/db"
	"github.com/gisanglee/gicoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"previous_hash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func createBlock(data string, prevHash string, height int) *Block {
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}

	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))

	block.persist()
	return &block
}

var ErrNotFound = errors.New("block not Found")

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}

	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
