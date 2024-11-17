package node

import (
	"fmt"
	"math/rand"
	"sync"

	"snowConsensusGo/internal/types"
)

type Node struct {
	NodeID                  int
	Network                 types.NetworkInterface
	Transactions            map[string]bool
	Confidences             map[string]float64
	BroadcastedTransactions map[string]bool
	Mutex                   sync.Mutex
}

func NewNode(id int, network types.NetworkInterface) *Node {
	return &Node{
		NodeID:                  id,
		Network:                 network,
		Transactions:            make(map[string]bool),
		Confidences:             make(map[string]float64),
		BroadcastedTransactions: make(map[string]bool),
	}
}

func (n *Node) GetNodeID() int {
	return n.NodeID
}

func (n *Node) HasTransaction(txID string) bool {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()
	return n.Transactions[txID]
}

func (n *Node) ReceiveTransaction(tx types.Transaction) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	if n.ValidateTransaction(tx) && !n.Transactions[tx.TxID] {
		n.Transactions[tx.TxID] = true
		n.BroadcastTransaction(tx)
	}
}

func (n *Node) BroadcastTransaction(tx types.Transaction) {
	if !n.BroadcastedTransactions[tx.TxID] {
		n.BroadcastedTransactions[tx.TxID] = true
		tx.BroadcastCount++
		n.Network.BroadcastTransaction(tx, n.NodeID)
	}
}

func (n *Node) ValidateTransaction(tx types.Transaction) bool {
	return rand.Float64() > 0.2
}

func (n *Node) QueryValidators(tx types.Transaction) float64 {
	selectedValidators := n.Network.SampleNodes(types.K)
	validVotes := 0
	for _, validator := range selectedValidators {
		if validator.HasTransaction(tx.TxID) {
			validVotes++
		}
	}
	return float64(validVotes) / float64(types.K)
}

// func (n *Node) StartConsensus() { // V1 thăm dò đồng thuận thế này độ tin tưởng sẽ thấp, qua thử nhiệm test thấy vẫn có trường hợp tất cả các node không hội tụ
// 	n.Mutex.Lock()
// 	txIDs := make([]string, 0, len(n.Transactions))
// 	for txID := range n.Transactions {
// 		txIDs = append(txIDs, txID)
// 	}
// 	n.Mutex.Unlock()

// 	for _, txID := range txIDs {
// 		tx := types.Transaction{TxID: txID}
// 		confidence := n.QueryValidators(tx)
// 		n.Confidences[txID] = confidence
// 		if confidence >= types.CONFIDENCE_THRESHOLD {
// 			fmt.Printf("Node %d accepted transaction %s with confidence %.2f\n", n.NodeID, txID, confidence)
// 		} else {
// 			fmt.Printf("Node %d rejected transaction %s with confidence %.2f\n", n.NodeID, txID, confidence)
// 		}
// 	}
// }


func (n *Node) StartConsensus() {
    n.Mutex.Lock()
    txIDs := make([]string, 0, len(n.Transactions))
    for txID := range n.Transactions {
        txIDs = append(txIDs, txID)
    }
    n.Mutex.Unlock()

    for _, txID := range txIDs {
        // Khởi tạo các biến cần thiết cho quá trình đồng thuận
        decisionMade := false
        preference := true // Giả sử ban đầu node ủng hộ giao dịch
        consecutiveSuccesses := 0  //số lượt thành công liên tiếp

        for round := 0; round < types.MAX_ROUNDS && !decisionMade; round++ {
            selectedValidators := n.Network.SampleNodes(types.K)
            votesFor := 0
            votesAgainst := 0

            for _, validator := range selectedValidators {
                if validator.HasTransaction(txID) {
                    votesFor++
                } else {
                    votesAgainst++
                }
            }
            majorityPreference := float64(votesFor/types.K) >= types.CONFIDENCE_THRESHOLD

            if majorityPreference == preference {
                consecutiveSuccesses++
                if consecutiveSuccesses >= types.BETA {
                    // Chốt quyết định
                    decisionMade = true
                    n.Mutex.Lock()
                    n.Confidences[txID] = float64(consecutiveSuccesses) / float64(round+1)
                    n.Mutex.Unlock()
                    if preference {
                        fmt.Printf("Node %d accepted transaction %s after %d rounds\n", n.NodeID, txID, round+1)
                    } else {
                        fmt.Printf("Node %d rejected transaction %s after %d rounds\n", n.NodeID, txID, round+1)
                    }
                }
            } else {
                // Thay đổi ưu tiên và đặt lại bộ đếm
                preference = majorityPreference
                consecutiveSuccesses = 0
            }
        }

        if !decisionMade {
            fmt.Printf("Node %d could not finalize decision on transaction %s after %d rounds\n", n.NodeID, txID, types.MAX_ROUNDS)
        }
    }
}

