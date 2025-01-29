package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Block represents a simple blockchain block
type Block struct {
	Index     int // Nonce
	Timestamp time.Time
	Data      string
	PrevHash  string
	Hash      string
}

// Validator represents a participating node in the PoS system
type Validator struct {
	ID           string
	Stake        int // Amount staked by the validator
	ResponseTime int // Simulated or measured response time in milliseconds
}

// Blockchain represents the state of the chain
type Blockchain struct {
	mu     sync.Mutex // Is used to ensure thread-safe operations on the blockchain state (Concurrency)
	Blocks []Block
}

// ProposeBlock represents a proposed block by the leader
func (bc *Blockchain) ProposeBlock(leader Validator, data string) Block {
	bc.mu.Lock()         // Acquires a lock before modifying the blockchain state.
	defer bc.mu.Unlock() // Ensures that the lock is released after the function's execution completes.

	prevHash := ""
	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].Hash
	}

	block := Block{
		Index:     len(bc.Blocks),
		Timestamp: time.Now(),
		Data:      data,
		PrevHash:  prevHash,
		Hash:      fmt.Sprintf("%x", rand.Int()),
	}
	fmt.Printf("Validator %s proposed block %d\n", leader.ID, block.Index)
	return block
}

// AppendBlock appends a valid block to the chain
func (bc *Blockchain) AppendBlock(block Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.Blocks = append(bc.Blocks, block)
	fmt.Printf("Block %d appended to the chain.\n", block.Index)
}

// SelectLeader selects a leader based on stake weight and response time
func SelectLeader(validators []Validator) Validator {
	// Sort validators by Stake (descending) and ResponseTime (descending)
	sort.SliceStable(validators, func(i, j int) bool {
		if validators[i].Stake == validators[j].Stake {
			return validators[i].ResponseTime > validators[j].ResponseTime // Higher response time wins
		}
		return validators[i].Stake > validators[j].Stake // Higher stake wins
	})

	// The first validator in the sorted list is the selected leader
	leader := validators[0]
	fmt.Printf("Selected Leader: %s (Stake: %d, Response Time: %dms)\n", leader.ID, leader.Stake, leader.ResponseTime)
	return leader
}

// ValidateBlock simulates block validation by other validators
func ValidateBlock(validators []Validator, block Block) bool {
	// Simulate all validators agreeing on the block
	fmt.Println("Validators validating block...")
	time.Sleep(1 * time.Second) // Simulate some processing time
	return true
}

func AssignStakes(validators []Validator) []Validator {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := range validators {
		validators[i].Stake = random.Intn(200) + 1       // Random stake between 1 and 200
		validators[i].ResponseTime = random.Intn(50) + 1 // Random Response time between 1 and 50
	}
	return validators
}

func main() {
	//TODO:01 Accept multiple validator nodes.
	validators := []Validator{
		{ID: "Validator1", Stake: 0, ResponseTime: 0},
		{ID: "Validator2", Stake: 0, ResponseTime: 0},
		{ID: "Validator3", Stake: 0, ResponseTime: 0},
	}

	//TODO:02 Assign stakes & responseTime dynamically
	validators = AssignStakes(validators)

	//Initialize blockchain
	chain := Blockchain{}
	// Simulate consensus rounds
	for i := 0; i < 5; i++ {
		fmt.Printf("\n--- Consensus Round %d ---\n", i+1)

		//TODO:03 Randomly select a leader based on stake amount and time.
		leader := SelectLeader(validators)

		//TODO:04 Propose a block and get it validated by the other nodes
		block := chain.ProposeBlock(leader, fmt.Sprintf("Block %d data", i))

		//TODO:05 Ensure the consensus is reached
		if ValidateBlock(validators, block) {
			//TODO:06 Update the state of the chain by appending the block.
			chain.AppendBlock(block)
		} else {
			fmt.Println("Block validation failed. Skipping...")
		}

		fmt.Printf("Current blockchain state: %+v\n", chain.Blocks)
	}
}
