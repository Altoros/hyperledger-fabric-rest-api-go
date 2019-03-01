#!/usr/bin/env bash

echo "Preparing basic network ..."

. ./scripts/fabric-samples.sh

cloneFabricSamples


cd test/fabric-samples/basic-network

./start.sh
