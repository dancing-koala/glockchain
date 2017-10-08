package main

import (
	"github.com/dancing-koala/glockchain/blockchain"
)

func main() {
	n := blockchain.NewNode()

	go blockchain.NewNode().StartListening()

	n.StartListening()
}
