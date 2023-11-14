package main

import (
	"fmt"
	"math/rand"
	"time"

	critical "github.com/marc9622/distributed-systems-handin4/src/critical"
	node "github.com/marc9622/distributed-systems-handin4/src/node"
)

func main() {
    fmt.Println("Hello, World!");

    for i := 1; i <= 3; i++ {
        go func() {
            var node = node.Node{ Id: uint(rand.Uint32()) }

            critical.EnterCriticalSection(&node)

            critical.LeaveCriticalSection(&node)
        }()
    }

    time.Sleep(1 * time.Second)
}

