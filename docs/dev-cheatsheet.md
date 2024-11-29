## **VitalEdge Blockchain Developer's Cheatsheet**

Blockchians are not for the faint of heart. To make the journey of working with Hyperledger Fabric, we have prepared a concise yet comprehensive list of all the handy commands you’ll need for setting up, managing, and developing on our VitalEdge Blockchain network. This cheatsheet ensures you have a quick reference for all critical steps and commands in your VitalEdge Blockchain project.

### **1. Environment Setup**
Add these to your `~/.bash_profile` for ease of use:
```bash
export PATH=/Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/bin:$PATH
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```
Reload the profile:
```bash
source ~/.bash_profile
```

---

### **2. Network Management**
#### **Start the Network**
```bash
./network.sh up createChannel -c vitaledgechannel -ca
```

#### **Stop the Network**
```bash
./network.sh down
```

#### **Restart the Network**
```bash
./network.sh down && ./network.sh up createChannel -c vitaledgechannel -ca
```

---

### **3. Docker Commands**
#### **Check Running Containers**
```bash
docker ps
```

#### **Inspect Logs**
```bash
docker logs <container_name>
```
Example:
```bash
docker logs peer0.org1.example.com
```

#### **Execute Commands Inside a Container**
```bash
docker exec -it <container_name> bash
```
Example:
```bash
docker exec -it peer0.org1.example.com bash
```

#### **Copy Files Between Host and Container**
Copy from host to container:
```bash
docker cp <host_path> <container_name>:<container_path>
```
Example:
```bash
docker cp organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/cacerts/localhost-9054-ca-orderer.pem peer0.org1.example.com:/etc/hyperledger/fabric-ca/ordererOrg/
```

Copy from container to host:
```bash
docker cp <container_name>:<container_path> <host_path>
```

---

### **4. Channel Management**
#### **List Channels**
```bash
docker exec -it peer0.org1.example.com peer channel list
```

#### **Fetch Channel Configuration**
```bash
docker exec -it peer0.org1.example.com peer channel fetch config /tmp/config_block.pb -o orderer.example.com:7050 -c vitaledgechannel --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem
```

Copy the block to your host:
```bash
docker cp peer0.org1.example.com:/tmp/config_block.pb .
```

Decode the configuration block:
```bash
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
```

Inspect the JSON:
```bash
jq '.' config_block.json
```

---

### **5. Chaincode Management**
#### **Deploy Example Chaincode**
```bash
./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/vitaledgechaincode -ccl go -c vitaledgechannel
```

#### **Verify Installed Chaincode**
```bash
docker exec -it peer0.org1.example.com peer lifecycle chaincode queryinstalled
```

#### **Approve Chaincode for Your Org**
```bash
docker exec -it peer0.org1.example.com peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem --channelID vitaledgechannel --name basic --version 1.0 --sequence 1
```

#### **Commit Chaincode**
```bash
docker exec -it peer0.org1.example.com peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem --channelID vitaledgechannel --name basic --version 1.0 --sequence 1
```

#### **Invoke Chaincode**
```bash
docker exec -it peer0.org1.example.com peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem -C vitaledgechannel -n basic -c '{"Args":["InitLedger"]}'
```

#### **Query Chaincode**
```bash
docker exec -it peer0.org1.example.com peer chaincode query -C vitaledgechannel -n basic -c '{"Args":["GetAllAssets"]}'
```

---

### **6. Housekeeping Commands**
#### **Fix Admin Certificate Issues**
1. Navigate to the admin MSP directory:
   ```bash
   cd /Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
   ```

2. Create the `admincerts` directory:
   ```bash
   mkdir -p admincerts
   ```

3. Copy the admin certificate:
   ```bash
   cp signcerts/cert.pem admincerts/
   ```

---

### **7. Debugging TLS Issues**
#### **Verify Peer-Orderer Communication**
Test TLS connectivity:
```bash
openssl s_client -connect orderer.example.com:7050 -CAfile /Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/cacerts/localhost-9054-ca-orderer.pem
```

---

### **8. General Utilities**
#### **Check Installed Binaries**
```bash
which configtxlator
which peer
which jq
```

#### **Test Binary Versions**
```bash
peer version
configtxlator version
jq --version
```

---

### **Future Extensions**
1. Commands for creating and deploying **custom chaincode**.
2. Steps for adding new organizations or peers to the network.
3. Integration with external systems for RESTful API interactions.

---
