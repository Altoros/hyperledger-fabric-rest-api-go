#!/usr/bin/env bash

echo "Preparing BYFN network ..."

. /scripts/fabric-samples.sh

cloneFabricSamples


cd test/fabric-samples/first-network

echo "Changing byfn.sh script to start without prompt"
sed -ie 's/^askProceed/# commented to start network without prompt - askProceed/' byfn.sh

./byfn.sh up
