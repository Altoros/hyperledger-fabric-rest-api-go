#!/usr/bin/env bash

echo "Preparing basic network ..."

. ./scripts/fabric-samples.sh

cloneFabricSamples


cd _tmp/fabric-samples/basic-network

./start.sh
