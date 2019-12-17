package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type config struct {
	relayInfo     peer.AddrInfo
	setConnection bool
	peerInfo      peer.AddrInfo
}

func parseFlags() *config {
	c := &config{}
	var relayID string
	var relayAddr string
	var peerID string
	var peerAddr string
	flag.StringVar(&relayID, "relayID", "foo", "ID of the relay node")
	flag.StringVar(&relayAddr, "relayAddr", "foo", "URL of the relay node")
	flag.StringVar(&peerID, "peerID", "foo", "ID of the peer to be connected to. (required if setConnection == true)")
	flag.StringVar(&peerAddr, "peerAddr", "foo", "URL of the peer to be connected to. (required if setConnection == true)")
	flag.BoolVar(&c.setConnection, "setConnection", false, "Should this node start the connection, or wait to be connected?")

	flag.Parse()
	if c.setConnection {
		c.peerInfo = peer.AddrInfo{
			ID:    parseID(peerID),
			Addrs: parseAddr(peerAddr),
		}
	}
	c.relayInfo = peer.AddrInfo{
		ID:    parseID(relayID),
		Addrs: parseAddr(relayAddr),
	}
	return c
}

func parseID(from string) peer.ID {
	ID, err := peer.IDHexDecode(from)
	if err != nil {
		fmt.Println("Invalid peer ID:", err)
		panic(err)
	}
	return ID
}

func parseAddr(from string) []multiaddr.Multiaddr {
	peerAddrString := strings.Replace(strings.Replace(from, "[", "", 1), "]", "", 1)
	peerAddrList := strings.Split(peerAddrString, " ")
	peerAddr := make([]multiaddr.Multiaddr, len(peerAddrList))
	for i := 0; i < len(peerAddrList); i++ {
		addr, err := multiaddr.NewMultiaddr(peerAddrList[i])
		if err != nil {
			fmt.Println("Invalid peer address:", err)
			panic(err)
		}
		peerAddr[i] = addr
	}
	return peerAddr
}
