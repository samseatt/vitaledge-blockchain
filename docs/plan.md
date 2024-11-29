### Comprehensive Plan for the VitalEdge Blockchain Implementation

**Plan** for developing the private blockchain (VitalEdge Blockchain) using **Hyperledger Fabric** framework and wrapping it in a Go-based microservice (VitalEdge Blockchain API):

---

### **Development Blueprint**

#### **1. Framework Selection**
- **Primary Framework**: **Hyperledger Fabric** (best suited for permissioned networks with modularity and enterprise-grade support).
- **Why**: High throughput, privacy-preserving channels, and compatibility with Go, fitting VitalEdge's requirements.
- **Public Blockchain**: Ethereum (to periodically commit hashes for external validation).

---

#### **2. Core Blockchain Setup**

##### **2.1 Installation & Network Configuration**
1. **Install Pre-Requisites**:
   - Docker and Docker Compose for containerization.
   - Hyperledger Fabric binaries and configuration tools.
2. **Launch Test Network**:
   - Create and start the Fabric test network:
     ```bash
     curl -sSL https://bit.ly/2ysbOFE | bash -s
     cd fabric-samples/test-network
     ./network.sh up createChannel -c vitaledgechannel -ca
     ```
   - Deploy the chaincode (`vitaledgechaincode`):
     ```bash
     ./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/ -ccl go
     ```

##### **2.2 Node and Channel Management**
- **Orderer Node**: Central node for transaction sequencing.
- **Peer Nodes**: Responsible for validating and committing transactions. Hosted in secure environments.
- **Channel**: `vitaledgechannel` for sensitive logging, with strict access controls.

##### **2.3 Chaincode Development**
- Create a **Go-based chaincode** to log actions securely.
- **Example Chaincode**:

```go
package main

import (
	"encoding/json"
	"time"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Log struct {
	EventID     string `json:"event_id"`
	Timestamp   string `json:"timestamp"`
	ActionType  string `json:"action_type"`
	ActorID     string `json:"actor_id"`
	Description string `json:"description"`
}

// Log an event
func (s *SmartContract) LogEvent(ctx contractapi.TransactionContextInterface, eventID, actionType, actorID, description string) error {
	timestamp := time.Now().Format(time.RFC3339)
	log := Log{
		EventID:     eventID,
		Timestamp:   timestamp,
		ActionType:  actionType,
		ActorID:     actorID,
		Description: description,
	}
	logAsBytes, _ := json.Marshal(log)
	return ctx.GetStub().PutState(eventID, logAsBytes)
}

// Query logs by event ID
func (s *SmartContract) QueryLog(ctx contractapi.TransactionContextInterface, eventID string) (*Log, error) {
	logAsBytes, err := ctx.GetStub().GetState(eventID)
	if err != nil || logAsBytes == nil {
		return nil, fmt.Errorf("Log not found: %s", eventID)
	}
	var log Log
	_ = json.Unmarshal(logAsBytes, &log)
	return &log, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(err)
	}
	if err := chaincode.Start(); err != nil {
		panic(err)
	}
}
```

---

#### **3. REST API Development**
Wrap the private blockchain functionality in a **Go-based microservice** using **Gin** or **Echo**.

##### **3.1 API Endpoints**
1. **Logging API** (`/log`):
   - **Method**: `POST`
   - Payload:
     ```json
     {
       "event_id": "123",
       "action_type": "data_ingress",
       "actor_id": "hashed_patient_id",
       "description": "PharmGKB annotations uploaded"
     }
     ```
   - **Functionality**: Logs events to the blockchain via chaincode.
   
2. **Query API** (`/query`):
   - **Method**: `GET`
   - Parameters: `event_id`
   - **Functionality**: Retrieves logs by event ID.

3. **Commit Hash API** (`/commit`):
   - **Method**: `POST`
   - **Functionality**: Commits the private blockchain hash to Ethereum.

##### **3.2 Example API Code**
```go
package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Log struct {
	EventID     string `json:"event_id"`
	ActionType  string `json:"action_type"`
	ActorID     string `json:"actor_id"`
	Description string `json:"description"`
}

func logEvent(c *gin.Context) {
	var log Log
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Call chaincode to log the event (stubbed for example)
	c.JSON(http.StatusOK, gin.H{"message": "Log successfully added"})
}

func queryLog(c *gin.Context) {
	eventID := c.Query("event_id")
	// Query chaincode for the log (stubbed for example)
	c.JSON(http.StatusOK, gin.H{"event": eventID, "message": "Log retrieved"})
}

func main() {
	r := gin.Default()
	r.POST("/log", logEvent)
	r.GET("/query", queryLog)
	r.Run(":8080")
}
```

---

#### **4. Public Blockchain Integration**
Commit hashes of the private blockchain to Ethereum for external validation.

- Use **web3.py** in a Python script or **web3.js** in a lightweight service.
- **Example Python Code**:
```python
from web3 import Web3

web3 = Web3(Web3.HTTPProvider('https://mainnet.infura.io/v3/YOUR_PROJECT_ID'))
private_key = "YOUR_PRIVATE_KEY"
account = web3.eth.account.from_key(private_key)

tx = {
    'to': '0xPublicBlockchainAddress',
    'value': 0,
    'gas': 2000000,
    'gasPrice': web3.toWei('50', 'gwei'),
    'nonce': web3.eth.getTransactionCount(account.address),
    'data': web3.toBytes(text="Hash_of_Private_Blockchain")
}
signed_tx = web3.eth.account.sign_transaction(tx, private_key)
tx_hash = web3.eth.sendRawTransaction(signed_tx.rawTransaction)
print(f"Transaction Hash: {web3.toHex(tx_hash)}")
```

---

#### **5. Deployment and Maintenance**
1. **Deployment**:
   - Use Docker Compose to deploy the blockchain and APIs.
   - For production, use Kubernetes to manage multiple nodes.

2. **Monitoring**:
   - Implement Prometheus and Grafana to monitor node performance and API usage.

3. **Archiving**:
   - Regularly archive old blockchain data for long-term storage.

---

### Consistency Review and Alignment with Earlier Plans

The **Design Document** is consistent with the vision and ideation previously discussed. It expands on the earlier concepts with more granular implementation details while retaining the original objectives. Below is a summary of the key points and alignment:

---

### **Key Consistencies**

1. **Core Objectives**:
   - **Immutable Logging**: Both documents emphasize the need for tamper-proof logs.
   - **Traceability and Auditability**: Providing an auditable trail for regulatory compliance is central to both plans.
   - **Integration-Friendly Interface**: REST API endpoints align well with the goal of seamless integration across microservices.

2. **Blockchain Structure and Implementation**:
   - **Block Design**:
     - The design document elaborates on the `Block` structure with fields like `Index`, `Timestamp`, `Data`, `PrevHash`, and `Hash`, consistent with the earlier ideation.
   - **Blockchain Validation**:
     - Functions for creating blocks, hashing, and validating the chain directly support the project's objectives.
   - **Consensus Algorithm**:
     - Proof of Authority (PoA) was suggested earlier; the simplified implementation supports VitalEdge's private network use case.

3. **REST API**:
   - Endpoints for logging (`/log-event`), querying (`/query-event`), and verifying (`/verify-data`) are aligned with previously described APIs.

4. **Security Considerations**:
   - The emphasis on TLS encryption, access control, and tamper detection is consistent with healthcare's regulatory needs.

5. **Implementation Phases**:
   - Phases for blockchain development, REST API setup, persistence, and optional admin dashboard align with the iterative approach previously suggested.

---

### **New or Expanded Details**
The design document introduces additional specifics that further clarify the implementation:
1. **Persistence Layer**:
   - Suggests **BoltDB** or **LevelDB** for lightweight and efficient data storage, a helpful addition for ensuring durability.
2. **Go Libraries**:
   - Recommendations like **crypto/sha256** for hashing, **Gin** or **Fiber** for REST API, and **Zap/Logrus** for structured logging enrich the technical roadmap.
3. **Deployment**:
   - Dockerized setup with orchestration options (e.g., AWS ECS or Kubernetes) provides clarity on deployment strategies.
4. **Admin Dashboard**:
   - Introduces the possibility of an optional web UI for administrators to interact with the blockchain data.

---

### **Next Steps for Setup**

To begin development on your Intel Mac (macOS Ventura) using **bash**, follow these steps:

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
  1. Download and install:
     ```bash
     curl -sSL https://bit.ly/2ysbOFE | bash -s
     cd fabric-samples/test-network
     ```

#### **3. Initial Blockchain Setup**
1. Launch the test network:
   ```bash
   ./network.sh up createChannel -c vitaledgechannel -ca
   ```
2. Deploy example chaincode:
   ```bash
   ./network.sh deployCC -ccn vitaledgechaincode -ccp ../chaincode/ -ccl go
   ```

---

### Creating the Blockchain Core in Go

I designed the core blockchain features, including block structure, hashing, validation, and the blockchain itself. This initial implementation focuses on modularity and scalability, ensuring future extensibility and alignment with VitalEdge requirements.

---

### **Plan**
1. **Block Structure**:
   - Defines the core components of a block (e.g., `Index`, `Timestamp`, `Data`, `PrevHash`, and `Hash`).
2. **Blockchain Structure**:
   - A chain of blocks with helper functions to append new blocks and validate the chain.
3. **Core Functions**:
   - `calculateHash`: Computes the hash of a block.
   - `createBlock`: Generates a new block.
   - `validateChain`: Ensures the integrity of the blockchain.

---

### **Go Implementation**

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index     int         `json:"index"`      // Position of the block in the chain
	Timestamp string      `json:"timestamp"`  // Time the block was created
	Data      interface{} `json:"data"`       // Data stored in the block
	PrevHash  string      `json:"prev_hash"`  // Hash of the previous block
	Hash      string      `json:"hash"`       // Hash of the current block
}

// Blockchain represents the full chain
type Blockchain struct {
	Blocks []Block `json:"blocks"` // Array of blocks
}

// NewBlockchain initializes a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	genesisBlock := Block{
		Index:     0,
		Timestamp: time.Now().UTC().String(),
		Data:      "Genesis Block",
		PrevHash:  "",
		Hash:      calculateHash(0, time.Now().UTC().String(), "Genesis Block", ""),
	}
	return &Blockchain{
		Blocks: []Block{genesisBlock},
	}
}

// calculateHash generates a SHA-256 hash for a block's data
func calculateHash(index int, timestamp string, data interface{}, prevHash string) string {
	record := fmt.Sprintf("%d%s%v%s", index, timestamp, data, prevHash)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data interface{}) error {
	// Get the last block in the chain
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	// Create the new block
	newBlock := Block{
		Index:     len(bc.Blocks),
		Timestamp: time.Now().UTC().String(),
		Data:      data,
		PrevHash:  lastBlock.Hash,
		Hash:      calculateHash(len(bc.Blocks), time.Now().UTC().String(), data, lastBlock.Hash),
	}
	// Append the block to the chain
	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

// ValidateChain verifies the integrity of the blockchain
func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		previousBlock := bc.Blocks[i-1]

		// Check if the hash is consistent
		if currentBlock.PrevHash != previousBlock.Hash {
			return false
		}
		// Recompute the hash to ensure integrity
		if currentBlock.Hash != calculateHash(currentBlock.Index, currentBlock.Timestamp, currentBlock.Data, currentBlock.PrevHash) {
			return false
		}
	}
	return true
}

// PrintChain prints the blockchain
func (bc *Blockchain) PrintChain() {
	for _, block := range bc.Blocks {
		fmt.Printf("Index: %d\nTimestamp: %s\nData: %v\nPrevHash: %s\nHash: %s\n\n",
			block.Index, block.Timestamp, block.Data, block.PrevHash, block.Hash)
	}
}

func main() {
	// Initialize blockchain
	chain := NewBlockchain()

	// Add some blocks
	chain.AddBlock(map[string]string{
		"event_type": "data_ingestion",
		"source":     "DataGate",
		"timestamp":  time.Now().UTC().String(),
	})

	chain.AddBlock(map[string]string{
		"event_type": "user_action",
		"actor_id":   "hashed_user_id",
		"timestamp":  time.Now().UTC().String(),
		"action":     "viewed_record",
	})

	// Print the blockchain
	chain.PrintChain()

	// Validate the blockchain
	if chain.ValidateChain() {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is invalid.")
	}
}
```

---

### **How It Works**

1. **Genesis Block**:
   - The blockchain starts with a predefined "Genesis Block" that serves as the foundation.

2. **Adding Blocks**:
   - New blocks are appended by referencing the hash of the previous block.
   - Data is stored as a generic `interface{}` for flexibility.

3. **Hash Calculation**:
   - Each block's hash depends on its `Index`, `Timestamp`, `Data`, and the `PrevHash`.

4. **Validation**:
   - Ensures continuity (`PrevHash` matches the previous block's `Hash`).
   - Verifies integrity by recalculating each block's hash.

---

### **Next Steps**
- **Persist the Blockchain**:
  - Use a storage solution (e.g., BoltDB) to save the blockchain to disk.
- **Expose Blockchain as REST API**:
  - Add endpoints for logging and querying blocks.

---

Integrating encryption directly into the **VitalEdge Blockchain** adds a significant layer of security, particularly because blockchain data is often stored and transmitted in multiple locations. Here's why having blockchain-specific encryption is advantageous:

---

### **Why Separate Encryption for the Blockchain?**

1. **Data Sovereignty**:
   - The blockchain contains sensitive patient and system data, which may require encryption tailored to blockchain's unique structure.
   - Relying on external microservices could create potential bottlenecks or points of failure.

2. **End-to-End Security**:
   - Data stored on the blockchain is immutable, but sensitive fields should still be encrypted to protect against unauthorized access.
   - Built-in encryption ensures that sensitive data is encrypted *before* it is written to the chain.

3. **Granular Encryption**:
   - You can encrypt individual fields within the blockchain (`Data`, for example) instead of encrypting the entire block.
   - This allows selective access control and efficient decryption for authorized users.

4. **Compliance**:
   - Integrating encryption directly into the blockchain aligns with healthcare standards like **HIPAA** and **GDPR** for handling sensitive data.

5. **Streamlined Design**:
   - By embedding encryption/decryption within the blockchain layer, you reduce interdependencies between services and simplify maintenance.

---

### **Recommended Encryption Approach**

1. **Algorithms**:
   - Use AES (Advanced Encryption Standard) for data encryption.
   - Combine with public/private key encryption (e.g., RSA) for role-based decryption.

2. **Hybrid Model**:
   - Encrypt block `Data` using AES.
   - Store the AES key itself encrypted with RSA for secure role-based decryption.

3. **Implementation Details**:
   - Sensitive data (e.g., `Data` field) is encrypted during the block creation process.
   - Decryption occurs when an authorized user queries the blockchain.

---

### **Implementation Plan in Go**

#### **Enhancements to the Blockchain Core**
1. Add **AES-256** encryption for sensitive `Data` fields in each block.
2. Implement **RSA** to protect AES keys for authorized decryption.

#### **Updated Block Structure**
Add an `EncryptedData` field and store encrypted AES keys alongside blocks.

```go
type Block struct {
	Index         int         `json:"index"`
	Timestamp     string      `json:"timestamp"`
	EncryptedData string      `json:"encrypted_data"` // Encrypted block data
	PrevHash      string      `json:"prev_hash"`
	Hash          string      `json:"hash"`
	AESKey        string      `json:"aes_key"`        // Encrypted AES key
}
```

#### **Go Implementation**
Below is the updated blockchain core with encryption:

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// Block represents a single block in the blockchain
type Block struct {
	Index         int    `json:"index"`
	Timestamp     string `json:"timestamp"`
	EncryptedData string `json:"encrypted_data"`
	PrevHash      string `json:"prev_hash"`
	Hash          string `json:"hash"`
	AESKey        string `json:"aes_key"`
}

// Blockchain represents the full chain
type Blockchain struct {
	Blocks []Block
}

// NewBlockchain initializes the blockchain
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []Block{},
	}
}

// AES Encryption
func encryptData(data string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES Key Generator
func generateAESKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// RSA Encryption for AES Key
func encryptAESKey(key []byte, publicKey *rsa.PublicKey) (string, error) {
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, key, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}

// Calculate Hash
func calculateHash(index int, timestamp, encryptedData, prevHash string) string {
	record := fmt.Sprintf("%d%s%s%s", index, timestamp, encryptedData, prevHash)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// AddBlock adds a new block with encrypted data
func (bc *Blockchain) AddBlock(data string, publicKey *rsa.PublicKey) error {
	aesKey, err := generateAESKey()
	if err != nil {
		return err
	}
	encryptedData, err := encryptData(data, aesKey)
	if err != nil {
		return err
	}
	encryptedAESKey, err := encryptAESKey(aesKey, publicKey)
	if err != nil {
		return err
	}

	var prevHash string
	if len(bc.Blocks) > 0 {
		prevHash = bc.Blocks[len(bc.Blocks)-1].Hash
	}

	newBlock := Block{
		Index:         len(bc.Blocks),
		Timestamp:     time.Now().UTC().String(),
		EncryptedData: encryptedData,
		PrevHash:      prevHash,
		Hash:          calculateHash(len(bc.Blocks), time.Now().UTC().String(), encryptedData, prevHash),
		AESKey:        encryptedAESKey,
	}
	bc.Blocks = append(bc.Blocks, newBlock)
	return nil
}

// ValidateChain validates the integrity of the blockchain
func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		current := bc.Blocks[i]
		previous := bc.Blocks[i-1]

		if current.PrevHash != previous.Hash {
			return false
		}
		if current.Hash != calculateHash(current.Index, current.Timestamp, current.EncryptedData, current.PrevHash) {
			return false
		}
	}
	return true
}

func main() {
	// Generate RSA keys for encryption
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey

	// Initialize blockchain
	chain := NewBlockchain()

	// Add blocks
	chain.AddBlock("Sensitive data for block 1", publicKey)
	chain.AddBlock("Sensitive data for block 2", publicKey)

	// Validate and print chain
	if chain.ValidateChain() {
		fmt.Println("Blockchain is valid.")
	} else {
		fmt.Println("Blockchain is invalid.")
	}

	for _, block := range chain.Blocks {
		fmt.Printf("Block %d\nEncrypted Data: %s\nAES Key: %s\n\n", block.Index, block.EncryptedData, block.AESKey)
	}
}
```

---

### **Next Steps**

1. **Decryption Workflow**:
   - Add decryption logic for querying and accessing data securely.

2. **Storage**:
   - Persist encrypted blockchain data in a local database.

---
