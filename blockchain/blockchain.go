package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
}

type Block struct {
	Index        uint32
	Timestamp    int64
	Transactions []Transaction
	Proof        uint32
	PreviousHash string
}

type Blockchain struct {
	Chain               []Block
	CurrentTransactions []Transaction
}

func NewBlockchain() *Blockchain {
	b := &Blockchain{
		CurrentTransactions: make([]Transaction, 0),
		Chain:               make([]Block, 0),
	}

	b.NewBlock(100, "1")

	return b
}

func (b *Blockchain) NewBlock(proof uint32, previousHash string) *Block {

	var hashToUse string = previousHash

	if len(hashToUse) < 1 {
		hashToUse = hash(b.Chain[len(b.Chain)-1])
	}

	block := Block{
		Index:        uint32(len(b.Chain) + 1),
		Timestamp:    time.Now().UnixNano(),
		Transactions: b.CurrentTransactions,
		Proof:        proof,
		PreviousHash: hashToUse,
	}

	b.CurrentTransactions = make([]Transaction, 0)

	b.Chain = append(b.Chain, block)

	return &block
}

func (b *Blockchain) NewTransaction(sender, recipient string, amount int) uint32 {
	b.CurrentTransactions = append(b.CurrentTransactions, Transaction{
		Recipient: recipient,
		Sender:    sender,
		Amount:    amount,
	})

	return b.LastBlock().Index + uint32(1)
}

func (b *Blockchain) LastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

func proofOfWork(lastProof uint32) uint32 {
	var proof uint32 = uint32(0)

	for !validProof(lastProof, proof) {
		proof++
	}

	return proof
}

func hash(block Block) string {

	encoded, err := json.Marshal(block)

	if err != nil {
		panic(err)
	}

	return bytesToSha256Hex(encoded)
}

func validProof(lastProof, proof uint32) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	result := bytesToSha256Hex([]byte(guess))

	return strings.HasSuffix(result, "0000")
}

var (
	hash256 = sha256.New()
)

func bytesToSha256Hex(data []byte) string {
	hash256.Reset()
	hash256.Write(data)

	return fmt.Sprintf("%x", hash256.Sum(nil))
}
