#!/usr/bin/env bash

NETWORK="net_byfn"

echo "Enrolling CA admin"

docker run -ti --rm --network $NETWORK alpine/httpie POST frag:8080/ca/enroll \
    login=admin \
    password=adminpw

echo "Registering CA user"

docker run -ti --rm --network $NETWORK alpine/httpie POST frag:8080/ca/register \
    login=UserCa \
    password=password

docker build -t lightcrypto ./demo/light-crypto-demo/.

docker run -d -p 3000:8080 --rm --name lightcrypto-demo lightcrypto

echo ""
echo "Demo app - http://localhost:3000/"
echo ""
