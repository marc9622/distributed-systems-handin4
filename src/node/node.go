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
}

func (node *Node) GiveToken(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
    //log.Printf("Received token at port %d\n", node.Port)
    select {
    case node.hasToken <- struct{}{}:
    default:
    }
    return &pb.Empty{}, nil
}

func Spawn(port uint, startsWithToken bool, allNodes []uint) {
    if len(allNodes) <= 1 {
        log.Fatalf("No nodes to connect to")
    }

    var node = Node {
        Port: port,
        hasToken: make(chan struct{}),
        allNodes: allNodes,
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
                select {
                case node.hasToken <- struct{}{}:
                //default:
                }
            }
        }()

        var opt = grpc.WithTransportCredentials(insecure.NewCredentials())

        var conn, connErr = grpc.Dial(fmt.Sprintf("localhost:%d", nextPort), opt)
        if connErr != nil {
            log.Panicf("Failed to dial server: %s", connErr)
        }
        defer conn.Close()

        var client = pb.NewTokenRingClient(conn)

        var ctx = context.Background()

        log.Printf("Running client at port %d\n", port);

        for {
            <- node.hasToken

            node.useToken()

            //log.Printf("Sending token from port %d to port %d\n", port, nextPort)
            var _, giveErr = client.GiveToken(ctx, &pb.Empty{})
            if giveErr != nil {
                log.Panicf("Failed to give token: %s", giveErr)
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

