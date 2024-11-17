package types

type Transaction struct {
	TxID           string
	BroadcastCount int
}

var (
	NUM_PROCESSES        = 10
	NODES_PER_PROCESS    = 20
	NUM_NODES            = NUM_PROCESSES * NODES_PER_PROCESS
	BROADCAST_COUNT      = 5     // Number of nodes to broadcast the transaction to - số lượng node để lan truyền khi 1 tx được tạo
	CONFIDENCE_THRESHOLD = 0.5  // Threshold for consensus - nguong quyet dinh
	BETA                 = 10    // Ngưỡng để chốt quyết định  - số vòng lặp liên tiếp mà một node cần quan sát đa số node đồng ý trước khi chốt trạng thái.
    MAX_ROUNDS           = 20   // Số vòng lặp tối đa - Đảm bảo node có đủ cơ hội để quan sát trạng thái mạng và hội tụ về quyết định cuối cùng.
	K                    = 5    // Number of validators to sample 2% - 5% để đảm bảo độ trung thực (quá ít thì độ tin cậy kém quá cao thì, letancy kém)
)

type NodeInterface interface {
	ReceiveTransaction(tx Transaction)
	StartConsensus()
	GetNodeID() int
	ValidateTransaction(tx Transaction) bool
	HasTransaction(txID string) bool
}

type NetworkInterface interface {
	AddNode(node NodeInterface)
	SampleNodes(k int) []NodeInterface
	BroadcastTransaction(tx Transaction, senderID int)
}
