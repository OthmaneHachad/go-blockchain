package main

import (
	"encoding/hex"
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

// NewBlockchain creates a new Blockchain (reference) with genesis Block
func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}


// creates a new blockchain DB
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

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	accumulated, validOutputs := bc.FindSpendableOutputs(from, amount)
	// note that TX Outputs are indivisible

	if accumulated < amount {
		log.Panic("Error: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if (err != nil) {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if accumulated > amount {
		outputs = append(outputs, TXOutput{accumulated - amount, from}) // remaining change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serialized, err := newBlock.Serialize()
		if (err != nil) {
			log.Panic(err)
		}

		err = b.Put(newBlock.Hash, serialized)
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// counting the unspent tx ouptuts
		Outputs:
			for outIdx, out := range tx.ValueOut {
				if spentTXOs[txID] != nil {
					// we search for tx outs that were locked by address
					// we skip those that were referenced in another tx input, 
						// as this means they were moved to other outputs
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// we store those that were never referenced elsewhere
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// we gather all inputs taht could unlock outputs locked with given address
			// coinbase tx cannot unlock outputs
			if tx.IsCoinbase() == false {
				for _, in := range tx.ValueIn {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.TxId)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.ValueOut)
					}
				}
			}

		}

		if (len(block.PrevBlockHash) == 0) {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int) // indices of the tx outputs accumulated grouped by Transactions
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.ValueOut {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.ValueOut {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
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
