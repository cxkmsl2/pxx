FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /raft-scheduler ./raft-scheduler/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /raft-scheduler /usr/local/bin/raft-scheduler
EXPOSE 7000
CMD ["raft-scheduler"]
