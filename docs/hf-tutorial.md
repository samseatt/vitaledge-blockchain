# Hyperledger Fabric Tutorial

This is a focused **Hyperledger Fabric Tutorial** to help understand its architecture, components, and how they interact, especially the parts we will use for **VitalEdge Blockchain**. Fabric has a steep learning curve, but breaking it down will help you navigate its complexities.

This is a brief primer to help familiarize with the general concepts that may be relevant to Vitaledge Blockchain. There are several other very helpful documentation interspersed withing this project or provided at Hypeerlogic Fabric site in the shape of references and focused tutorials.

---

## **What is Hyperledger Fabric?**
Hyperledger Fabric is a **permissioned blockchain framework** designed for enterprise use. Unlike public blockchains (e.g., Bitcoin, Ethereum), Fabric provides:
- **Modularity**: Customizable components like consensus, identity, and policies.
- **Privacy**: Channels and private data collections for selective data sharing.
- **Permissioned Access**: Only authorized participants can join and interact with the network.

---

## **Fabric Components**
Here’s a high-level overview of the key components and their roles:

### **1. Membership Service Provider (MSP)**
- **What It Does**: Manages identities and permissions in the network.
- **Key Files**:
  - **Certificates**: Identify entities (e.g., peers, orderers, admins).
  - **Policies**: Define who can read/write/approve transactions.
- **Directory**:
  - `organizations/peerOrganizations`: Stores peer MSPs.
  - `organizations/ordererOrganizations`: Stores orderer MSPs.
  
**Why it Matters**: MSPs enforce trust in the network by linking cryptographic identities to organizational roles (e.g., admin, peer).

---

### **2. Certificate Authority (CA)**
- **What It Does**: Issues and manages certificates for network participants.
- **Key Files**:
  - `ca-cert.pem`: The CA’s root certificate.
  - Private Key (`*_sk`): Used to sign certificates.
- **Directory**:
  - `organizations/fabric-ca`: Stores CA data for peers and orderers.

**Why it Matters**: Every participant (peers, orderers, clients) must have a valid certificate issued by the CA.

---

### **3. Peers**
- **What They Do**:
  - Hold copies of the ledger.
  - Execute and validate smart contracts (chaincode).
- **Key Directories**:
  - `organizations/peerOrganizations/org1.example.com/peers`: Peer-specific MSPs and TLS certificates.

**Types of Peers**:
- **Endorsing Peer**: Simulates and endorses transactions.
- **Committing Peer**: Commits validated transactions to the ledger.
- **Anchor Peer**: Connects organizations in a channel.

---

### **4. Ordering Service**
- **What It Does**:
  - Sequences transactions into blocks.
  - Ensures consistent ledger state across peers.
- **Key Files**:
  - `organizations/ordererOrganizations`: MSP for orderers.
- **Consensus Options**:
  - **Raft** (default): Leader-based consensus.
  - **Kafka**: External message broker.

---

### **5. Ledger**
- **What It Does**: Stores the blockchain data.
  - **World State**: Current state of the ledger.
  - **Blockchain**: Immutable log of transactions.
- **Where It Lives**: On each peer node.

---

### **6. Chaincode**
- **What It Does**: Smart contracts that define transaction logic.
- **Key Files**:
  - Go, Node.js, or Java source files.
  - Example: `assetTransfer.go` (manages assets in the example project).
- **Where It Lives**:
  - `chaincode` directories for your application logic.

---

### **7. Channels**
- **What They Do**: Partition the blockchain into separate ledgers.
- **Key Files**:
  - `config_block.pb`: Channel configuration block.
  - Policies in `config_block.json`.
- **Directory**:
  - Each channel has its own ledger.

---

### **8. Policies**
- **What They Do**:
  - Define access control for the network.
  - Example: Who can invoke chaincode, add blocks, or update configurations.
- **Defined In**:
  - Channel configuration (`Admins`, `Writers`, etc.).
  - MSP policies (e.g., NodeOUs).

---

## **How Components Work Together**
### **1. Initial Setup**
1. Use the CA to issue certificates for:
   - Admins
   - Peers
   - Orderers
2. Create an MSP directory structure with:
   - `admincerts`
   - `cacerts`
   - `keystore`
   - `signcerts`

---

### **2. Network Initialization**
1. Start the Fabric network (peers, orderers, and CAs).
2. Create a channel to define the participants and policies.
3. Join peers to the channel.

---

### **3. Deploy Chaincode**
1. Package and install the chaincode on each peer.
2. Approve the chaincode definition from each organization.
3. Commit the chaincode to the channel.

---

### **4. Invoke Transactions**
1. A client application submits a transaction proposal.
2. Endorsing peers simulate the proposal and generate endorsements.
3. The ordering service creates a block of transactions.
4. Committing peers validate and add the block to their ledger.

---

## **Directory Walkthrough for Example Network**
Let’s map the directories in the `test-network` to Fabric concepts:

| **Directory**                                | **Purpose**                                                                 |
|----------------------------------------------|-----------------------------------------------------------------------------|
| `fabric-samples/test-network/organizations`  | Stores MSPs, CAs, and related certificates for peers and orderers.         |
| `fabric-samples/test-network/chaincode`      | Chaincode source files for smart contracts.                                |
| `fabric-samples/test-network/docker-compose` | Docker configurations for peers, orderers, and CA servers.                |
| `fabric-samples/test-network/scripts`        | Helper scripts (`network.sh`) to start, stop, and configure the network.   |

---

## **Hands-On: Testing Fabric**
### **1. Check the Network Status**
Verify the running containers:
```bash
docker ps
```

Check logs for a peer or orderer:
```bash
docker logs peer0.org1.example.com
docker logs orderer.example.com
```

---

### **2. Query the Ledger**
Fetch the blockchain info:
```bash
docker exec -it peer0.org1.example.com peer channel getinfo -c vitaledgechannel
```

Query the chaincode:
```bash
docker exec -it peer0.org1.example.com peer chaincode query -C vitaledgechannel -n basic -c '{"Args":["QueryAllAssets"]}'
```

---

### **3. Invoke Transactions**
Invoke a chaincode function:
```bash
docker exec -it peer0.org1.example.com peer chaincode invoke -C vitaledgechannel -n basic -c '{"Args":["CreateAsset","asset1","blue","10","tom","100"]}' -o orderer.example.com:7050 --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem
```
