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
    var port = flag.Uint("port", 0, "Port to listen on")
    var first = flag.Bool("first", false, "This node starts with the token")
    var debug = flag.Bool("debug", false, "Enable debug logging")
    flag.Parse()

    if *port == 0 {
        log.Fatal("Port must be specified with -port")
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

    node.Spawn(*port, *first, ports, *debug)

    var scanner = bufio.NewScanner(os.Stdin)
    scanner.Scan()
}

