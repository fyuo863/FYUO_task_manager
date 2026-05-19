# ==============================
# 构建阶段
# ==============================
FROM golang:1.26-alpine AS builder

# 设置 Go 代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache git

WORKDIR /app

# 先复制依赖文件，利用 Docker 层缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .

# 构建三个入口 (静态编译，去除调试信息)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/all      . \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/user     ./cmd/user \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/worker   ./cmd/worker


# ==============================
# 运行阶段 — 基础镜像 (alpine)
# ==============================
FROM alpine:3.21 AS base

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY configs/ ./configs/


# ==============================
# 目标 1: 一体化 (开发/测试用)
# ==============================
FROM base AS all
COPY --from=builder /app/bin/all /app/server
EXPOSE 8090
CMD ["/app/server"]


# ==============================
# 目标 2: 仅 API 服务
# ==============================
FROM base AS user
COPY --from=builder /app/bin/user /app/server
EXPOSE 8090
CMD ["/app/server"]


# ==============================
# 目标 3: 仅 Worker
# ==============================
FROM base AS worker
COPY --from=builder /app/bin/worker /app/worker
CMD ["/app/worker"]
