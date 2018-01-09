# YaGo - Yet Another Go Coin
A proof-of-concept cryptocurrency written in Go.
Based on a simple blockchain as described at https://jeiwan.cc

## Installation
Pull down the package as normal
```bash
$ go get -u github.com/danmrichards/yagocoin
```

YaGo uses [dep](https://github.com/golang/dep) to manage it's dependencies. From the package directory install
them like so:
```bash
$ dep ensure -update
```

## Usage
```bash
Usage:
  yagocoin [command]

Available Commands:
  createblockchain Create a new blockchain
  createwallet     Generates a new key-pair and saves it into the wallet file
  getbalance       Get balance of adress
  help             Help about any command
  printchain       Print all the blocks of the blockchain
  send             Send an amount of coins from one address to another

Flags:
  -h, --help   help for yagocoin

Use "yagocoin [command] --help" for more information about a command.
```