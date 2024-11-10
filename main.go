package main

import (
	"encoding/binary"
	"fmt"
)



func main() {
	bc, err := NewBlockchain()
	if err != nil {
		panic(fmt.Sprintf("Error when creating blockchain: %s", err))
	}
	
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}

func IntToHex(num int64) []byte {
    buff := make([]byte, 8) // Create a buffer of 8 bytes (since int64 is 8 bytes)
    binary.BigEndian.PutUint64(buff, uint64(num))
    return buff
}