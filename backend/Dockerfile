FROM golang:1.23.4 AS builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    TZ=Asia/Shanghai

ENV LANG=C.UTF-8

WORKDIR /build

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auto-deploy-backend .

FROM alpine:3.18

ENV TZ=Asia/Shanghai
ENV RUN_MODE=prod

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /app

COPY --from=builder /build/auto-deploy-backend .

EXPOSE 8080

ENTRYPOINT ["/app/auto-deploy-backend"]
