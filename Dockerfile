FROM golang:alpine AS builder
ENV GOPROXY https://mirrors.aliyun.com/goproxy/
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" --mod=vendor -tags timetzdata

FROM alpine:latest
ENV TZ=Asia/Shanghai
COPY --from=builder /app/sloth /app/sloth
ENTRYPOINT ["/app/sloth"]

# docker build -t "sloth:latest" .
