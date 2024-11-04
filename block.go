package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

type Block struct {
	Data 			[]byte
	Timestamp 		int64
	PrevBlockHash 	[]byte
	Hash 			[]byte
	Nonce			int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	newBlock := &Block{[]byte(data), time.Now().Unix(), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	return newBlock
}

func (block *Block) SetHash() {
	// hash together all block fields
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	dataToHash := bytes.Join([][]byte{block.Data, block.PrevBlockHash, timestamp}, []byte{})
	hash := sha256.Sum256(dataToHash)

	block.Hash = hash[:] // all values in hash form index 0, length - 1

}


func (block *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if (err != nil) {
		return nil, err
	}

	return result.Bytes(), nil
}

func DeserializeBlock(sb []byte) (*Block, error) {
	var block Block


	decoder := gob.NewDecoder(bytes.NewReader(sb))
	err := decoder.Decode(&block)
	if (err != nil) {
		return nil, err
	}

	return &block, nil
}