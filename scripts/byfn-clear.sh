#!/usr/bin/env bash

if [ -d "_tmp/fabric-samples/first-network" ]; then
    cd _tmp/fabric-samples/first-network
    ./byfn.sh down
fi

