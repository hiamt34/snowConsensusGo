# SnowConsensusGo - Avalanche Consensus Algorithm Implementation
## Tài liệu tự phân tích và implement thuật toán:
```
    https://docs.google.com/document/d/1DwAp0uueI-Tb5qTO8Sg4LC3yX_Tb6uSO3p_XOg7Bcko/edit?tab=t.0
```
## **Giới thiệu**

**SnowConsensusGo** là một dự án triển khai thuật toán đồng thuận **Avalanche Consensus** trong Go. Avalanche là một giao thức đồng thuận phân tán, nổi bật với khả năng mở rộng cao, độ trễ thấp và tính an toàn. Thuật toán này dựa trên việc lấy mẫu ngẫu nhiên và hội tụ xác suất để đạt đến sự đồng thuận của các node trong mạng.

Dự án này mô phỏng quá trình đồng thuận bằng cách tạo ra các node, phát tán giao dịch giữa chúng và sử dụng cơ chế lấy mẫu để xác định xem các node có thể đạt đồng thuận về trạng thái của giao dịch hay không.

## **Tính năng chính**

- Triển khai thuật toán đồng thuận Avalanche với khả năng phát tán giao dịch giữa các node.
- Các node thực hiện đồng thuận thông qua việc lấy mẫu các node khác và đánh giá sự tin tưởng vào giao dịch.
- Khả năng kiểm thử để đảm bảo tất cả các node hội tụ về cùng một trạng thái cho một giao dịch.

## **Cấu trúc dự án**

- `internal/network/`: Quản lý mạng lưới và giao tiếp giữa các node.
- `internal/node/`: Triển khai các node và quy trình đồng thuận.
- `internal/types/`: Định nghĩa các kiểu dữ liệu và tham số cấu hình.
- `cmd/main/`: Chương trình chính để chạy mô phỏng.
- `tests/`: Các bài kiểm thử để đảm bảo tính đúng đắn của thuật toán.

## **Yêu cầu hệ thống**

- **Go 1.18+**

## **Cài đặt**

1. **Clone dự án từ GitHub**:

   ```bash
   git clone https://github.com/hiamt34/snowConsensusGo.git
   cd snowConsensusGo
   ```

2. **Cài đặt các phụ thuộc**:

   Tải các phụ thuộc cần thiết bằng lệnh `go mod tidy`:

   ```bash
   go mod tidy
   ```

## **Cách sử dụng**

1. **Chạy mô phỏng**:

   Để chạy mô phỏng đồng thuận giữa các node, sử dụng lệnh:

   ```bash
   go run cmd/main/main.go
   ```

   Mô phỏng sẽ tạo ra các node, phát tán giao dịch và tiến hành đồng thuận theo thuật toán Avalanche.

2. **Kiểm thử**:

   Để chạy kiểm thử và xác minh tính hội tụ của đồng thuận, sử dụng lệnh:

   ```bash
   go test ./tests -v
   ```

   Bài kiểm thử sẽ xác minh rằng tất cả các node đều đạt đồng thuận về cùng một trạng thái giao dịch (chấp nhận hoặc từ chối).

## **Tham số cấu hình**

Các tham số của thuật toán có thể được cấu hình trong tệp `internal/types/types.go`. Các tham số chính bao gồm:

- `NUM_NODES`: Số lượng node trong mạng.
- `K`: Kích thước của tập lấy mẫu trong mỗi vòng đồng thuận.
- `BETA`: Ngưỡng số lượt đồng thuận thành công liên tiếp để chốt quyết định.
- `MAX_ROUNDS`: Số vòng lặp tối đa để đạt đồng thuận.
- `CONFIDENCE_THRESHOLD`: Ngưỡng tin cậy để quyết định chấp nhận giao dịch.

Bạn có thể thay đổi các giá trị này để kiểm tra và tối ưu quá trình đồng thuận.

## **Ví dụ**

Để chạy mô phỏng với **200 node**, kích thước lấy mẫu là **10**, và số vòng lặp tối đa là **20**, bạn có thể cấu hình trong tệp `types.go` như sau:

```go
var (
    NUM_NODES            = 200
    K                    = 10
    BETA                 = 10
    MAX_ROUNDS           = 20
    CONFIDENCE_THRESHOLD = 0.8
)
```

