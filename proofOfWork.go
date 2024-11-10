package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// we are implementing a fixed target algorithm
// target bits is the difficulty at which a block is mined
const targetBits = 24
/* 24 is an arbitrary number, our goal 
is to have a target that takes less than 256 bits in memory 
the greater 256 - targetbits is, the harder it is to mine a block */

const maxNonce = math.MaxInt64

type ProofOfWork struct {
	block *Block
	target *big.Int // refers to the requirement above
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits)) // left shift by 256 - targetBits

	return &ProofOfWork{block, target}
}


func (pow *ProofOfWork) prepareDataToHashing(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}


func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Minting the following block: \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareDataToHashing(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break 
		} else {
			nonce++
		}	
		
	}
	fmt.Printf("\n\n")

	return nonce, hash[:]

}


func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareDataToHashing(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}