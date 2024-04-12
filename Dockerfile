# 使用官方Golang鏡像作為構建環境
FROM golang:latest AS builder

# 設置工作目錄
WORKDIR /app

# 複製go mod和sum文件
COPY go.mod ./
COPY go.sum ./

# 下載依賴包
RUN go mod download

# 複製源代碼到容器中
COPY . .

WORKDIR /app/cmd/Dcard-Backend-HW

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 使用scratch作為運行環境
FROM scratch

# 從builder階段複製構建的二進制文件
COPY --from=builder /app/cmd/Dcard-Backend-HW .

# 啟動應用
CMD ["./main"]
