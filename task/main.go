package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	blockchainSize = 5
)

type Block struct {
	blockID      int
	timestamp    int64
	previousHash string
}

type Blockchain struct {
	blocks []Block
}

func main() {
	// declare a blockchain of 5 blocks in sequence
	blocks := Blockchain{make([]Block, blockchainSize)}
	for i := 0; i < blockchainSize; i++ {
		if i == 0 {
			blocks.blocks[i] = createBlock(i+1, time.Now().UnixNano(), "0")
		} else {
			blocks.blocks[i] = createBlock(i+1, time.Now().UnixNano(), calculateHash(blocks.blocks[i-1]))
		}
	}
	printBlockchain(blocks)
}

func createBlock(blockID int, timestamp int64, prevBlockHash string) Block {
	return Block{blockID, timestamp, prevBlockHash}
}

func calculateHash(block Block) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d%d%s", block.blockID, block.timestamp, block.previousHash)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func printBlockchain(blocks Blockchain) {
	for i, block := range blocks.blocks {
		if i == 0 {
			fmt.Printf("Genesis Block:\n")
		} else {
			fmt.Printf("Block:\n")
		}
		fmt.Printf("Id: %d\n", block.blockID)
		fmt.Printf("Timestamp: %d\n", block.timestamp)
		fmt.Printf("Hash of the previous block:\n%s\n", block.previousHash)
		fmt.Printf("Hash of the block:\n%s\n", calculateHash(block))
		fmt.Println()
	}
}
