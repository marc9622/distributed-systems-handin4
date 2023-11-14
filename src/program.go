package main

import (
	"bufio"
	"os"

	node "github.com/marc9622/distributed-systems-handin4/src/node"
)

func main() {
    var ports = []uint{1111,2222,3333}

    for i, port := range ports {
        node.Spawn(port, i == 0, ports)
    }

    var scanner = bufio.NewScanner(os.Stdin)
    scanner.Scan()
}

