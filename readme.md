# Fabric REST API Go

REST control system for [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) network.

Features
* Invoke and query chaincode
* Chaincode detailed info
* Channels list
* Channel peers
* CA and transactions light cryptography
* ...

## Local test instructions

Requirements:
* https://hyperledger-fabric.readthedocs.io/en/release-1.4/install.html
* https://github.com/hyperledger/fabric-sdk-go
* Go 1.12 

Clone repo outside GOPATH in order to go module work correctly.

Build and run
```
make run
```

*NOTE: only one test network could be up at the same time*

Run all unit and end-2-end tests

```
make test
```

Stop and clear all test networks
```
make clear
```

### Basic network test

Start local basic test network
```
make basic_up
```

To create and install/instantiate chaincode do POST request to 
*localhost:8080/init_test_fixtures*
 
Run E2E tests against Basic network
```
make basic_e2e_test
```

Postman collection

*/test/FabricApiBasic.postman_collection.json* 


### BYFN test

Additional requirements:
* [qj](https://stedolan.github.io/jq/) command-line JSON processor

Start local basic test network
```
make byfn_up
```

Run E2E tests against BYFN
```
make byfn_e2e_test
```

Postman collection

*/test/FabricApiBYFN.postman_collection.json*

## Light cryptography demo

Start BYFN network with API and demo app at http://localhost:3000/.
Also script will register demo user in ORG1 CA with creds - CaUser:password

```
make demo_light_up
```

Stop demo app and clear containers

```
make demo_light_down
```


Example sequence of actions to interact with standard A to B transfer chaincode:

#### Enroll user:

1) "GET TBS CSR /CA/TBSCSR"
2) "SIGN TBS CSR"
3) "ENROLL TO CA WITH CREDS, TBS CSR AND SIGNATURE - /CA/ENROLL_CSR"

#### Query:

1) Fill out five fields with: mychannel / mycc / Org1MSP / query / a
2) "CALL /TX/PROPOSAL"
3) "SIGN PROPOSAL"
4) "QUERY PROPOSAL"

#### Invoke:

1) Fill out five fields with: mychannel / mycc / Org1MSP / invoke / a,b,15
2) "CALL /TX/PROPOSAL"
3) "SIGN PROPOSAL"
4) Fill out endorsement peers with: org1/peer0,org2/peer0
5) "GET BROADCAST PAYLOAD (PROPOSAL+ENDORSMENTS)"
6) "SIGN BROADCAST PAYLOAD"
7) "BROADCAST PAYLOAD WITH SIGNATURE"

After finishing invoke sequence you may run query sequence again, to check out that value have been changed.


## Work progress 

### Close plans

* Channels creation
* Network management 
* Documentation 

### Strategic plans

* Cover code with unit test (partially done)

