package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db *bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (blockchain *Blockchain) AddBlock(data string) {
	var lastHash []byte 
	var err error

	err = blockchain.db.View(func (tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(data, lastHash)

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		serialized, serializedErr := newBlock.Serialize()
		if (serializedErr != nil) {
			return serializedErr
		}

		// save serialized represntation of block in DB
		if err := bucket.Put(newBlock.Hash, serialized); err != nil {
			return err
		}

		// update the last block hash
		if err := bucket.Put([]byte("l"), newBlock.Hash); err != nil {
			return err
		}


		blockchain.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("Error updating internal DB: %s", err))
	}
}

func CreateGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() (*Blockchain, error) {
	var tip []byte // last block added to the BC
	db, err := bolt.Open(dbFile, 0600, nil)

	if (err != nil) {
		return nil, err
	}

	err = db.Update(func (tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		if bucket == nil {
			// means we have no block in the bucket
			genesis := CreateGenesisBlock()
			bucket, err := tx.CreateBucket([]byte(blocksBucket))
			if (err != nil) {
				return err
			}

			serialized, serializedErr := genesis.Serialize()
			if (serializedErr != nil) {
				return serializedErr
			}

			err = bucket.Put(genesis.Hash, serialized)
			err = bucket.Put([]byte("l"), genesis.Hash)

			tip = genesis.Hash

		} else {
			tip = bucket.Get([]byte("l"))
			// 'l' -> 4-byte file number: the last block file number used
		}

		return nil
	})

	bc := &Blockchain{tip, db}

	return bc, nil
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	iterator := &BlockchainIterator{blockchain.tip, blockchain.db}

	return iterator
}



type BlockchainIterator struct {
	currentHash []byte
	db *bolt.DB
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var SerializationErr error
	var err error

	err = iter.db.View(func (tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(iter.currentHash)
		block, SerializationErr = DeserializeBlock(encodedBlock)

		if (SerializationErr != nil) {
			return SerializationErr
		}

		return nil

	})

	if (err != nil) {
		return nil
	}

	iter.currentHash = block.PrevBlockHash
	return block


}

