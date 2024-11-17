package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"snowConsensusGo/internal/network"
	"snowConsensusGo/internal/node"
	"snowConsensusGo/internal/types"
)

func nodeProcess(startID, endID int, net *network.Network, wg *sync.WaitGroup) {
	defer wg.Done()

	nodes := make([]*node.Node, 0)
	for nodeID := startID; nodeID < endID; nodeID++ {
		n := node.NewNode(nodeID, net)
		nodes = append(nodes, n)
		net.AddNode(n)
	}

	for _, n := range nodes {
		tx := types.Transaction{TxID: "tx_" + uuid.New().String()}
		n.ReceiveTransaction(tx)
	}

	time.Sleep(1 * time.Second)

	for _, n := range nodes {
		n.StartConsensus()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	net := network.NewNetwork()
	var wg sync.WaitGroup

	for i := 0; i < types.NUM_PROCESSES; i++ {
		startID := i * types.NODES_PER_PROCESS
		endID := (i + 1) * types.NODES_PER_PROCESS
		wg.Add(1)
		go nodeProcess(startID, endID, net, &wg)
	}

	wg.Wait()
}
