package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

var p host.Host

func main() {
	// create peer
	config := parseFlags()
	var err error
	p, err = libp2p.New(
		context.Background(),
		config.listen,
		config.relay,
	)
	if err != nil {
		panic(err)
	}

	// listen to other peers
	p.SetStreamHandler("/ping", func(s network.Stream) {
		fmt.Println("________________________________________")
		fmt.Println("Ping!")
		printStreamInfo(s)
		fmt.Println("________________________________________")
		s.Close()
	})

	// App loop
	fmt.Println("\n\n\n========================================")
	fmt.Println("Hi!")
	printWhoami()
	fmt.Println("========================================")
	fmt.Println("You can add me by running this command in another instance of this app:")
	printAddme()
	fmt.Println("========================================")
	if config.noInteractive {
		// wait for a SIGINT or SIGTERM signal
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("Received signal, shutting down...")

		// shut the node down
		if err := p.Close(); err != nil {
			panic(err)
		}
		return
	}
	printHelp()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("________________________________________")
		fmt.Println("Enter command:")
		input, err := reader.ReadString('\n')
		if err != nil {
			panic("Error receiving input: " + err.Error())
		}
		cmd := parseInput(input)
		switch cmd[0] {
		case "help":
			printHelp()
		case "whoami":
			printWhoami()
		case "addme":
			printAddme()
		case "add":
			addPeer(cmd)
		case "ls":
			printPeerList()
		case "ping":
			ping(cmd)
		case "quit":
			return
		default:
			fmt.Println("Unrecognized command")
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("APP FLAGS (passed when starting the app):")
	fmt.Println("\n-listen\t\t\tList of addrs the peer will listen to, quoted (\"\") and separated by space. Useful if behind proxy, DNS, portforwarding, ...")
	fmt.Println("-relay\t\t\tShould this peer relay connections to other peers? Useful if you want this peer to act as a gateway. If not passed NO is assumed")
	fmt.Println("-no_interactive\t\t\tRun without user intervention. For ping and relay purposes")
	fmt.Println("\nAPP USAGE (once the app is running):")
	fmt.Println("\nhelp\t\t\t Prints this message.")
	fmt.Println("whoami\t\t\t Prints info about this peer (ID, addr).")
	fmt.Println("addme\t\t\t Prints command to be used by others in order to find you.")
	fmt.Println("add [peer ID] [addr]\t Adds a peer with the given ID and addr.")
	fmt.Println("ls\t\t\t Print list of connected peers.")
	fmt.Println("ping [peer ID] \t\t Send a ping to peer ID.")
	fmt.Println("quit\t\t\t Stop app.")
}

func printWhoami() {
	fmt.Println("I'm the peer:", p.ID().Pretty())
}

func printAddme() {
	addrs := p.Addrs()
	addrString := "[/p2p-circuit/p2p/" + p.ID().Pretty()
	for _, v := range addrs {
		addrString += " " + v.String()
	}
	addrString += "]"
	fmt.Println("add", peer.IDHexEncode(p.ID()), addrString)
}

func addPeer(args []string) {
	if len(args) != 3 {
		fmt.Println("You need to specify peer ID and Addr")
		return
	}
	peerInfo := parsePeerInfo(args[1], args[2])
	if err := p.Connect(context.Background(), peerInfo); err != nil {
		fmt.Println("Error connecting to the peer")
		return
	}
	fmt.Println("Peer added:", peerInfo.ID.Pretty())
}

func printPeerList() {
	fmt.Println(p.Peerstore().Peers())
}

func ping(args []string) {
	s, err := p.NewStream(context.Background(), parseID(args[1]), "/ping")
	if err != nil {
		fmt.Println("Error stablishing connection: ", err)
		return
	}
	fmt.Println("Pong!")
	printStreamInfo(s)
	s.Read(make([]byte, 1))
}

// TODO: import / export
