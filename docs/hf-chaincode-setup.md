
## **Hyperledger Fabric Setup and Chaincode Deployment**

The following can serve as a guide to set up Hyperledger Fabric, create a private blockchain, and deploy the example chaincode to `vitaledgechannel`. This is tailored to our setup and includes all the necessary steps, file edits, environment variables, and options to avoid pitfalls. 

---
### **1. Install Prerequisites**
Ensure the following are installed on your machine:
- **Docker and Docker Compose** (Ensure Docker is running):
  ```bash
  docker --version
  docker-compose --version
  ```

- **Go** (Version 1.20 or higher):
  ```bash
  go version
  ```
  If Go is not installed, download and install from [https://go.dev/dl/](https://go.dev/dl/).

---

### **2. Set Up Fabric Samples**
1. Clone the `fabric-samples` repository:
   ```bash
   cd /Users/samseatt/projects/vitaledge/vitaledge-blockchain
   curl -sSL https://bit.ly/2ysbOFE | bash -s
   ```

2. Navigate to the `test-network` directory:
   ```bash
   cd fabric-samples/test-network
   ```

3. Ensure the `bin` directory (Fabric binaries) is in your `PATH`. Add the following to your `~/.bash_profile`:
   ```bash
   export PATH=/Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/bin:$PATH
   ```
   Reload your shell:
   ```bash
   source ~/.bash_profile
   ```

---

### **3. Launch the Network**
1. Start the test network and create the `vitaledgechannel`:
   ```bash
   ./network.sh up createChannel -c vitaledgechannel -ca
   ```

2. Verify the network is running:
   ```bash
   docker ps
   ```
   You should see containers for the orderer, CA, and peer nodes.

---

### **4. Fix Admin Certificate Issues**
To ensure proper MSP configuration for admin privileges, edit the files and directories as follows:

1. Navigate to the `Admin@org1.example.com` MSP directory:
   ```bash
   cd /Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
   ```

2. Create the `admincerts` directory (if missing):
   ```bash
   mkdir -p admincerts
   ```

3. Copy the admin certificate:
   ```bash
   cp signcerts/cert.pem admincerts/
   ```

4. Ensure `config.yaml` has `NodeOUs` enabled:
   ```yaml
   NodeOUs:
     Enable: true
     ClientOUIdentifier:
       Certificate: cacerts/localhost-7054-ca-org1.pem
       OrganizationalUnitIdentifier: client
     PeerOUIdentifier:
       Certificate: cacerts/localhost-7054-ca-org1.pem
       OrganizationalUnitIdentifier: peer
     AdminOUIdentifier:
       Certificate: cacerts/localhost-7054-ca-org1.pem
       OrganizationalUnitIdentifier: admin
     OrdererOUIdentifier:
       Certificate: cacerts/localhost-7054-ca-org1.pem
       OrganizationalUnitIdentifier: orderer
   ```

---

### **5. Prepare the Orderer CA Certificate**
Ensure the orderer CA certificate is accessible by the peer for TLS communication.

1. Copy the orderer CA certificate to the peer container:
   ```bash
   docker exec -it peer0.org1.example.com mkdir -p /etc/hyperledger/fabric-ca/ordererOrg
   docker cp /Users/samseatt/projects/vitaledge/vitaledge-blockchain/fabric-samples/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/cacerts/localhost-9054-ca-orderer.pem peer0.org1.example.com:/etc/hyperledger/fabric-ca/ordererOrg/
   ```

2. Confirm the file is copied:
   ```bash
   docker exec -it peer0.org1.example.com ls /etc/hyperledger/fabric-ca/ordererOrg/
   ```

---

### **6. Fetch and Verify Channel Configuration**
1. Fetch the channel configuration block:
   ```bash
   docker exec -it peer0.org1.example.com peer channel fetch config /tmp/config_block.pb -o orderer.example.com:7050 -c vitaledgechannel --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem
   ```

2. Copy the block to your host machine:
   ```bash
   docker cp peer0.org1.example.com:/tmp/config_block.pb .
   ```

3. Decode the block with `configtxlator`:
   ```bash
   configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
   ```

4. Verify the decoded configuration for:
   - `OrdererAddresses`
   - `MSP Definitions`
   - `Policies`

---

### **7. Deploy Example Chaincode**
1. Deploy the `asset-transfer-basic` example chaincode:
   ```bash
   ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go/ -ccl go -c vitaledgechannel
   ```

2. Verify the chaincode installation:
   ```bash
   docker exec -it peer0.org1.example.com peer lifecycle chaincode queryinstalled
   ```

3. Approve the chaincode definition for your organization:
   ```bash
   docker exec -it peer0.org1.example.com peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem --channelID vitaledgechannel --name basic --version 1.0 --sequence 1
   ```

4. Commit the chaincode:
   ```bash
   docker exec -it peer0.org1.example.com peer lifecycle chaincode commit -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem --channelID vitaledgechannel --name basic --version 1.0 --sequence 1
   ```

5. Invoke and query the chaincode to test it:
   - **Invoke**:
     ```bash
     docker exec -it peer0.org1.example.com peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem -C vitaledgechannel -n basic -c '{"Args":["InitLedger"]}'
     ```
   - **Query**:
     ```bash
     docker exec -it peer0.org1.example.com peer chaincode query -C vitaledgechannel -n basic -c '{"Args":["GetAllAssets"]}'
     ```

---

### **8. Develop VitalEdge Blockchain Customizations and Features**
Once the example chaincode is deployed successfully using these steps, I'm using it to:
- Build and deploy custom chaincode for the VitalEdge Blockchain project.
- Updating the chaincode with additional functionality (e.g., tokens, proof of participation).

---
