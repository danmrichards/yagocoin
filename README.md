# YaGo - Yet Another Go Coin
A proof-of-concept cryptocurrency written in Go.
Based on a simple blockchain as described at https://jeiwan.cc

## Installation
Pull down the package as normal
```bash
$ go get -u github.com/danmrichards/yagocoin
```

YaGo uses dep to manage it's dependencies. From the package directory install
them like so:
```bash
$ dep ensure -update
```

## Usage
```bash
Usage:
  yagocoin [command]

Available Commands:
  addblock    Add a block to the crypto
  help        Help about any command
  printchain  Print all the blocks of the blockchain

Flags:
  -h, --help   help for yagocoin

Use "yagocoin [command] --help" for more information about a command.

```