# Glockchain

> Disclaimer: This is my first Go program so it obviously contains mistakes and/or bad practices. Don't hesitate to correct those.

## What is it ?

Glockchain is a simplified blockchain mechanism developped in Go following [this article](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46).

## Why

The article uses Python to implement the basic blockchain and I'm learing Golang right now so it was the perfect occasion to discover Golang features while getting a better understanding of what blockchain is.

## How to use it

This repo includes a *main.go* you can run with the following command `go run main.go`.
It starts 2 Node servers you can access through *http://localhost:5000* and *http://localhost:5001*.

## Available routes

* **/close** [GET]

Close the server

* **/chain** [GET]

Get the current list of blocks of the node

* **/mine** [GET]

Mine new blocks

* **/transactions/new** [POST]

Create new transactions by posting data as a JSON of the follow form:

```json
{
    "Sender": "ffff-fff-fff", //Identifier of the sender
    "Recipient": "aaaa-aaa-aaa", //Identifier of the recipient
    "Amount": 5, //Amount for the transaction
}
```

* **/nodes/register/<addressOfNeighbour>** [GET]

Registers a neighbour node to be used for consensus

* **/nodes/resolve** [GET]

Resolve the conflicts by inspecting every neighbour's chain
and replacing (if necessary) its own chain with the authoritative one.