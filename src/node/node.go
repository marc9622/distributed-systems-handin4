package node

import (
	"log"
	"math/rand"
	"time"

	critical "github.com/marc9622/distributed-systems-handin4/src/critical"
)

type Node struct {
    Id       uint
    hasLock  bool // The token
    allNodes []uint
}

func Spawn(port string, hasLock bool, allNodes []uint) {
    var id = uint(rand.Uint32() % 1000)

    log.Printf("Spawning node %3d at port %s\n", id, port);

    go func() {
        var node = Node {
            Id: id,
            hasLock: false,
            allNodes: allNodes,
        }

        critical.EnterCriticalSection(node.Id)

        time.Sleep(100 * time.Millisecond)

        critical.LeaveCriticalSection(node.Id)
    }()
}

