package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	config := parseFlags()
	// create peer
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptDiscovery))
	if err != nil {
		panic(err)
	}
	fmt.Println("Hi I am a peer:", peer.IDHexEncode(host.ID()), "\nYou can find me at:", host.Addrs())

	// connect to the relay
	if err := host.Connect(context.Background(), config.relayInfo); err != nil {
		panic(err)
	}
	// connect to the peer who is laready waiting
	if config.setConnection {
		// Connect to the peer
		if err := host.Connect(context.Background(), config.peerInfo); err != nil {
			fmt.Println("Failed to connect to peer:", err)
			return
		}
		// Create stream
		s, err := host.NewStream(context.Background(), config.peerInfo.ID, "/relay-test")
		if err != nil {
			fmt.Println("Error stablishing connection: ", err)
			return
		}
		fmt.Println("Connection stablished :D (I have knocked on someones door)")
		s.Read(make([]byte, 1)) // block until the handler closes the stream
	} else {
		// prepare to receive stream
		host.SetStreamHandler("/relay-test", func(s network.Stream) {
			fmt.Println("Connection stablished :D (someone knocked on my door)")
			s.Close()
		})
		fmt.Println("Waiting for someone to connect with me...")
	}

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

	// shut the node down
	if err := host.Close(); err != nil {
		panic(err)
	}
}
