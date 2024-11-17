//snowConsensusGo/cmd/main/main.go
package main

import (
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"snowConsensusGo/internal/network"
	"snowConsensusGo/internal/node"
	"snowConsensusGo/internal/types"
)

func nodeProcess(startID, endID int, net *network.Network) {
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
	var wg sync.WaitGroup

	for i := 0; i < types.NUM_PROCESSES; i++ {
		startID := i * types.NODES_PER_PROCESS
		endID := (i + 1) * types.NODES_PER_PROCESS
		wg.Add(1)
		go func(startID, endID int) {
			defer wg.Done()
			cmd := exec.Command(os.Args[0], strconv.Itoa(startID), strconv.Itoa(endID))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}
		}(startID, endID)
	}

	wg.Wait()
}

func init() {
	if len(os.Args) > 2 {
		startID, _ := strconv.Atoi(os.Args[1])
		endID, _ := strconv.Atoi(os.Args[2])
		net := network.NewNetwork()
		nodeProcess(startID, endID, net)
		os.Exit(0)
	}
}
