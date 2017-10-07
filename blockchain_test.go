package blockchain

import (
	"fmt"
	"testing"
)

func TestNewBlockchain(t *testing.T) {

	b := NewBlockchain()

	if len(b.CurrentTransactions) != 0 {
		t.Error("CurrentTransactions should be empty")
	}

	if len(b.Chain) != 1 {
		t.Error("Chain should have one block")
	}

	block := b.LastBlock()

	if block.Index != 1 {
		t.Error("Index of default block should be 1")
	}

}

func TestNewTransaction(t *testing.T) {

	sender := "A"
	recipient := "B"
	amount := 5

	b := NewBlockchain()

	if len(b.CurrentTransactions) != 0 {
		t.Error("CurrentTransactions should be empty")
	}

	b.NewTransaction(sender, recipient, amount)

	if len(b.CurrentTransactions) != 1 {
		t.Error("CurrentTransactions should now have one item")
	}

	tr := &b.CurrentTransactions[0]

	if tr.Sender != sender {
		t.Error(fmt.Printf("Sender is '%s', expected '%s'", tr.Sender, sender))
	}

	if tr.Recipient != recipient {
		t.Error(fmt.Printf("Recipient is '%s', expected '%s'", tr.Recipient, recipient))
	}

	if tr.Amount != amount {
		t.Error(fmt.Printf("Amount is '%d', expected '%d'", tr.Amount, amount))
	}
}

func TestNewBlock(t *testing.T) {

	b := NewBlockchain()

	if len(b.Chain) != 1 {
		t.Error("Chain should have one block")
	}

	block := b.LastBlock()

	if block.Index != 1 {
		t.Error("Index of default block should be 1")
	}

	newBlock := b.NewBlock("", block.Proof)

	if newBlock.Index == block.Index {
		t.Error("Block indices should be different")
	}

	if newBlock.Timestamp <= block.Timestamp {
		t.Error("New block timestamp should be greater than the first block timestamp")
	}

	if newBlock.PreviousHash == block.PreviousHash {
		t.Error("New block previousHash should be different from the first block previousHash")
	}
}
