FROM golang:1.9 as builder
WORKDIR /go/src/github.com/DOSNetwork/core
COPY . .
#ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/local/bin/dep
#RUN chmod +x /usr/local/bin/dep
#RUN dep ensure --vendor-only
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && go build -ldflags "-linkmode external -extldflags -static" -a -o clientNode main.go

# STEP 2 build a small image
FROM scratch
COPY --from=builder /go/src/github.com/DOSNetwork/core/clientNode /
COPY --from=builder /go/src/github.com/DOSNetwork/core/onChain.json /
COPY --from=builder /go/src/github.com/DOSNetwork/core/offChain.json /
COPY --from=builder /go/src/github.com/DOSNetwork/core/testAccounts/bootCredential /credential/
CMD ["/clientNode"]