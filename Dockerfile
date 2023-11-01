FROM golang:alpine-1.20.10 AS builder
LABEL stage=4everland-ipfs-top
ARG APP_NAME
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /build/ipfs
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o /runner/ipfs/server ./app/${APP_NAME}/cmd/main.go ./app/${APP_NAME}/cmd/wire_gen.go

FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ UTC
WORKDIR /runner/ipfs
COPY --from=builder /runner/ipfs/server /runner/ipfs/server
CMD ["./server"]
