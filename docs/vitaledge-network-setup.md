
# HyperLedger Fabric - VitalEdge Network

This document provides a detailed breakdown of **"vitaledge-network"** which is based on Hyperledger Fabric. This example is designed to help developers set up a local vitaledge network quickly, providing a working baseline for understanding its current implementation and developing new features. We will examine its architecture, directory structure, and how everything fits together.

The vitaledge-network structure is modular and closely follows and acknowledges the patterns developed by Hyperlogic Fabric architects. This includes how cryptographic material is stored, how channels are created, and how chaincode is deployed. Understanding that will help you effectively work with this Fabric-based blockchain.

## **Overview of the VitalEdge Network**
The **vitaledge-network** provides:
1. A **single-channel** network contaiinng a custom channel named `vitaledgechannel`.
2. Two organizations:
   - **Org1** (Clinicians) and **Org2** (Data Scientists), each with:
     - A CA (Certificate Authority)
     - A peer node (peer0.clinicians.xmed.ai, peer0.data-scientists.xnome.net)
   - Admin identities and MSPs for both organizations.
3. An **orderer** node (represented by a decentralized doctors' network - drx.network) for transaction ordering.
4. Scripts for managing the network (`network.sh`) and deploying vitaledge chaincode.

---

## **Test Network Architecture**
### **1. Nodes**
- **Orderer Node**:
  - Sequences transactions into blocks and distributes them to peers.
  - Single node: `orderer.drx.network`.
- **Peer Nodes**:
  - Store the ledger and execute chaincode.
  - Two peers: `peer0.clinicians.xmed.ai`, `peer0.data-scientists.xnome.net`.

### **2. MSPs**
- Each organization has an MSP, which defines its cryptographic identity.
- Includes admin, peer, and user certificates.

### **3. Certificate Authorities**
- Fabric CA servers issue certificates for peers, orderers, and admins.
- Two CAs:
  - **Org1 CA**: Issues certificates for Org1 participants.
  - **Org2 CA**: Issues certificates for Org2 participants.

### **4. Channel**
- VitalEdge channel: `vitaledgechannel`.
- Configured to include Org1 and Org2 with specified policies.

---

## **Directory Structure**
Here’s a walkthrough of the **vitaledge-network** directory structure:

### **Root Directory**
```plaintext
fabric/vitaledge-network/
```
This is the main directory containing scripts, configuration files, and subdirectories for running the vitaledge network.

| **File/Directory**        | **Purpose**                                                                                          |
|----------------------------|------------------------------------------------------------------------------------------------------|
| `network.sh`               | Main script to start, stop, and manage the network.                                                 |
| `docker-compose-test-net.yaml` | Docker Compose file defining the peer, orderer, and CA containers.                                 |
| `docker-compose-ca.yaml`   | Docker Compose file for the Certificate Authorities.                                                |
| `scripts/`                 | Scripts for setting up organizations, channels, and chaincode.                                      |
| `organizations/`           | Contains the cryptographic material for peers, orderers, and users.                                |
| `channel-artifacts/`       | Stores the configuration blocks and other channel-related files.                                    |
| `chaincode/`               | Directory for chaincode source files (e.g., `asset-transfer-basic/chaincode-go`).                   |

---

### **Key Subdirectories**
#### **1. `organizations/`**
Contains all cryptographic material, MSPs, and CA configurations.

| **Path**                                               | **Purpose**                                                                                   |
|--------------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `organizations/peerOrganizations/`                    | Cryptographic material and MSPs for peer nodes.                                               |
| `organizations/ordererOrganizations/`                 | Cryptographic material and MSPs for the orderer.                                              |
| `organizations/fabric-ca/org1/`                       | Certificate Authority server files for Org1 (includes `fabric-ca-server-config.yaml`).        |
| `organizations/fabric-ca/org2/`                       | Certificate Authority server files for Org2.                                                 |

**Important MSP Subdirectories**:
- `admincerts/`: Stores admin certificates.
- `cacerts/`: Stores the root CA certificate.
- `keystore/`: Stores the private keys.
- `signcerts/`: Stores the public certificates.

#### **2. `scripts/`**
Contains utility scripts for setting up organizations, channels, and chaincode.

| **Script**            | **Purpose**                                                                                             |
|------------------------|---------------------------------------------------------------------------------------------------------|
| `createChannel.sh`    | Script to create and join channels.                                                                     |
| `deployCC.sh`         | Deploys chaincode to the network.                                                                       |
| `registerEnroll.sh`   | Handles registration and enrollment of identities (admin, peers, users) with the Fabric CA.             |

#### **3. `channel-artifacts/`**
This directory contains files generated during channel creation.

| **File**                 | **Purpose**                                                                                          |
|--------------------------|------------------------------------------------------------------------------------------------------|
| `genesis.block`          | Genesis block for the orderer.                                                                       |
| `mychannel.block`        | Channel configuration block.                                                                         |
| `mychannel.tx`           | Channel transaction file.                                                                            |

#### **4. `chaincode/`**
Contains chaincode implementations for smart contracts.

| **Subdirectory**                   | **Purpose**                                                                                 |
|------------------------------------|---------------------------------------------------------------------------------------------|
| `asset-transfer-basic/chaincode-go/` | Chaincode written in Go for basic asset transfer functionality (used here).                              |
| `asset-transfer-basic/chaincode-java/` | Chaincode written in Java (we don't use this).                                      |
| `asset-transfer-basic/chaincode-node/` | Chaincode written in Node.js (we don't use this).                                   |

---

## **How Test Network Components Interact**
1. **MSP Creation**:
   - CAs issue certificates, and `registerEnroll.sh` registers identities (admin, peers, orderers).
   - Certificates are placed in the `organizations/` directory.

2. **Channel Creation**:
   - `createChannel.sh` creates a channel (e.g., `vitaledgechannel`) using the `configtxgen` tool.
   - The channel configuration is stored in `channel-artifacts/`.

3. **Joining Peers to the Channel**:
   - Peers use their MSPs to authenticate and join the channel.

4. **Deploying Chaincode**:
   - Chaincode is installed on all peers, approved by the organizations, and committed to the channel.

5. **Executing Transactions**:
   - A client submits a transaction proposal.
   - Endorsing peers simulate and sign the transaction.
   - The orderer creates a block, which is validated and added to the ledger by all peers.

---

## **Key Commands for Managing VitaEdge Network**
Here are some commonly used commands:

### **Starting the Network**
```bash
./network.sh up createChannel -c vitaledgechannel -ca
```

### **Deploying VitalEdge Chaincode**
```bash
./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/vitaledgechaincode -ccl go -c vitaledgechannel
```

### **Querying the Ledger**
```bash
docker exec -it peer0.clinicians.xmed.ai peer chaincode query -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["ReadAsset","asset1"]}'

peer chaincode query -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["GetAllAssets"]}'
```

### **Invoking a Transaction**
```bash
docker exec -it peer0.clinicians.xmed.ai peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/ca-cert.pem -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["TransferAsset","asset1","newOwner"]}'

peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/ca-cert.pem -C vitaledgechannel -n vitaledgechaincode -c '{"Args":["TransferAsset","asset1","newOwner"]}'
```
