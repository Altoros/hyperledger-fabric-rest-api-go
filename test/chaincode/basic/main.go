package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
)

type BasicTestChaincode struct {
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *BasicTestChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("-=-=- BasicTestChaincode Init -=-=-")

	// Put in the ledger the key/value hello/world
	err := stub.PutState("value", []byte("Hello world!"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return a successful message
	return shim.Success(nil)
}

// All future requests named invoke will arrive here.
func (t *BasicTestChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("-=-=- BasicTestChaincode Invoke -=-=-")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient")
	}

	switch function {
	case "query":
		return t.query(stub, args)
	case "update":
		return t.update(stub, args)
	case "user":
		return t.userInfo(stub)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// Every readonly functions in the ledger will be here
func (t *BasicTestChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("-=-=- BasicTestChaincode query -=-=-")

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Get the state of the value matching the key hello in the ledger
	state, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get state of value")
	}

	// Return this value in response
	return shim.Success(state)
}

// Every functions that read and write in the ledger will be here
func (t *BasicTestChaincode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("-=-=- BasicTestChaincode invoke -=-=-")

	if len(args) != 2 {
		return shim.Error("The number of arguments is incorrect, need 2 arguments.")
	}

	// Write the new value in the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error("Failed to update state of hello")
	}

	err = stub.SetEvent("eventInvoke", []byte("Big payload"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return this value in response
	return shim.Success(nil)
}

// Fetch user info
func (t *BasicTestChaincode) userInfo(stub shim.ChaincodeStubInterface) pb.Response {
	creatorBytes, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	name, org := getCreator(creatorBytes)

	// Return this value in response
	return shim.Success([]byte(fmt.Sprintf("Name: %s, Org: %s, Pem: %s", name, org, creatorBytes)))
}

func getCreator(certificate []byte) (commonName, organization string) {
	data := certificate[strings.Index(string(certificate), "-----") : strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode(data)
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization = cert.Issuer.Organization[0]
	commonName = cert.Subject.CommonName
	return
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(BasicTestChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
