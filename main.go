package main

import (
	"encoding/binary"
	"fmt"
	"strconv"
)



func main() {
	bc := NewBlockchain()

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate(block)))
		fmt.Println()
		fmt.Println()
	}
}

func IntToHex(num int64) []byte {
    buff := make([]byte, 8) // Create a buffer of 8 bytes (since int64 is 8 bytes)
    binary.BigEndian.PutUint64(buff, uint64(num))
    return buff
}