# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
