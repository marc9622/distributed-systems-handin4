package critical

import (
	"log"
	"os"
)

var currentNode uint
var occupied bool

func EnterCriticalSection(id uint) {
    if !occupied {
        currentNode = id
        occupied = true
        log.Printf(">>> Node %d entered critical section.\n", id)
    } else {
        log.Printf("!!! Node %d tried to enter critical section, but it is already occupied by node %d.\n", id, currentNode)
        os.Exit(1)
    }
}

func LeaveCriticalSection(id uint) {
    if occupied {
        occupied = false
        log.Printf("<<< Node %d left critical section.\n", id)
    } else {
        log.Printf("!!! Node %d tried to leave critical section, but was not occupying it.\n", id)
        os.Exit(1)
    }
}

