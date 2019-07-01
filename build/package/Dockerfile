FROM ubuntu:bionic AS multistage

RUN apt-get update \
    && apt-get install gnupg -y

RUN echo "deb [arch=amd64] http://repo.sawtooth.me/ubuntu/ci bionic universe" >> /etc/apt/sources.list \
 && echo 'deb http://ppa.launchpad.net/gophers/archive/ubuntu bionic main' >> /etc/apt/sources.list \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 8AA7AF1F1091A5FD) \
 && (apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 308C15A29AD198E9 \
 || apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 308C15A29AD198E9) \
 && apt-get update \
 && apt-get install -y -q \
    golang-1.11 \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    protobuf-compiler \
    python3 \
    python3-grpcio \
    python3-grpcio-tools \
    python3-pkg-resources \
    ca-certificates \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /go

ENV GOPATH=/go
ENV PATH=$PATH:/usr/lib/go-1.11/bin:/go/bin

RUN go get -u \
    github.com/btcsuite/btcd/btcec \
    github.com/golang/protobuf/proto \
    github.com/golang/protobuf/protoc-gen-go \
    github.com/golang/mock/gomock \
    github.com/golang/mock/mockgen \
    github.com/pebbe/zmq4 \
    github.com/satori/go.uuid \
    github.com/hyperledger/sawtooth-sdk-go

RUN ls -al /go/bin

WORKDIR /go/src/github.com/hyperledger/sawtooth-sdk-go
RUN go generate

RUN mkdir -p /go/src/github.com/propsproject/props-transaction-processor
WORKDIR /go/src/github.com/propsproject/props-transaction-processor
COPY . .

RUN make deps
RUN make build

FROM ubuntu:18.04
RUN apt-get update \
    && apt-get install -y -q libzmq3-dev ca-certificates curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY --from=multistage /go/src/github.com/propsproject/props-transaction-processor/build/bin/pending-props-tp /go/bin/
COPY --from=multistage /go/src/github.com/propsproject/props-transaction-processor/configs/default.json /go/bin/configs/
EXPOSE 3000
ENTRYPOINT ["/go/bin/pending-props-tp -c -f /go/bin/configs/default.json"]

EXPOSE 4004/tcp
