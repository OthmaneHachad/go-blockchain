package main

import (
	"encoding/binary"
)



func main() {
	bc, err := NewBlockchain()
	if (err != nil) {
		panic(err)
	}

	bc.AddBlock("Send 1 BTC to Ivan")
	bc.AddBlock("Send 2 more BTC to Ivan")

}

func IntToHex(num int64) []byte {
    buff := make([]byte, 8) // Create a buffer of 8 bytes (since int64 is 8 bytes)
    binary.BigEndian.PutUint64(buff, uint64(num))
    return buff
}