package critical

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var fd *os.File
var writer *bufio.Writer

func EnterCriticalSection(id uint, file string) {
    var f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        log.Fatalf("!!! Node %d failed to open file %s: %v\n", id, file, err)
    }
    fd = f

    writer = bufio.NewWriter(f)
    writer.WriteString(fmt.Sprintf(">>> Node %d entered critical section.\n", id))
    writer.Flush()
}

func LeaveCriticalSection(id uint) {
    writer.WriteString(fmt.Sprintf("<<< Node %d left critical section.\n", id))
    writer.Flush()

    var err = fd.Close()
    if err != nil {
        log.Fatalf("!!! Node %d failed to close file: %v\n", id, err)
    }
}

