package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/config"
	"github.com/multiformats/go-multiaddr"
)

type conf struct {
	listen        config.Option
	relay         config.Option
	noInteractive bool
}

func parseFlags() *conf {
	c := &conf{}
	var listen string
	var relay bool
	flag.StringVar(&listen, "listen", "", "List of addrs the peer will listen to, separated by space")
	flag.BoolVar(&relay, "relay", false, "Should this peer relay connections to other peers? (Useful if you want this peer to act as a gateway)")
	flag.BoolVar(&c.noInteractive, "no-interactive", false, "Run without user intervention. For ping and relay purposes")

	flag.Parse()

	if relay {
		c.relay = libp2p.EnableRelay(circuit.OptHop)
	} else {
		c.relay = libp2p.EnableRelay(circuit.OptDiscovery)
	}
	if listen == "" {
		c.listen = libp2p.ListenAddrStrings()
	} else {
		addrs := strings.Split(listen, " ")
		c.listen = libp2p.ListenAddrStrings(addrs...)
	}
	return c
}

func parseInput(inp string) []string {
	noNewLine := strings.Split(inp, "\n")[0]
	splitAddr := strings.Split(noNewLine, "[")
	splitOpts := strings.Split(splitAddr[0], " ")
	if 1 < len(splitAddr) {
		return append(splitOpts[:len(splitOpts)-1], "["+splitAddr[1])
	}
	return splitOpts
}

func parseID(from string) peer.ID {
	ID, err := peer.IDHexDecode(from)
	if err != nil {
		fmt.Println("Invalid peer ID:", err)
	}
	return ID
}

func getAddrList(from string) []string {
	peerAddrString := strings.Replace(strings.Replace(from, "[", "", 1), "]", "", 1)
	return strings.Split(peerAddrString, " ")
}

func parseAddr(from string) []multiaddr.Multiaddr {
	peerAddrList := getAddrList(from)
	peerAddr := make([]multiaddr.Multiaddr, len(peerAddrList))
	for i := 0; i < len(peerAddrList); i++ {
		addr, err := multiaddr.NewMultiaddr(peerAddrList[i])
		if err != nil {
			fmt.Println("Invalid peer address:", err)
		}
		peerAddr[i] = addr
	}
	return peerAddr
}

func parsePeerInfo(ID, addr string) peer.AddrInfo {
	return peer.AddrInfo{
		ID:    parseID(ID),
		Addrs: parseAddr(addr),
	}
}

func printStreamInfo(s network.Stream) {
	fmt.Println("Message comes through:", s.Conn().RemoteMultiaddr())
	fmt.Println("By: ", s.Conn().RemotePeer(),
	"(", peer.IDHexEncode(s.Conn().RemotePeer()), ")")
}
