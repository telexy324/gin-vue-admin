#FROM golang:1.19-alpine as builder
FROM golang:alpine as builder

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server
COPY . .

#RUN apk add --no-cache git
ARG CURRENTVERSION
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=0
RUN go env
RUN go mod tidy
RUN go build -ldflags="-X 'github.com/flipped-aurora/gin-vue-admin/server/core.buildTime=$(date '+%Y-%m-%dT%H:%M:%S')' \
    -X 'github.com/flipped-aurora/gin-vue-admin/server/core.currentVersion=$CURRENTVERSION' \
    " -o server .
#RUN go build -o server .

FROM alpine:latest
LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server

COPY --from=0 /go/src/github.com/flipped-aurora/gin-vue-admin/server ./

EXPOSE 8888

ENTRYPOINT ./server -c config.docker.yaml
