package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Data 			[]byte
	Timestamp 		int64
	PrevBlockHash 	[]byte
	Hash 			[]byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	newBlock := &Block{[]byte(data), time.Now().Unix(), prevBlockHash, []byte{}}
	newBlock.SetHash()
	return newBlock
}

func (block *Block) SetHash() {
	// hash together all block fields
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	dataToHash := bytes.Join([][]byte{block.Data, block.PrevBlockHash, timestamp}, []byte{})
	hash := sha256.Sum256(dataToHash)

	block.Hash = hash[:] // all values in hash form index 0, length - 1

}