package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	node "github.com/marc9622/distributed-systems-handin4/src/node"
)

func main() {
    var port = flag.Uint("port", 0, "The port of this node")
    var first = flag.Bool("first", false, "Whether this node starts with the token")
    var debug = flag.Bool("debug", false, "Enable debug logging")
    var file = flag.String("file", "", "The file which the critical section writes to")
    flag.Parse()

    if *port == 0 {
        log.Fatal("Port must be specified with -port")
    }
    if *file == "" {
        log.Fatal("File must be specified with -file")
    }

    var ports = []uint{}
    {
        var portStrs = flag.Args();
        for _, portStr := range portStrs {
            var port uint

            var _, err = fmt.Sscanf(portStr, "%d", &port)
            if err != nil {
                log.Fatalf("Failed to parse port %s: %v", portStr, err)
            }
            ports = append(ports, port)
        }
    }

    node.Spawn(*port, *first, ports, *debug, *file)

    var scanner = bufio.NewScanner(os.Stdin)
    scanner.Scan()
}

