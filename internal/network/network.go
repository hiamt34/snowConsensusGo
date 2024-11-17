package network

import (
	"math/rand"
	"sync"

	"snowConsensusGo/internal/types"
)

type Network struct {
	Nodes []types.NodeInterface
	Mutex sync.Mutex
}

func NewNetwork() *Network {
	return &Network{
		Nodes: make([]types.NodeInterface, 0),
	}
}

func (net *Network) AddNode(node types.NodeInterface) {
	net.Mutex.Lock()
	defer net.Mutex.Unlock()
	net.Nodes = append(net.Nodes, node)
}

func (net *Network) SampleNodes(k int) []types.NodeInterface {
	net.Mutex.Lock()
	defer net.Mutex.Unlock()

	if len(net.Nodes) < k {
		k = len(net.Nodes)
	}

	shuffledNodes := make([]types.NodeInterface, len(net.Nodes))
	copy(shuffledNodes, net.Nodes)
	rand.Shuffle(len(shuffledNodes), func(i, j int) {
		shuffledNodes[i], shuffledNodes[j] = shuffledNodes[j], shuffledNodes[i]
	})
	return shuffledNodes[:k]
}

func (net *Network) BroadcastTransaction(tx types.Transaction, senderID int) {
	targetNodes := net.SampleNodes(types.BROADCAST_COUNT)
	for _, node := range targetNodes {
		if node.GetNodeID() != senderID {
			go node.ReceiveTransaction(tx)
		}
	}
}
