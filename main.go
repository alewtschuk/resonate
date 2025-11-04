package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background() // create background context
	// Create listener node
	listenerNode := createNode()
	getNodeInfo(listenerNode)

	listenerAddr := listenerNode.Addrs()[0]                                                  // gets the first node address from the multiaddress list
	var fullAddrString string = listenerAddr.String() + "/p2p/" + listenerNode.ID().String() // dynamically creates the address string
	log.Info("Full addr to dial is: ", "addr", fullAddrString)

	// Create dialer node that will connect to listener
	dialerNode := createNode()
	getNodeInfo(dialerNode)

	// Converts the address string into a valid multiaddress
	m, err := multiaddr.NewMultiaddr(fullAddrString)
	checkErr(err)

	// Creates the peer address to that will be used for connection
	paddr, err := peer.AddrInfoFromP2pAddr(m)
	checkErr(err)

	// Connect the dialer node to the listener node
	connErr := dialerNode.Connect(ctx, *paddr)
	if connErr == nil {
		log.Info("CONNECTION SUCESSFULL")
	} else {
		checkErr(connErr)
	}
	select {}
}

// Creates and returns a new node with included options
func createNode() host.Host {
	log.Info("Creating Node...")
	// Sets options array for the node
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/udp/0/quic-v1"), // sets the node to use ip4 with random address via udp quic-v1 protocol
		libp2p.Transport(quic.NewTransport),                    // sets the transport
	}

	// creates the new node
	node, err := libp2p.New(opts...)
	checkErr(err)
	log.Info("Node Created Successfully!")
	return node
}

// Checks any passed error and panics on error
func checkErr(err error) {
	if err != nil {
		log.Info("An error has occured", "err", err)
		panic(err)
	}
}

// Prints node info
func getNodeInfo(node host.Host) {
	log.Infof("Node ID: %s", node.ID())
	log.Infof("Node Address: %s\n", node.Addrs())
}

// Sets logger level
func init() {
	log.SetReportTimestamp(false)
}
