FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/lsc-chocos/
COPY . .
RUN go mod tidy
RUN go build -o /exe main.go

FROM golang:1.13
COPY --from=builder /go/src/github.com/lsc-chocos/ssl/mainflux-server.crt /mainflux-server.crt
COPY --from=builder /go/src/github.com/lsc-chocos/configs/config.json /config.json
COPY --from=builder /exe /
ENTRYPOINT [ "/exe", "-f", "/config.json", "--cacert", "/mainflux-server.crt"]
