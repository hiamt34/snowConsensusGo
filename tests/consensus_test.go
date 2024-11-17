package tests

import (
    "math/rand"
    "sync"
    "testing"
    "time"

    "github.com/google/uuid"
    "snowConsensusGo/internal/network"
    "snowConsensusGo/internal/node"
    "snowConsensusGo/internal/types"
)

func TestConsensusConvergence(t *testing.T) {
    // Thiết lập hạt giống ngẫu nhiên
    rand.Seed(time.Now().UnixNano())

    // Cấu hình tham số
    types.NUM_NODES = 200
    types.K = 200 * 0.5
    types.BETA = 10
    types.MAX_ROUNDS = 20
    types.CONFIDENCE_THRESHOLD = 0.8

    // Khởi tạo mạng lưới
    net := network.NewNetwork()
    var wg sync.WaitGroup

    // Khởi tạo các node
    totalNodes := types.NUM_NODES
    nodes := make([]*node.Node, totalNodes)
    for i := 0; i < totalNodes; i++ {
        nodes[i] = node.NewNode(i, net)
        net.AddNode(nodes[i])
    }

    // Tạo một giao dịch chung cho tất cả các node
    txID := "tx_" + uuid.New().String()
    tx := types.Transaction{TxID: txID}

    // Mỗi node nhận giao dịch
    for _, n := range nodes {
        n.ReceiveTransaction(tx)
    }

    // Thực hiện quá trình đồng thuận song song trên tất cả các node
    wg.Add(len(nodes))
    for _, n := range nodes {
        go func(n *node.Node) {
            defer wg.Done()
            n.StartConsensus()
        }(n)
    }

    wg.Wait()

    // Kiểm tra kết quả đồng thuận
    var finalDecision bool
    var decisionMade bool

    for _, n := range nodes {
        n.Mutex.Lock()
        confidence, exists := n.Confidences[txID]
        n.Mutex.Unlock()

        if !exists {
            t.Errorf("Node %d did not make a decision on transaction %s", n.NodeID, txID)
            continue
        }

        accepted := confidence >= types.CONFIDENCE_THRESHOLD

        if !decisionMade {
            finalDecision = accepted
            decisionMade = true
        } else {
            if accepted != finalDecision {
                t.Errorf("Consensus failed: Node %d has a different decision on transaction %s", n.NodeID, txID)
            }
        }
    }

    if !decisionMade {
        t.Errorf("No node made a decision on transaction %s", txID)
    } else {
        if finalDecision {
            t.Logf("All nodes accepted transaction %s", txID)
        } else {
            t.Logf("All nodes rejected transaction %s", txID)
        }
    }
}
