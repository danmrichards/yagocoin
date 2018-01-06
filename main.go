package main

import (
	"log"

	"github.com/danmrichards/yagocoin/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
