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
	NextNodeAddress string
	mu              sync.Mutex
}

func (n *Node) ReceiveToken(ctx context.Context, in *proto.Token) (*proto.Empty, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("Node received token %d\n", n.id)
	n.HasToken = true

	time.Sleep(time.Second * 2)
	fmt.Printf("Node %d going to %s\n", n.id, n.NextNodeAddress)

	n.passTokenToNext()

	return &proto.Empty{}, nil
}

func (n *Node) passTokenToNext() {
	n.HasToken = false

	conn, err := grpc.Dial(n.NextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := proto.NewTokenringClient(conn)
	token := &proto.Token{
		Id:      n.id,
		Message: "Token passing",
	}

	fmt.Printf("Node %d sending token to %s\n", n.id, n.NextNodeAddress)
	_, err = client.Send(context.Background(), token)
	if err != nil {
		log.Fatalf("Could not send token: %s", err)
	}

}

func (n *Node) SendToken() {
	n.HasToken = false

	conn, err := grpc.Dial(n.NextNodeAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	defer conn.Close()

	client := proto.NewTokenringClient(conn)
	token := &proto.Token{
		Id:      n.id,
		Message: "Token passing",
	}

	fmt.Printf("Node %d sending token to %s\n", n.id, n.NextNodeAddress)
	_, err = client.Send(context.Background(), token)
	if err != nil {
		log.Fatalf("Could not send token: %s", err)
	}
}

func RunNode(id int32, port string, nextNodeAddress string, HasToken bool) *Node {
	node := &Node{
		id:              id,
		HasToken:        HasToken,
		NextNodeAddress: nextNodeAddress,
	}

	lis, err := net.Listen("Tcp", fmt.Sprintf(":%s", port))
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

	if HasToken {
		time.Sleep(time.Second * 2)
		node.passTokenToNext()
	}
	select {}
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
}
