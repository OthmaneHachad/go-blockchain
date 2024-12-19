package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 10; // constant for now


type Transaction struct {
	ID []byte
	ValueIn []TXInput // an input must reference an output
	ValueOut []TXOutput // an output may not reference a future input
}


type TXOutput struct {
	Value int // this output value is indivisble
	ScriptPublicKey string // this is where the coins are actually stored. Since we don't have addresses yet, will be arbitrary string
}

type TXInput struct {
	TxId []byte
	ValueOut int
	ScriptSignature string // also known as the Input Data

}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSignature == unlockingData // will be improved later on after implementing addresses
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPublicKey == unlockingData // will be improved later on after implementing addresses
}


// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.ValueIn) == 1 && len(tx.ValueIn[0].TxId) == 0 && tx.ValueIn[0].ValueOut == -1
}


// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}


func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{reward, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}



