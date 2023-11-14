package critical

import (
	"fmt"
	"os"

	node "github.com/marc9622/distributed-systems-handin4/src/node"
)

var currentNode uint
var occupied bool

func EnterCriticalSection(node *node.Node) {
    if !occupied {
        currentNode = node.Id
        occupied = true
        fmt.Printf(">>> Node %d entered critical section.\n", node.Id)
    } else {
        fmt.Printf("!!! Node %d tried to enter critical section, but it is already occupied by node %d.\n", node.Id, currentNode)
        os.Exit(1)
    }
}

func LeaveCriticalSection(node *node.Node) {
    if occupied {
        occupied = false
        fmt.Printf("<<< Node %d left critical section.\n", node.Id)
    } else {
        fmt.Printf("!!! Node %d tried to leave critical section, but was not occupying it.\n", node.Id)
        os.Exit(1)
    }
}

