#!/usr/bin/env bash

echo "Preparing BYFN network ..."

cd test/byfn

if [ ! -d "fabric-samples" ]; then
    echo "Cloning https://github.com/hyperledger/fabric-samples@release-1.4 into ./test/fabric-samples"
    git clone --single-branch --branch release-1.4 https://github.com/hyperledger/fabric-samples
fi

echo "Changing byfn.sh script to start without prompt"
cd fabric-samples/first-network
sed -ie 's/^askProceed/# commented to start network without prompt - askProceed/' byfn.sh

./byfn.sh up
