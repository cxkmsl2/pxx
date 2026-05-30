FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /gateway-service ./gateway/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /gateway-service /usr/local/bin/gateway-service
EXPOSE 8080
CMD ["gateway-service"]
