package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	blockchainSize = 5
	lowerTimeLimit = 5 * time.Second
	upperTimeLimit = 10 * time.Second
	minerNum       = 10
)

type Block struct {
	blockID      int
	creator      string
	timestamp    int64
	magicNumber  int32
	previousHash string
	thisHash     string
	message      string
	duration     time.Duration
}

type BlockChain struct {
	blocks     []Block
	zeroString string
}

type MiningResult struct {
	magicNumber int32
	hash        string
	miner       int
}

func main() {
	// initialization
	blockChain := createBlockChain(blockchainSize, "")

	for i := 0; i < blockchainSize; i++ {
		if i == 0 {
			blockChain.blocks[i] = *createBlock(i+1, time.Now().UnixNano(),
				blockChain.zeroString, "0")
		} else {
			blockChain.blocks[i] = *createBlock(i+1, time.Now().UnixNano(),
				blockChain.zeroString, blockChain.blocks[i-1].thisHash)
		}
		change := blockChain.adjustZero(blockChain.blocks[i].duration)
		blockChain.blocks[i].printBlock()
		printChange(change, blockChain.zeroString)
	}
}

func createBlock(blockID int, timestamp int64, zeroString string, prevBlockHash string) *Block {
	start := time.Now()
	block := &Block{blockID, "miner", timestamp, -1,
		prevBlockHash, "", "No messages", time.Duration(0)}
	if block.blockID != 1 {
		block.updateMessage()
	}
	block.calculateMagicNumber(zeroString)
	block.duration = time.Since(start)
	return block
}

func (block *Block) updateMessage() {
	fmt.Println("\nEnter a single message to send to the Blockchain:")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error when scanning message input: ", err)
	}
	block.message = strings.TrimSuffix(message, "\n")
}

func createBlockChain(length int, zeroString string) *BlockChain {
	blockChain := BlockChain{make([]Block, length), zeroString}
	return &blockChain
}

func (blockChain *BlockChain) adjustZero(duration time.Duration) int {
	flag := 0
	if duration < lowerTimeLimit {
		blockChain.zeroString += "0"
		flag = 1
	} else if duration > upperTimeLimit {
		blockChain.zeroString = strings.TrimSuffix(blockChain.zeroString, "0")
		flag = -1
	}
	return flag
}

func (block *Block) calculateHashWithMagic(magic int32) string {
	blockData := fmt.Sprintf("%d%d%d%s", block.blockID, block.timestamp, magic, block.previousHash)
	hash := sha256.New()
	hash.Write([]byte(blockData))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (block *Block) calculateHash() string {
	return block.calculateHashWithMagic(block.magicNumber)
}

func (block *Block) calculateMagicNumber(zeroString string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resultCh := make(chan MiningResult, 1)

	for i := 0; i < minerNum; i++ {
		go block.mining(ctx, i, zeroString, resultCh)
	}
	result := <-resultCh
	cancel()

	block.magicNumber = result.magicNumber
	block.thisHash = result.hash
	block.creator = fmt.Sprintf("miner%d", result.miner)
}

func (block *Block) mining(ctx context.Context, miner int, zeroString string, resultCh chan<- MiningResult) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		hushString := block.thisHash
		magic := block.magicNumber
		for !strings.HasPrefix(hushString, zeroString) || magic == -1 {
			magic = rand.Int31()
			hushString = block.calculateHashWithMagic(magic)
		}
		resultCh <- MiningResult{magic, hushString, miner}
	}
}

func (block *Block) printBlock() {

	if block.blockID == 1 {
		fmt.Printf("Genesis Block:\n")
	} else {
		fmt.Printf("\nBlock:\n")
		fmt.Printf("Created by %s\n", block.creator)
	}
	fmt.Printf("Id: %d\n", block.blockID)
	fmt.Printf("Timestamp: %d\n", block.timestamp)
	fmt.Printf("Magic number: %d\n", block.magicNumber)
	fmt.Printf("Hash of the previous block:\n%s\n", block.previousHash)
	fmt.Printf("Hash of the block:\n%s\n", block.thisHash)
	fmt.Printf("Block data:\n%s\n", block.message)
	fmt.Printf("Block was generating for %.0f seconds\n", block.duration.Seconds())
}

func printChange(change int, zeros string) {
	switch change {
	case 1:
		fmt.Printf("N was increased to %d\n", len(zeros))
	case 0:
		fmt.Print("N stays the same\n")
	case -1:
		fmt.Print("N was decreased by 1\n")
	}
}
