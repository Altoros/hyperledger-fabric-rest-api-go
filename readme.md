## Local test instructions

Requirements:
* https://hyperledger-fabric.readthedocs.io/en/release-1.4/install.html
* Golang 11.0 

Put in code in ~/go/src/fabric-rest-api-go
```
cd ~/go/src/fabric-rest-api-go
make run
```

*NOTE: only one test network could be up at the same time*

Stop and clear all test networks
```
make clear
```

### Basic network test

Start local network
```
make basic_up
```

To create and install/instantiate chaincode do POST request to 
*localhost:8080/init_test_fixtures*
 

Postman collection

*/test/FabricApi.postman_collection.json* 


### BYFN test

Additional requirements:
* [qj](https://stedolan.github.io/jq/) command-line JSON processor

Run integration tests against BYFN
```
make byfn
```

Postman collection

*/test/FabricApiBYFN.postman_collection.json*


#### Close plans

* Test with BYFN
* Cover code with unit test
* Create full integration test with BYFN, makefile and shell scripts
* Automate postman collections tests with Newman
