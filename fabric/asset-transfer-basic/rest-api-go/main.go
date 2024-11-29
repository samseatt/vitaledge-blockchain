package main

import (
	"fmt"
	"rest-api-go/web"
)

func main() {
	//Initialize setup for Org1
	cryptoPath := "../../test-network/organizations/peerOrganizations/clinicians.xmed.ai"
	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/User1@clinicians.xmed.ai/msp/signcerts/cert.pem",
		KeyPath:      cryptoPath + "/users/User1@clinicians.xmed.ai/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.clinicians.xmed.ai/tls/ca.crt",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.clinicians.xmed.ai",
	}

	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org1: ", err)
	}
	web.Serve(web.OrgSetup(*orgSetup))
}
