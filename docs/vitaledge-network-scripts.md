### **Hyperledger Fabric Test Network Scripts: A Detailed Breakdown**

This document explains the **`network.sh`** script, its arguments, and other related scripts/files in the **vitaledge-network** directory. It will help you understand how to manage the Hyperledger Fabric based vitaledge network and its components. This should cover most of what you need to understand vitaledge network and its supporting files.

---

## **1. `network.sh` Script**
The `network.sh` script is the primary utility for setting up, managing, and interacting with the vitaledge network. It automates common tasks like creating channels, deploying chaincode, and tearing down the network.

### **Arguments and Options**
Here’s a detailed breakdown of the arguments and options supported by `network.sh`:

| **Argument**          | **Purpose**                                                                                          |
|------------------------|------------------------------------------------------------------------------------------------------|
| `up`                  | Starts the Fabric network by creating containers for peers, orderers, and certificate authorities.   |
| `down`                | Stops and removes all Fabric-related containers, networks, and artifacts.                            |
| `restart`             | Restarts the network without removing artifacts (like channels and chaincode definitions).           |
| `createChannel`       | Creates a new channel.                                                                               |
| `deployCC`            | Deploys chaincode to the network.                                                                    |
| `-ca`                 | Enables the use of Certificate Authorities (CAs) for identity and MSP management.                   |
| `-c <channel name>`   | Specifies the name of the channel to create or use (default: `mychannel`).                           |
| `-ccn <chaincode name>` | Specifies the name of the chaincode to deploy.                                                      |
| `-ccp <chaincode path>` | Specifies the path to the chaincode source files.                                                   |
| `-ccl <chaincode language>` | Specifies the chaincode language (e.g., `go`, `java`, `node`).                                     |
| `-ccv <version>`      | Specifies the version of the chaincode (default: `1.0`).                                             |
| `-ccs <sequence>`     | Specifies the chaincode sequence (default: `1`).                                                     |
| `-verbose`            | Enables verbose logging for debugging purposes.                                                      |

---

### **Key Workflow Examples**

#### **1. Start the Network**
```bash
./network.sh up createChannel -ca -c vitaledgechannel
```
- **Purpose**: Starts the network, initializes CAs, and creates the `vitaledgechannel` channel.
- **Options**:
  - `-ca`: Use certificate authorities.
  - `-c vitaledgechannel`: Specify the channel name as `vitaledgechannel`.

---

#### **2. Deploy Chaincode**
```bash
./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/vitaledgechaincode -ccl go -c vitaledgechannel
```
- **Purpose**: Deploys the vitaledge chaincode `vitaledgechaincode` to the channel.
- **Options**:
  - `-ccn vitaledgechaincode`: Name the chaincode `vitaledgechaincode`.
  - `-ccp ../chaincode/vitaledgechaincode`: Specify the path to this Go-based chaincode.
  - `-ccl go`: Use the Go programming language for vitaledgechaincode.

---

#### **3. Stop the Network**
```bash
./network.sh down
```
- **Purpose**: Stops all Fabric containers and cleans up the generated artifacts.

---

### **How `network.sh` Works**
Internally, the `network.sh` script calls various other scripts and tools to accomplish tasks.

#### **Workflow for `up`**
1. **Start CA Servers**:
   - Uses `docker-compose-ca.yaml` to spin up CA containers for Org1 and Org2.
2. **Generate MSP Artifacts**:
   - Calls `registerEnroll.sh` to register identities and generate certificates for peers and orderers.
3. **Start the Network**:
   - Uses `docker-compose-test-net.yaml` to start peer and orderer containers.
4. **Create Channels**:
   - Calls `createChannel.sh` to create the default or specified channel.

#### **Workflow for `deployCC`**
1. **Package Chaincode**:
   - Packages the chaincode into a `.tar.gz` file.
2. **Install Chaincode**:
   - Installs the chaincode on all peers using `peer lifecycle chaincode install`.
3. **Approve Chaincode**:
   - Each organization approves the chaincode definition.
4. **Commit Chaincode**:
   - Commits the chaincode to the channel.

---

## **2. Other Scripts**

### **a. `registerEnroll.sh`**
This script handles identity registration and enrollment with the Fabric CA. It is called automatically by `network.sh` during the `up` process.

#### **Key Functions**
- **`createOrg1`**:
  - Registers and enrolls Org1’s admin, peers, and users.
  - Generates MSP structure for Org1.
- **`createOrg2`**:
  - Similar to `createOrg1` but for Org2.
- **`createOrderer`**:
  - Handles the registration and enrollment of the orderer identity.

---

### **b. `createChannel.sh`**
This script creates and joins a channel for the network.

#### **Key Steps**
1. **Fetch Channel Configuration**:
   - Uses `peer channel fetch` to retrieve the configuration block.
2. **Create the Channel**:
   - Uses `peer channel create` with the specified channel name and configuration file.
3. **Join Peers to the Channel**:
   - Executes `peer channel join` for each peer to join the channel.

---

## **3. Key Docker Compose Files**

### **a. `docker-compose-test-net.yaml`**
Defines the Docker containers for the test network, including:
- **Orderers**:
  - `orderer.example.com`: The ordering node.
- **Peers**:
  - `peer0.org1.example.com`
  - `peer0.org2.example.com`

Key services:
```yaml
services:
  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    ...
  orderer.example.com:
    container_name: orderer.example.com
    ...
```

---

### **b. `docker-compose-ca.yaml`**
Defines the CA containers for Org1 and Org2.

Key services:
```yaml
services:
  ca_org1:
    container_name: ca_org1
    ...
  ca_org2:
    container_name: ca_org2
    ...
```

---

## **4. Key Configuration Files**

### **a. `configtx.yaml`**
This file defines the channel configuration, including:
- **Organizations**:
  - MSPs for Org1 and Org2.
- **Policies**:
  - Admin, writer, and reader policies.
- **Capabilities**:
  - Defines supported Fabric features.

---

### **b. `core.yaml`**
Defines the peer configuration, such as:
- **Peer Address**:
  ```yaml
  peer:
    address: peer0.org1.example.com:7051
  ```
- **Gossip Protocol**:
  - Handles peer-to-peer communication.

---

## **Summary**

| **Script/File**                  | **Purpose**                                                                                       |
|----------------------------------|---------------------------------------------------------------------------------------------------|
| `network.sh`                     | Main script for managing the test network.                                                       |
| `registerEnroll.sh`              | Registers and enrolls identities with the CA.                                                    |
| `createChannel.sh`               | Creates and joins a channel.                                                                     |
| `docker-compose-test-net.yaml`   | Defines Docker containers for peers, orderers, and other services.                               |
| `docker-compose-ca.yaml`         | Defines CA containers.                                                                           |
| `configtx.yaml`                  | Defines channel configuration and policies.                                                      |
| `core.yaml`                      | Defines peer-specific configurations.                                                            |
