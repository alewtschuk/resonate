package nodeconnection

import (
	"context"
	"time"

	"github.com/charmbracelet/log"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

var booststrapPeersMap map[peer.ID]bool = make(map[peer.ID]bool)
var peerAddrInfos []peer.AddrInfo = dht.GetDefaultBootstrapPeerAddrInfos()

// Creates the DHT for the node, sets the mode for DHT handling (server or auto), and spins up the discover function in a new goroutine.
func InitDHTDiscovery(ctx context.Context, node host.Host, DISCOVERY_SERVICE_TAG string, bootstrapPeers []multiaddr.Multiaddr) {
	log.Info("Configuring the node's DHT...")
	var opts []dht.Option = []dht.Option{dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...)} // sets the bootstrap peers as the libp2p default bootstrap nodes

	// Setup the bootstrap peers map
	for _, peerAddrInfo := range peerAddrInfos {
		booststrapPeersMap[peerAddrInfo.ID] = true
	}

	// Sets the nodes dht mode to auto to select if it will be a server or a client
	//TODO: Add cli flags to specify if the user is planning on running a persistant node and set persistant nodes to server mode.
	opts = append(opts, dht.Mode(dht.ModeAuto))
	log.Info("Node DHT mode: Auto")

	// Create new DHT with specified options
	dht, err := dht.New(ctx, node, opts...)
	if err == nil {
		log.Info("DHT creation successfull!")
	} else {
		checkErr(err)
	}

	log.Info("Populating node DHT")
	err = dht.Bootstrap(ctx)
	if err == nil {
		log.Info("Node DHT population satisfactory")
	} else {
		checkErr(err)
	}

	go discover(ctx, node, dht, DISCOVERY_SERVICE_TAG)
}

func discover(ctx context.Context, node host.Host, dht *dht.IpfsDHT, DISCOVERY_SERVICE_TAG string) {
	routingDiscovery := routing.NewRoutingDiscovery(dht)

	discovery.Advertise(ctx, routingDiscovery, DISCOVERY_SERVICE_TAG)
	log.Info("Node now advertizing via DHT\n")

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers, err := discovery.FindPeers(ctx, routingDiscovery, DISCOVERY_SERVICE_TAG)
			checkErr(err)

			for _, peer := range peers {
				//dhtPeerLen := len(peers)
				if peer.ID == node.ID() {
					//dhtPeerLen--
					continue
				}

				//log.Infof("Current num of dht peers: %d", dhtPeerLen)
				log.Infof("Discovered DHT peer: %s", peer.ID.String())
				if booststrapPeersMap[peer.ID] {
					log.Infof("Found Bootstrap peer %s", peer.ID.String())
				} else {
					log.Infof("Found Application peer %s\n", peer.ID.String())
				}

				if node.Network().Connectedness(peer.ID) != network.Connected {
					_, err := node.Network().DialPeer(ctx, peer.ID)
					if err != nil {
						log.Errorf("Error connecting to peer %v", err)
						continue
					} else {
						if booststrapPeersMap[peer.ID] {
							log.Infof("Connected to Bootstrap peer %s", peer.ID.String())
						} else {
							log.Infof("Connected to Application peer node %s\n", peer.ID.String())
						}
					}
				}
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Errorf("An error has occured: %v", err)
		panic(err)
	}
}
