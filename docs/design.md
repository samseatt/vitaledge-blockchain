# **Design Document for VitalEdge Blockchain**

This design provides a robust, secure, and efficient framework for implementing VitalEdge Blockchain, ensuring the system achieves traceability, immutability, and ease of integration.

## **1. Overview**

**VitalEdge Blockchain** is a private blockchain-based microservice developed in Go to provide immutable and secure logging for the VitalEdge ecosystem. Its primary purpose is to record critical events, ensuring traceability, tamper-proof auditability, and compliance with regulatory requirements. The service wraps the blockchain with a REST API for ease of integration and may optionally provide administrative and visualization UI components.

This design document outlines the architecture, blockchain implementation, microservice interfaces, and the operational framework for VitalEdge Blockchain.

---

## **2. Objectives**

- **Immutable Logging**:
  - Record sensitive and critical system events (e.g., data provenance, security events, or pipeline actions).
  - Ensure logs cannot be tampered with once recorded.
- **Traceability and Auditability**:
  - Provide an audit trail for system interactions.
  - Enable authorized administrators to query and verify logs.
- **Integration-Friendly Interface**:
  - Offer a REST API for microservices to interact with the blockchain.
- **Efficient Blockchain Implementation**:
  - Design a lightweight, private blockchain optimized for high write and read performance.

---

## **3. Core Responsibilities**

1. **Blockchain Implementation**:
   - Create a private blockchain that maintains an immutable ledger of events.
   - Support basic blockchain features like proof-of-ownership, timestamping, and chaining blocks.
2. **REST API**:
   - Expose endpoints for writing events, querying logs, and verifying data integrity.
3. **Admin Tools**:
   - Optional UI for visualizing blockchain data and querying logs.
4. **Security and Scalability**:
   - Ensure the blockchain is resistant to tampering and scalable for high-frequency writes.

---

## **4. Architecture**

### **Core Components**

#### **4.1. Blockchain Layer**
Implements the private blockchain with the following key features:
- **Block Structure**:
  - Holds event data, hash of the previous block, and a timestamp.
- **Consensus Algorithm**:
  - Simplified Proof-of-Authority (PoA) since it is a private blockchain with trusted participants.
- **Chain Validation**:
  - Ensures the integrity of the chain using cryptographic hashing.

#### **4.2. REST API Layer**
Provides external interfaces for microservices and administrators to interact with the blockchain.

#### **4.3. Persistence Layer**
Stores the blockchain data persistently in a local or distributed database, ensuring durability and availability.

#### **4.4. Admin Interface (Optional)**
A web-based dashboard for administrators to query and visualize blockchain events.

---

## **5. Interfaces**

### **5.1. REST API**

#### **1. Write Event**
- **Endpoint**: `/log-event`
- **Method**: POST
- **Description**: Records a new event in the blockchain.
- **Request**:
  ```json
  {
      "event_type": "data_ingestion",
      "event_data": {
          "data_id": "12345",
          "source": "DataGate",
          "timestamp": "2023-11-17T10:45:00Z"
      }
  }
  ```
- **Response**:
  ```json
  {
      "status": "success",
      "block_hash": "abc123",
      "block_index": 42
  }
  ```

#### **2. Query Event**
- **Endpoint**: `/query-event`
- **Method**: POST
- **Description**: Queries blockchain data by event type or other criteria.
- **Request**:
  ```json
  {
      "filters": {
          "event_type": "data_ingestion",
          "date_range": ["2023-11-16", "2023-11-17"]
      }
  }
  ```
- **Response**:
  ```json
  [
      {
          "block_index": 42,
          "event_type": "data_ingestion",
          "event_data": {
              "data_id": "12345",
              "source": "DataGate",
              "timestamp": "2023-11-17T10:45:00Z"
          }
      }
  ]
  ```

#### **3. Verify Data**
- **Endpoint**: `/verify-data`
- **Method**: POST
- **Description**: Verifies the integrity of the blockchain or a specific event.
- **Request**:
  ```json
  {
      "data_id": "12345"
  }
  ```
- **Response**:
  ```json
  {
      "status": "valid",
      "block_index": 42,
      "block_hash": "abc123"
  }
  ```

---

### **5.2. Admin Interface (Optional)**
- **Log Explorer**:
  - Allows admins to browse blockchain logs.
  - Filters by event type, date range, or source.
- **Verification Tool**:
  - Verifies the integrity of specific data or the entire chain.

---

## **6. Blockchain Implementation**

### **6.1. Block Structure**
```go
type Block struct {
    Index        int          `json:"index"`
    Timestamp    string       `json:"timestamp"`
    Data         EventData    `json:"data"`
    PrevHash     string       `json:"prev_hash"`
    Hash         string       `json:"hash"`
}

type EventData struct {
    EventType string            `json:"event_type"`
    EventData map[string]string `json:"event_data"`
}
```

### **6.2. Blockchain Structure**
```go
type Blockchain struct {
    Blocks []Block `json:"blocks"`
}
```

---

### **6.3. Functions**

#### **1. Create Block**
Generates a new block and appends it to the blockchain.
```go
func CreateBlock(data EventData, prevHash string) Block {
    newBlock := Block{
        Index:     len(blockchain.Blocks) + 1,
        Timestamp: time.Now().UTC().String(),
        Data:      data,
        PrevHash:  prevHash,
        Hash:      calculateHash(newBlock),
    }
    blockchain.Blocks = append(blockchain.Blocks, newBlock)
    return newBlock
}
```

#### **2. Calculate Hash**
Generates a hash for a block.
```go
func calculateHash(block Block) string {
    record := strconv.Itoa(block.Index) + block.Timestamp + block.PrevHash + fmt.Sprintf("%v", block.Data)
    hash := sha256.Sum256([]byte(record))
    return fmt.Sprintf("%x", hash)
}
```

#### **3. Validate Chain**
Ensures all blocks are linked correctly and no tampering has occurred.
```go
func ValidateChain(chain Blockchain) bool {
    for i := 1; i < len(chain.Blocks); i++ {
        if chain.Blocks[i].PrevHash != chain.Blocks[i-1].Hash {
            return false
        }
    }
    return true
}
```

---

## **7. Frameworks and Libraries**

1. **Go HTTP Framework**:
   - **Fiber** or **Gin** for REST API development.
     - [Fiber](https://github.com/gofiber/fiber)
     - [Gin](https://github.com/gin-gonic/gin)

2. **Hashing**:
   - **crypto/sha256**: Built-in library for hashing.
     - [crypto/sha256](https://pkg.go.dev/crypto/sha256)

3. **Storage**:
   - **BoltDB** or **LevelDB** for lightweight persistent storage.
     - [BoltDB](https://github.com/etcd-io/bbolt)
     - [LevelDB](https://github.com/google/leveldb)

4. **Logging**:
   - **Zap** or **Logrus** for structured logging.
     - [Zap](https://github.com/uber-go/zap)
     - [Logrus](https://github.com/sirupsen/logrus)

---

## **8. Deployment**

### **Dockerfile**
```dockerfile
FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o blockchain-service
EXPOSE 8085
CMD ["./blockchain-service"]
```

---

### **Orchestration**
- Deploy on AWS ECS or Kubernetes.
- Configure service discovery for dependent services.

---

## **9. Security Considerations**

1. **TLS Encryption**:
   - Use HTTPS for all REST API communications.
2. **Access Control**:
   - Restrict access to admin endpoints.
3. **Tamper Detection**:
   - Validate the blockchain on startup and periodically.

---

## **10. Development Phases**

### **Phase 1: Blockchain Implementation**
- Implement core blockchain logic (e.g., create, hash, validate).

### **Phase 2: REST API Development**
- Build REST endpoints for writing, querying, and verifying data.

### **Phase 3: Persistence Layer**
- Add BoltDB or LevelDB for persistent storage.

### **Phase 4: Admin Dashboard**
- Develop optional UI for querying and visualizing blockchain events.
