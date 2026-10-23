package nodeconnection

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// DiscoveryNotify gets notified whenever a new mDNS node is discovered
// This implements the mdns.Notifee interface by implementing the method
// HandlePeerFound()
type discoveryNotifee struct {
	currentNode host.Host
}

// Creates a new mDNS discovery service passing the local node, the discovery tag, and the struct that implements mdns.Notifee
func InitMDNSDiscovery(node host.Host, DISCOVERY_SERVICE_TAG string) error {
	// mDNS service fills the notifee struct with a reference to the local node.
	// This struct serves as a watcher for the mDNS address.
	// When another node with the service tag is found the service notifys the referenced node in the watcher (local node)
	// and invokes the HandlePeerFound method to connect the peers.
	log.Info("Starting mDNS service...")
	mdnsService := mdns.NewMdnsService(node, DISCOVERY_SERVICE_TAG, &discoveryNotifee{currentNode: node})
	return mdnsService.Start() // starts the mDNS service and returns any error
}

// HandlePeerFound connects to discovered mDNS peers on the local network
// This is automatically handled by the started mDNSService instance
func (dn *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	// Check to ignore discovery of self
	if dn.currentNode.ID() == peerInfo.ID {
		return
	}

	log.Infof("Discovered new mDNS peer: %v\n", peerInfo.ID.ShortString())
	err := dn.currentNode.Connect(context.Background(), peerInfo)
	if err != nil {
		log.Errorf("Error connecting to peer %s", peerInfo.ID.ShortString())
	} else {
		log.Infof("Connected to mDNS peer node: %s\n", peerInfo.ID.ShortString())
	}
}
