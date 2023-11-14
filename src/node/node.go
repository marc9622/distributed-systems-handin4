package node

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	critical "github.com/marc9622/distributed-systems-handin4/src/critical"

	pb "github.com/marc9622/distributed-systems-handin4/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
    pb.UnimplementedTokenRingServer
    Port     uint
    hasToken chan struct{}
    allNodes []uint
    debug    bool
}

func (node *Node) GiveToken(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
    if node.debug {
        log.Printf("Received token at port %d\n", node.Port)
    }
    node.hasToken <- struct{}{}
    return &pb.Empty{}, nil
}

func Spawn(port uint, startsWithToken bool, allNodes []uint, debug bool) {
    if len(allNodes) <= 1 {
        log.Fatalf("No nodes to connect to")
    }

    var node = Node {
        Port: port,
        hasToken: make(chan struct{}),
        allNodes: allNodes,
        debug: debug,
    }

    // Setting up gRPC server
    go func() {

        var grpcServer = grpc.NewServer()

        pb.RegisterTokenRingServer(grpcServer, &node)

        var list, listErr = net.Listen("tcp", fmt.Sprintf("localhost:%d", port));
        if listErr != nil {
            log.Fatalf("Failed to listen: %v", listErr)
        }
        defer list.Close()

        log.Printf("Running server at port %d\n", port);

        var err = grpcServer.Serve(list)
        if err != nil {
            grpcServer.Stop()
            log.Panicf("Failed to serve: %s", err)
        }
    }()

    var nextPort = findNextPort(port, allNodes)

    // Setting up gRPC client
    go func() {
        go func() {
            if startsWithToken {
                node.hasToken <- struct{}{}
            }
        }()

        var opt = grpc.WithTransportCredentials(insecure.NewCredentials())

        var conn *grpc.ClientConn
        for {
            var connAttempt, connErr = grpc.Dial(fmt.Sprintf("localhost:%d", nextPort), opt)
            if connErr != nil {
                log.Panicf("Failed to dial server: %s trying again in 3 seconds", connErr)
                time.Sleep(3 * time.Second)
                continue
            }
            conn = connAttempt
            break
        }
        defer conn.Close()

        var client = pb.NewTokenRingClient(conn)

        var ctx = context.Background()

        log.Printf("Running client at port %d\n", port);

        for {
            <- node.hasToken

            if rand.Int() % 2 == 0 {
                node.useToken()
            } else {
                // We didn't need the token
            }

            if node.debug {
                log.Printf("Sending token from port %d to port %d\n", port, nextPort)
            }
            for {
                var _, giveErr = client.GiveToken(ctx, &pb.Empty{})
                if giveErr != nil {
                    log.Printf("Failed to give token trying again in 3 seconds: %s", giveErr)
                    time.Sleep(3 * time.Second)
                    continue
                }
                break
            }
        }
    }()
}

func findNextPort(port uint, allNodes []uint) uint {
    var nextPort uint = math.MaxUint
    var lowestPort uint = math.MaxUint

    for _, otherPort := range allNodes {
        if otherPort > port && otherPort < nextPort {
            nextPort = otherPort
        }
        if otherPort < lowestPort {
            lowestPort = otherPort
        }
    }
    if nextPort == math.MaxUint {
        if lowestPort == port {
            log.Fatalf("Could not find next port")
        } else {
            nextPort = lowestPort
        }
    }

    return nextPort
}

func (node *Node) useToken() {
    critical.EnterCriticalSection(node.Port)

    time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

    critical.LeaveCriticalSection(node.Port)
}

