#!/usr/bin/env bash

function cloneFabricSamples() {
    if [ ! -d "./test/fabric-samples" ]; then
        echo "Cloning https://github.com/hyperledger/fabric-samples@release-1.4.0 into ./test/fabric-samples"
        git clone --single-branch --branch release-1.4 https://github.com/hyperledger/fabric-samples ./test/fabric-samples
    fi
}
