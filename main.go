package main

import (
	"encoding/binary"
)



func main() {
	cli := CLI{}
	cli.Run()
}

func IntToHex(num int64) []byte {
    buff := make([]byte, 8) // Create a buffer of 8 bytes (since int64 is 8 bytes)
    binary.BigEndian.PutUint64(buff, uint64(num))
    return buff
}