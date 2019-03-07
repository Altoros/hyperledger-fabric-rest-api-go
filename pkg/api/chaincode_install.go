package api

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"io/ioutil"
)

func ChaincodeInstall(channelClientProvider AdminProvider, _ fab.Peer, channelId, chaincodeId, chaincodeVersion string) (string, error) {

	ccTarBytes, _ := ioutil.ReadFile("/home/avk/www/cc.tar.gz")
	ccPkg := &resource.CCPackage{Type: pb.ChaincodeSpec_GOLANG, Code: ccTarBytes}

	// TODO research path usage inside peer, seems like it is only comment of some kind
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: chaincodeId, Path: "chaincode/", Version: chaincodeVersion, Package: ccPkg}
	installCcResponses, err := channelClientProvider.Admin().InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

	if err != nil {
		return "", errors.WithMessage(err, "failed to install chaincode")
	}

	fmt.Println("Peers responses in install request")
	for _, installCcResponse := range installCcResponses {
		fmt.Printf("Info: %s  Target: %s  Status %d\n", installCcResponse.Info, installCcResponse.Target, installCcResponse.Status)
	}
	fmt.Println("Chaincode installed")

	return channelId + chaincodeId + chaincodeVersion + "ok", nil
}
