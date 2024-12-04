### **VitalEdge Blockchain: Detailed Development Plan**

This plan outlines manageable steps to evolve the **VitalEdge Blockchain** into a full-fledged system supporting personalized medicine and decentralized healthcare. We will leverage Hyperledger Fabric's framework, starting with a basic blockchain and gradually adding functionality to achieve the vision outlined.

This phased approach balances immediate functionality with long-term scalability, ensuring the system evolves in alignment with the decentralized healthcare vision.

---

### **Phase 1: Foundational Functionality**
1. **Logging Transactions into the Blockchain**
    - **Objective:** Create a mechanism to log all relevant VitalEdge system activities immutably.
    - **Actions:**
      - Modify `smartcontract.go` to include a new function: `LogTransaction()`.
      - Each transaction will record:
        - **Actor ID:** The participant responsible for the activity (clinician, scientist, etc.).
        - **Activity Type:** E.g., data contribution, ML pipeline configuration, clinical interaction.
        - **Timestamp.**
        - **Metadata:** Additional optional information like context, location, or importance.
      - Update the chaincode to write transaction logs into the blockchain as key-value pairs.
    - **Testing:** Invoke `LogTransaction()` using the CLI and verify transactions using `peer chaincode query`.

2. **Flagging Activities for Immutability and Proof of Participation**
    - **Objective:** Identify activities that need permanent recording or serve as a basis for token rewards.
    - **Actions:**
      - Add a `FlaggedTransaction` structure:
        ```go
        type FlaggedTransaction struct {
          ActorID      string    `json:"actor_id"`
          ActivityType string    `json:"activity_type"`
          Importance   int       `json:"importance"` // Scale of 1-10
          Immutable    bool      `json:"immutable"`
          Timestamp    string    `json:"timestamp"`
        }
        ```
      - Create a `LogFlaggedActivity()` function in the chaincode:
        - Records flagged transactions.
        - Adds a flag for immutability or marks it for token credit.
        - Uses a separate keyspace (`FlaggedTransactions`) in the ledger.
    - **Testing:** Log flagged activities and verify their immutability.

3. **Tokenization Framework**
    - **Objective:** Introduce a basic framework for creating and managing tokens.
    - **Actions:**
      - Add support for four tokens: GitCoin, HitCoin, KitCoin, FitCoin.
      - Implement `CreateToken()` and `AllocateToken()` functions:
        - `CreateToken()` initializes a token type (e.g., GitCoin).
        - `AllocateToken()` credits tokens to participants based on flagged activities.
        - Tokens will be stored as:
          ```go
          type Token struct {
            TokenID   string `json:"token_id"`
            ActorID   string `json:"actor_id"`
            Amount    int    `json:"amount"`
            TokenType string `json:"token_type"`
          }
          ```
      - Explore a potential fifth `EdgeToken` that aggregates contributions.
    - **Testing:** Allocate tokens and query token balances for participants.

4. **View Blockchain State**
    - **Objective:** Enable users to query blockchain data programmatically and visually.
    - **Actions:**
      - Add `GetTransaction()` and `GetAllTransactions()` functions to chaincode.
      - Integrate CLI commands for querying blockchain state.
      - Develop a lightweight web-based dashboard for viewing logged transactions and token balances.
    - **Testing:** Query and verify ledger states through CLI and UI.

---

### **Phase 2: RESTful API Development**
1. **REST API Wrapper**
    - **Objective:** Build a REST API for external systems to interact with the blockchain.
    - **Actions:**
      - Use `fiber` or `gin` framework for the Go-based REST API.
      - Create endpoints:
        - `/log-transaction`: Logs a transaction into the blockchain.
        - `/flag-activity`: Flags an activity for immutability or token allocation.
        - `/allocate-token`: Allocates tokens to a participant.
        - `/query-transactions`: Fetches transaction logs.
        - `/query-tokens`: Fetches token balances.
      - Secure endpoints with basic authentication and API keys.
    - **Testing:** Test API endpoints using tools like Postman.

2. **Integration with VitalEdge Crypt API**
    - **Objective:** Use the existing encryption/decryption service to secure sensitive data before blockchain logging.
    - **Actions:**
      - Add middleware to encrypt sensitive data fields before logging transactions.
      - Decrypt data upon querying.

---

### **Phase 3: Advanced Features**
1. **Enhancing Tokenization**
    - Introduce token-burning and transfer functionalities.
    - Implement rules for EdgeToken calculation as a weighted sum of other tokens.

2. **Multi-Channel Support**
    - Create additional channels for specific purposes (e.g., ML pipelines, genomic data sharing).

3. **Audit and Reporting**
    - Implement an audit trail for flagged immutable activities.
    - Add reporting capabilities for token balances and activity summaries.

4. **Security Hardening**
    - Shift from **solo** ordering to **etcdraft** for high availability.
    - Explore external identity providers for authentication.

---

### **Phase 4: Ecosystem Expansion**
1. **Stakeholder-Specific Features**
    - Clinicians:
      - Create a dashboard for querying patient-specific logs.
    - Scientists:
      - Add support for contributing ML models and data pipelines.
    - Patients:
      - Introduce an app for tracking participation and health progress.

2. **Integration with Other Networks**
    - Use inter-ledger protocols to integrate with other blockchain ecosystems.

---

### **Immediate Next Steps**
1. Finalize directory structure and commit current state to Git.
2. Modify `smartcontract.go` to add basic transaction logging functionality.
3. Add CLI commands and test initial blockchain interactions.
4. Set up the REST API skeleton.
