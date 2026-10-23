package main

import (
	"context"
	"flag"
	"io"
	"resonate/internal/nodeconnection"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
)

const DISCOVERY_SERVICE_TAG string = "resonate-app-node-v0.0.1-adjfhadkljfhalkduhfalkjdfhalekj"

func main() {
	setFlags()
	ctx := context.Background()
	// Create listener node
	node := createNode()
	getNodeInfo(node)
	// err := nodeconnection.InitMDNSDiscovery(node, DISCOVERY_SERVICE_TAG) // creates an mDNS discovery service for the node
	// if err == nil {
	// 	log.Info("Successfully started mDNS service!")
	// } else {
	// 	checkErr(err)
	// }
	nodeconnection.InitDHTDiscovery(ctx, node, DISCOVERY_SERVICE_TAG, nil)
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
		log.Errorf("An error has occured: %v", err)
		panic(err)
	}
}

// Prints node info
func getNodeInfo(node host.Host) {
	log.Infof("Node ID: %s", node.ID().String())
	log.Infof("Node Address: %s\n", node.Addrs())
}

func setFlags() {
	debugFlag := flag.Bool("debug", true, "Set debug flag level")
	flag.Parse()

	if !*debugFlag {
		log.Warn("Running without debug output")
		log.SetOutput(io.Discard)
	}
}

// Sets logger level
func init() {
	log.SetReportTimestamp(false)
}
