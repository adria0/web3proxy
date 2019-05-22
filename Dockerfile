FROM golang:latest as builder 
RUN apt update -y
RUN apt install -y upx
RUN mkdir /build 
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web3proxy . 
RUN upx --best web3proxy

FROM scratch
COPY --from=builder /build/web3proxy /web3proxy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/web3proxy"]
