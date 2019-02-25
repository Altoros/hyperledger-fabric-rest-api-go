package api

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/common/ccprovider"
	fabriccmn "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/msp"
	"github.com/pkg/errors"
)

func (fsc *FabricSdkClient) InstalledChaincodes() (string, error) {

	queryInstalledChaincodesResponse, err := fsc.admin.QueryInstalledChaincodes(resmgmt.WithTargets(fsc.GetCurrentPeer()))
	installedChaincodes := queryInstalledChaincodesResponse.GetChaincodes()

	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(installedChaincodes)
	return string(jsonBytes), nil
}

func (fsc *FabricSdkClient) InstantiatedChaincodes(channelId string) (string, error) {

	queryInstantiatedChaincodesResponse, err := fsc.admin.QueryInstantiatedChaincodes(channelId, resmgmt.WithTargets(fsc.GetCurrentPeer()))

	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(queryInstantiatedChaincodesResponse.Chaincodes)
	return string(jsonBytes), nil
}

func (fsc *FabricSdkClient) ChaincodeInfo(channelId, chainCodeId string) (string, error) {

	var args [][]byte
	args = append(args, []byte(channelId))
	args = append(args, []byte(chainCodeId))

	client, err := fsc.ChannelClient("chainhero")
	if err != nil {
		return "", fmt.Errorf("failed to create channel client: %v", err)
	}

	response, err := client.Query(
		channel.Request{ChaincodeID: "lscc", Fcn: "getccdata", Args: args},
		channel.WithTargetEndpoints(fsc.GetCurrentPeer().URL()),
	)
	if err != nil {
		return "", errors.Errorf("error querying for chaincode info: %v", err)
	}

	ccData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(response.Payload, ccData)
	if err != nil {
		return "", errors.Errorf("error unmarshalling chaincode data: %v", err)
	}

	cdsData := &ccprovider.CDSData{}
	unmarshalOrPanic(ccData.Data, cdsData)

	policy := &fabriccmn.SignaturePolicyEnvelope{}
	unmarshalOrPanic(ccData.Policy, policy)

	instPolicy := &fabriccmn.SignaturePolicyEnvelope{}
	unmarshalOrPanic(ccData.InstantiationPolicy, instPolicy)

	jsonObj := gabs.New()
	jsonObj.Set(ccData.Name, "name")
	jsonObj.Set(ccData.Version, "version")
	jsonObj.Set(ccData.Escc, "escc")
	jsonObj.Set(ccData.Vscc, "vscc")

	jsonObj.Set(SignaturePolicyEnvelopeToJsonObject(policy).Data(), "policy")

	jsonObj.SetP(cdsData.CodeHash, "data.CodeHash")
	jsonObj.SetP(cdsData.MetaDataHash, "data.MetaDataHash")

	jsonObj.Set(ccData.Id, "id")

	jsonObj.Set(SignaturePolicyEnvelopeToJsonObject(instPolicy).Data(), "instantiationPolicy")

	return jsonObj.String(), nil
}

// SignaturePolicyEnvelopeToJsonObject converts a SignaturePolicyEnvelope to json object
func SignaturePolicyEnvelopeToJsonObject(sigPolicyEnv *fabriccmn.SignaturePolicyEnvelope) *gabs.Container {

	jsonObj := gabs.New()

	jsonObj.SetP(SignaturePolicyToJsonObject(sigPolicyEnv.Rule).Data(), "rules")

	jsonObj.ArrayP("identities")

	for _, identity := range sigPolicyEnv.Identities {
		jsonObjIdentity := MSPPrincipalToJsonObject(identity)
		jsonObj.ArrayAppendP(jsonObjIdentity.Data(), "identities")
	}

	return jsonObj
}

// MSPPrincipalToJsonObject converts a MSPPrincipal to json object
func MSPPrincipalToJsonObject(principal *msp.MSPPrincipal) *gabs.Container {

	jsonObj := gabs.New()
	jsonObj.SetP(principal.PrincipalClassification, "PrincipalClassification")

	switch principal.PrincipalClassification {
	case msp.MSPPrincipal_ROLE:
		// Principal contains the msp role
		mspRole := &msp.MSPRole{}
		unmarshalOrPanic(principal.Principal, mspRole)
		jsonObj.SetP(mspRole.Role, "Role")
		jsonObj.SetP(MSPRoleTypeLookup[mspRole.Role], "RoleName")
		jsonObj.SetP(mspRole.MspIdentifier, "MspIdentifier")
	case msp.MSPPrincipal_IDENTITY:
		jsonObj.SetP(principal.Principal, "Identity")
	case msp.MSPPrincipal_ORGANIZATION_UNIT:
		// Principal contains the OrganizationUnit
		unit := &msp.OrganizationUnit{}
		unmarshalOrPanic(principal.Principal, unit)

		jsonObj.SetP(unit.MspIdentifier, "MspIdentifier")
		jsonObj.SetP(unit.OrganizationalUnitIdentifier, "OrganizationalUnitIdentifier")
		jsonObj.SetP(unit.CertifiersIdentifier, "CertifiersIdentifier")
	default:
		jsonObj.Set("unknown PrincipalClassification", "Error")
	}

	return jsonObj
}

var MSPRoleTypeLookup = map[msp.MSPRole_MSPRoleType]string{
	msp.MSPRole_MEMBER: `MSPRole_MEMBER`,
	msp.MSPRole_ADMIN:  `MSPRole_ADMIN`,
	msp.MSPRole_CLIENT: "MSPRole_CLIENT",
	msp.MSPRole_PEER:   "MSPRole_PEER",
}

// SignaturePolicyToJsonObject converts a SignaturePolicy policy to json object
func SignaturePolicyToJsonObject(sigPolicy *fabriccmn.SignaturePolicy) *gabs.Container {
	jsonObj := gabs.New()

	switch t := sigPolicy.Type.(type) {
	case *fabriccmn.SignaturePolicy_SignedBy:
		jsonObj.SetP("SignedBy", "Type")
		jsonObj.SetP(t.SignedBy, "SignedBy")
	case *fabriccmn.SignaturePolicy_NOutOf_:
		jsonObj.SetP("NOutOf", "Type")
		jsonObj.SetP(SignaturePolicyNOutOfToJsonObject(t.NOutOf).Data(), "NOutOf")
	default:
		jsonObj.Set(fmt.Sprintf("!!!!!!! Don't know how to print signature policy: %s", t), "Error")
	}

	return jsonObj
}

// SignaturePolicyNOutOfToJsonObject converts a SignaturePolicy_NOutOf policy to json object
func SignaturePolicyNOutOfToJsonObject(sigPolicy *fabriccmn.SignaturePolicy_NOutOf) *gabs.Container {

	jsonObj := gabs.New()

	jsonObj.SetP(sigPolicy.N, "N")

	jsonObj.Array("Rules")

	for _, policy := range sigPolicy.Rules {
		jsonObj.ArrayAppend(SignaturePolicyToJsonObject(policy).Data(), "Rules")
	}

	return jsonObj
}

func unmarshalOrPanic(buf []byte, pb proto.Message) {
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		panic(err)
	}
}
