package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	blockchainSize = 5
)

type Block struct {
	blockID      int
	timestamp    int64
	magicNumber  int32
	previousHash string
	duration     time.Duration
}

func main() {
	// declare a blockchain of 5 blocks in sequence
	var zeros int
	fmt.Print("Enter how many zeros the hash must start with:\n")
	_, err := fmt.Scan(&zeros)
	if err != nil {
		fmt.Println("Error when scanning input: ", err)
	}

	zeroString := strings.Repeat("0", zeros)
	blocks := make([]Block, blockchainSize)
	for i := 0; i < blockchainSize; i++ {
		if i == 0 {
			blocks[i] = createBlock(i+1, time.Now().UnixNano(), zeroString, "0")
		} else {
			blocks[i] = createBlock(i+1, time.Now().UnixNano(), zeroString, calculateHash(blocks[i-1]))
		}
	}

	printBlockchain(blocks)
}

func createBlock(blockID int, timestamp int64, zeroString string, prevBlockHash string) Block {
	start := time.Now()
	block := Block{blockID, timestamp, -1, prevBlockHash,
		time.Duration(0)}
	calculateMagicNumber(&block, zeroString)
	block.duration = time.Since(start)
	return block
}

func calculateHash(block Block) string {
	blockData := fmt.Sprintf("%d%d%d%s", block.blockID, block.timestamp, block.magicNumber, block.previousHash)
	hash := sha256.New()
	hash.Write([]byte(blockData))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func calculateMagicNumber(block *Block, zeroString string) {
	hushString := calculateHash(*block)
	for !strings.HasPrefix(hushString, zeroString) || block.magicNumber == -1 {
		block.magicNumber = rand.Int31()
		hushString = calculateHash(*block)
	}
}

func printBlockchain(blocks []Block) {
	for i, block := range blocks {
		if i == 0 {
			fmt.Printf("\nGenesis Block:\n")
		} else {
			fmt.Printf("Block:\n")
		}
		fmt.Printf("Id: %d\n", block.blockID)
		fmt.Printf("Timestamp: %d\n", block.timestamp)
		fmt.Printf("Magic number: %d\n", block.magicNumber)
		fmt.Printf("Hash of the previous block:\n%s\n", block.previousHash)
		fmt.Printf("Hash of the block:\n%s\n", calculateHash(block))
		fmt.Printf("Block was generating for %.0f seconds\n", block.duration.Seconds())
		fmt.Println()
	}
}
