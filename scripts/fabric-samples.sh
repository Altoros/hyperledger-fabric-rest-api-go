#!/usr/bin/env bash

function cloneFabricSamples() {
    if [ ! -d "./_tmp" ]; then
        mkdir ./_tmp
    fi

    if [ ! -d "./_tmp/fabric-samples" ]; then
        echo "Cloning https://github.com/hyperledger/fabric-samples@release-1.4 into ./_tmp/fabric-samples"
        git clone --single-branch --branch release-1.4 https://github.com/hyperledger/fabric-samples ./_tmp/fabric-samples
    fi
}
