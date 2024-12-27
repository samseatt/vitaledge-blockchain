

############ ONLY THE FIRST TIME ##################
cd /Users/samseatt/projects/vitaledge/vitaledge-blockchain

curl -sSL https://bit.ly/2ysbOFE | bash -s

cd fabric/vitaledge-network
mkdir -p ../chaincode/vitaledgechaincode
cp -r ../asset-transfer-basic/chaincode-go/* ../chaincode/vitaledgechaincode/
###################################################

cd /Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric/vitaledge-network

./network.sh down

# Start the blockchain with the channel (use one)
# ./network.sh createChannel -c vitaledgechannel
# ./network.sh up createChannel -c vitaledgechannel
./network.sh up createChannel -c vitaledgechannel -ca

# Just to verity the channel is up
docker exec -it peer0.clinicians.xmed.ai peer channel list
docker exec -it peer0.scientists.xnome.net peer channel list

# Deploy the chaincode
./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/vitaledgechaincode -ccl go -c vitaledgechannel

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1
# export CORE_PEER_TLS_ENABLED=true
# export CORE_PEER_LOCALMSPID=Org1MSP
# export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
# export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
# export CORE_PEER_ADDRESS=localhost:7051


# Environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/clinicians.xmed.ai/peers/peer0.clinicians.xmed.ai/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/clinicians.xmed.ai/users/Admin@clinicians.xmed.ai/msp
export CORE_PEER_ADDRESS=localhost:7051

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C vitaledgechannel -n vitaledgechaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.drx.network --tls --cafile "${PWD}/organizations/ordererOrganizations/drx.network/orderers/orderer.drx.network/msp/tlscacerts/tlsca.drx.network-cert.pem" -C vitaledgechannel -n vitaledgechaincode --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/clinicians.xmed.ai/peers/peer0.clinicians.xmed.ai/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/scientists.xnome.net/peers/peer0.scientists.xnome.net/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

peer chaincode query -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["GetAllAssets"]}'

# Environment variables for Org2
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org2MSP
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/scientists.xnome.net/peers/peer0.scientists.xnome.net/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/scientists.xnome.net/users/Admin@scientists.xnome.net/msp
export CORE_PEER_ADDRESS=localhost:9051

peer chaincode query -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["ReadAsset","asset6"]}'


curl -X POST http://localhost:8082/api/log-event \
     -H "Content-Type: application/json" \
     -d '{"event_id": "E1", "event_type": "Info", "event_details": "Test event", "user": "Admin", "timestamp": "2024-11-23T10:00:00Z"}'


# Get into running docker containers
docker exec -it vitaledge-rest-api sh
docker exec -it peer0.clinicians.xmed.ai sh

# Get logs of running docker conatainers
docker logs vitaledge-rest-api sh
docker logs peer0.clinicians.xmed.ai sh

# Use the openssl command to inspect the certificate
openssl x509 -in /path/to/cert.pem -text -noout


# Verify the Channel Configuration
# Ensure the Org1MSP definition is correctly applied to the channel configuration.
# If you suspect configuration drift, fetch the channel configuration and inspect it:
peer channel fetch config config_block.pb -o localhost:7050 -c vitaledgechannel --tls --cafile /path/to/orderer/ca.pem
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
configtxlator proto_decode --input config.block --type common.Block --output config.json



# ---- ONE TIME
fabric-ca-client register --id.name appUser --id.secret appUserpw --id.type client --tls.certfiles /path/to/ca-cert.pem
fabric-ca-client enroll -u https://appUser:appUserpw@localhost:7054 --caname ca-clinicians --csr.names "C=US,ST=North Carolina,O=Hyperledger,OU=client" --tls.certfiles /path/to/ca-cert.pem
