package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type Block struct {
	Transactions 	[]*Transaction
	Timestamp 		int64
	PrevBlockHash 	[]byte
	Hash 			[]byte
	Nonce			int
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	newBlock := &Block{transactions, time.Now().Unix(), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	return newBlock
}

func CreateGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
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