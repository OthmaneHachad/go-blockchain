package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)



type CLI struct {
	bc *Blockchain
}


func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)


	addBlockData := addBlockCmd.String("data", "", "Block data")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	switch os.Args[1] {

		case "addblock" : 
			err := addBlockCmd.Parse(os.Args[2:])
			if err != nil {
				panic(fmt.Sprintf("Error parsing addblock data flag: %s", err))
			}
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			if err != nil {
				panic(fmt.Sprintf("Error parsing printchain flags: %s", err))
			}
		case "createblockchain":
			err := createBlockchainCmd.Parse(os.Args[2:])
			if err != nil {
				panic(fmt.Sprintf("Error parsing createblockchain flags: %s", err))
			}
		default :
			cli.printUsage()
			os.Exit(1)
	}

	if (createBlockchainCmd.Parsed()) {
		if (*createBlockchainAddress == "") {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

}


func (cli *CLI) createBlockchain(address string) {
	bc, err := CreateBlockchain(address)
	if err != nil {
		panic(fmt.Sprintf("Error creating blockchain: %s", err))
	}
	bc.db.Close()
	fmt.Println("Done!")
}

func (cli *CLI) addBlock(data string) {
	//cli.bc.AddBlock(data)
	fmt.Println("Successfully added a Block!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Transactions)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

}



func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}