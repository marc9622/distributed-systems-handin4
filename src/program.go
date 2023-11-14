package main

import (
	"time"

	node    "github.com/marc9622/distributed-systems-handin4/src/node"
)

func main() {
    node.Spawn("8080", true,  []uint{0})
    node.Spawn("8080", false, []uint{0})
    node.Spawn("8080", false, []uint{0})

    time.Sleep(1 * time.Second)
}

