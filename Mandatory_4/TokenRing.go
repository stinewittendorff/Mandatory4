package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	proto "TokenRing/proto"

	"google.golang.org/grpc"
)

type Node struct {
	proto.UnimplementedTokenringServer
	id              int32
	HasToken        bool
	RequestAccess   bool
	NextNodeAddress string
	mu              sync.Mutex
}

func (n *Node) Send(ctx context.Context, in *proto.Token) (*proto.Empty, error) {
	fmt.Printf("")
	fmt.Printf("Node received token %d\n", n.id)
	n.HasToken = true
	if n.RequestAccess == true {
	fmt.Println()

	n.EnterCriticalSection()
	}

	fmt.Printf("Node %d going to %s\n", n.id, n.NextNodeAddress)
	time.Sleep(time.Second * 5)

	n.passTokenToNext()
	

	return &proto.Empty{}, nil
}

func (n *Node) passTokenToNext() {
	n.HasToken = false

	conn, err := grpc.Dial(n.NextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}

	client := proto.NewTokenringClient(conn)
	token := &proto.Token{
		Id:      n.id,
		Message: "Token passing",
	}

	defer conn.Close()

	fmt.Printf("Node %d sending token to %s\n at time %s", n.id, n.NextNodeAddress, time.Now().Format("15:04:05"))
	_, err = client.Send(context.Background(), token)
	if err != nil {
		log.Fatalf("Could not send token: %s", err)
	}
}

func (n *Node) EnterCriticalSection() {
    fmt.Printf("Node %d is entering critical section\n", n.id)
    time.Sleep(time.Second * 5)
	n.RequestAccess = false
    fmt.Printf("Node %d is leaving critical section\n", n.id)
}

func (n *Node) RequestCriticalSection() {
	if (n.RequestAccess){} else {
    n.mu.Lock()
    n.RequestAccess = true
	fmt.Printf("Node %d requests access to critical section\n at time %s", n.id, time.Now().Format("15:04:05"))
    n.mu.Unlock()
	}
}

func RunNode(id int32, port string, nextNodeAddress string, HasToken bool) *Node {
	node := &Node{
		id:              id,
		HasToken:        HasToken,
		NextNodeAddress: nextNodeAddress,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterTokenringServer(grpcServer, node)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go func() {
        for {
            time.Sleep(time.Second * 45)
            node.RequestCriticalSection()
        }
    }()

	if HasToken {
		time.Sleep(time.Second * 10)
		node.passTokenToNext()
	}
	for{}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: TokenRing <id> <port> <nextNodeAddress> <HasToken>")
	}

	id, err := strconv.Atoi(os.Args[1])
	port := os.Args[2]
	nextNodeAddress := os.Args[3]
	HasToken := os.Args[4] == "true"
	if err != nil {
		log.Fatalf("Conversion error:, %v", err)
	}

	RunNode(int32(id), port, nextNodeAddress, HasToken)

	for{}
}
