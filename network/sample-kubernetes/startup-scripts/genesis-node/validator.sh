#!/bin/sh
apt-get update
apt-get install curl -y 

echo ${OPENTSDB_USERNAME} >> /etc/sawtooth/validator.toml 
echo ${OPENTSDB_PW} >> /etc/sawtooth/validator.toml

if [ ! -e /etc/sawtooth/keys/validator.priv ]; then 
    sawadm keygen; 
fi

poet enclave basename --enclave-module simulator
poet registration create --enclave-module simulator 

if [ ! -e "$SAWTOOTH_HOME/data/block-chain-id" ]; then
    sawadm genesis
fi

sawtooth-validator -vv \
    --endpoint tcp://$PROPSCHAIN_GENESIS_NODE_SERVICE_HOST:8800 \
    --bind component:tcp://eth0:4004 \
    --bind network:tcp://eth0:8800 \
#    --opentsdb-url http://propschain-metrics:8086 \
#    --opentsdb-db metrics \