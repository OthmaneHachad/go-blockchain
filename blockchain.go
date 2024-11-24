package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func (blockchain *Blockchain) AddBlock(tx *Transaction) {
	var lastHash []byte
	var err error

	err = blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock([]*Transaction{tx}, lastHash)

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		serialized, serializedErr := newBlock.Serialize()
		if serializedErr != nil {
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

func CreateBlockchain(address string) (*Blockchain, error) {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := CreateGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		ser, err := genesis.Serialize()

		err = b.Put(genesis.Hash, ser)
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc, nil
}

func (bc *Blockchain) FindUnspentTransactions(address string) []*Transaction {
	// TODO: Implement
	return nil
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	iterator := &BlockchainIterator{blockchain.tip, blockchain.db}

	return iterator
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var SerializationErr error
	var err error

	err = iter.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(iter.currentHash)
		block, SerializationErr = DeserializeBlock(encodedBlock)

		if SerializationErr != nil {
			return SerializationErr
		}

		return nil

	})

	if err != nil {
		return nil
	}

	iter.currentHash = block.PrevBlockHash
	return block

}
