package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"

	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	relay, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"), libp2p.EnableRelay(circuit.OptHop))
	if err != nil {
		panic(err)
	}
	fmt.Println("Hi I am a relay:", peer.IDHexEncode(relay.ID()), "\nYou can find me at:", relay.Addrs())

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := relay.Close(); err != nil {
		panic(err)
	}
}
