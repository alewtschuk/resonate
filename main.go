package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
)

const DISCOVERY_SERVICE_TAG string = "resonate-local-node"

// DiscoveryNotify gets notified whenever a new mDNS node is discovered
// This implements the mdns.Notifee interface by implementing the method
// HandlePeerFound()
type discoveryNotifee struct {
	currentNode host.Host
}

func main() {
	// Create listener node
	node := createNode()
	getNodeInfo(node)
	err := initMDNSDiscovery(node) // creates an mDNS discovery service for the node
	checkErr(err)
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

// HandlePeerFound connects to discovered mDNS peers on the local network
func (dn *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	// Check to ignore discovery of self
	if dn.currentNode.ID() == peerInfo.ID {
		return
	}

	log.Infof("Discovered new peer: %v\n", peerInfo.ID.String())
	err := dn.currentNode.Connect(context.Background(), peerInfo)
	if err != nil {
		log.Errorf("Error connecting to peer %s", peerInfo.ID.String())
	} else {
		log.Infof("Connected to peer node: %s\n", peerInfo.ID.String())
	}
}

// Creates a new mDNS discovery service passing the local node, the discovery tag, and the struct that implements mdns.Notifee
func initMDNSDiscovery(node host.Host) error {
	// mDNS service fills the notifee struct with a reference to the local node.
	// This struct serves as a watcher for the mDNS address.
	// When another node with the service tag is found the service notifys the referenced node in the watcher (local node)
	// and invokes the HandlePeerFound method to connect the peers.
	mdnsService := mdns.NewMdnsService(node, DISCOVERY_SERVICE_TAG, &discoveryNotifee{currentNode: node})
	return mdnsService.Start() // starts the mDNS service and returns any error
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
