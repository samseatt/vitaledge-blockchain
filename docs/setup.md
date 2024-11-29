### VitalEdge Blockchain Setup

**`TODO`**:
- This document is to be merged with the more comprehensive setup/installation document for vitaledgeblockchain.
- Hyperlogic Fabric framework installation should be discussed separately - or omitted altogether to avoid confusion and overwriting.

#### **1. Set Up the Development Environment**
- **Install Golang**:
  1. Download the Go binary for macOS: [Golang Downloads](https://golang.org/dl/).
  2. Verify installation:
     ```bash
     go version
     ```
  3. Set up environment variables in `~/.bash_profile`:
     ```bash
     export GOPATH=$HOME/go
     export PATH=$PATH:$GOPATH/bin
     ```
     Apply changes:
     ```bash
     source ~/.bash_profile
     ```

- **Install IDE**:
  - Recommended: **Visual Studio Code (VS Code)**.
  - Install the Go extension:
    ```bash
    code --install-extension golang.Go
    ```

#### **2. Prepare Dependencies**
- **Docker**:
  1. Install Docker Desktop for macOS: [Docker](https://www.docker.com/products/docker-desktop/).
  2. Verify installation:
     ```bash
     docker --version
     ```

- **Hyperledger Fabric Binaries**:
  1. Download and install. (Only do it once. However, if you're getting your vitaledge-blockchain project from GitHub, the you don't need that step. You may still decide to install it at alternate location to learn form or reuse its many example projects as guides to developing new features and expansions to the vitaledgeblockchain fabric:
     ```bash
     curl -sSL https://bit.ly/2ysbOFE | bash -s
     cd fabric-samples/test-network
     ```

#### **3. Initial Blockchain Setup**
0. jq installation (if not installed on your system):
   ```bash
   brew --version
   brew install jq
   jq --version
   ```
1. Launch the vitaledge network:
   ```bash
   cd fabric/vitaledge-network
   ./network.sh up createChannel -c vitaledgechannel -ca
   ```
2. Deploy vitaledge chaincode:
   ```bash
   ./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/vitaledgechaincode -ccl go -c vitaledgechannel
   ```

