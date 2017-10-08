package blockchain

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Node struct {
	Identifier string
	Blockchain *Blockchain
	Server     *http.Server
	Nodes      map[string]struct{}
}

type ChainResponse struct {
	Chain  []Block
	Length int
}

func NewNode() *Node {

	n := &Node{
		Identifier: pseudoUuid(),
		Blockchain: NewBlockchain(),
		Nodes:      make(map[string]struct{}),
	}

	n.Server = newNodeServer(n)

	return n
}

func (n *Node) registerNode(address string) {
	n.Nodes[address] = struct{}{}
}

func (n *Node) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to node #%s#", n.Identifier)
}

func (n *Node) mineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Error")
		return
	}

	lb := n.Blockchain.LastBlock()

	lastProof := lb.Proof
	proof := proofOfWork(lastProof)

	n.Blockchain.NewTransaction("0", n.Identifier, 1)

	block := n.Blockchain.NewBlock(proof, "")

	encoded, err := json.Marshal(block)

	if err != nil {
		fmt.Println("Could not serialize mining response", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not serialize mining response"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(encoded)
}

func (n *Node) newTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Error")
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var tr Transaction

	err := decoder.Decode(&tr)

	if err != nil {
		panic(err)
		return
	}

	if len(tr.Sender) == 0 || len(tr.Recipient) == 0 || tr.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Missing fields in POST data")
		return
	}

	index := n.Blockchain.NewTransaction(tr.Sender, tr.Recipient, tr.Amount)

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, fmt.Sprintf("Transaction created for block with index %d", index))
}

func (n *Node) closeHandler(w http.ResponseWriter, r *http.Request) {
	n.Server.Close()
}

func (n *Node) chainHandler(w http.ResponseWriter, r *http.Request) {

	resp := &ChainResponse{
		Chain:  n.Blockchain.Chain,
		Length: len(n.Blockchain.Chain),
	}

	encoded, err := json.Marshal(resp)

	if err != nil {
		fmt.Println("Could not serialize chain response", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Could not serialize chain response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
}

func (n *Node) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.RequestURI(), "/")
	addr := parts[len(parts)-1]

	n.registerNode(addr)
}

func (n *Node) resolveConflictsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	replaced := n.resolveConflicts()

	if replaced {
		io.WriteString(w, "Our chain was replaced")
		return
	}

	io.WriteString(w, "Our chain is authoritative")
}

func (n *Node) StartListening() {
	fmt.Println("Listening on " + n.Server.Addr)
	n.Server.ListenAndServe()
}

func (n *Node) StopListening() {
	fmt.Println("Closing " + n.Server.Addr)
	n.Server.Close()
}

func (n *Node) resolveConflicts() bool {
	var newChain []Block
	neighbours := n.Nodes

	maxLen := len(n.Blockchain.Chain)

	for neighbour := range neighbours {
		chainRes, err := getChainResponse(neighbour)

		if err != nil {
			fmt.Println("Error during conflict resolution:", err)
		}

		if chainRes.Length > maxLen && validChain(chainRes.Chain) {
			maxLen = chainRes.Length
			newChain = chainRes.Chain
		}

	}

	if newChain != nil {
		n.Blockchain.Chain = newChain
		return true
	}

	return false
}

func getChainResponse(addr string) (*ChainResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/chain", addr))

	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var chain ChainResponse

	err = decoder.Decode(&chain)

	if err != nil {
		return nil, err
	}

	return &chain, nil
}

func pseudoUuid() string {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)

	if err != nil {
		panic(err)
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", bytes[:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])

	return uuid
}

var (
	currentPort = 4999
)

func newNodeServer(n *Node) *http.Server {
	currentPort += 1

	sm := http.NewServeMux()

	sm.HandleFunc("/", n.homeHandler)
	sm.HandleFunc("/close", n.closeHandler)
	sm.HandleFunc("/nodes/register/", n.registerHandler)
	sm.HandleFunc("/nodes/resolve/", n.resolveConflictsHandler)
	sm.HandleFunc("/mine", n.mineHandler)
	sm.HandleFunc("/transactions/new", n.newTransactionHandler)
	sm.HandleFunc("/chain", n.chainHandler)

	return &http.Server{
		Addr:           fmt.Sprintf(":%d", currentPort),
		Handler:        sm,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
