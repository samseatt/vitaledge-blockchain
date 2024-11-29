
## **Tutorial: How to Use and Work with the Running Example Blockchain**

This tutorial helps you explore how to **test, verify, and inspect the behavior of vitaledge blockchain**. We'll focus on how to confirm that our added features (e.g., logging events, managing EdgeCoins) work as expected and gain insights into how the blockchain operates (e.g., block creation, ledger state, transaction details).

### **1. Validate Your Chaincode Features**
After deploying our updated chaincode (`vitaledgechaincode`) to the `vitaledgechannel`, you can invoke its functions and query its state using the CLI or scripts.

#### **a. Test Chaincode Functions**
Use the following commands to test the new features:

1. **Log an Event**:
   ```bash
   docker exec -it peer0.org1.example.com peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem -C vitaledgechannel -n basic -c '{"Args":["LogEvent","event1","CreateEdgeCoin","Initial token assignment for user1"]}'
   ```
   This should write a new `EventLog` object to the ledger.

2. **Create an EdgeCoin**:
   ```bash
   docker exec -it peer0.org1.example.com peer chaincode invoke -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem -C vitaledgechannel -n basic -c '{"Args":["CreateEdgeCoin","user1","100"]}'
   ```

3. **Query EdgeCoin Balance**:
   ```bash
   docker exec -it peer0.org1.example.com peer chaincode query -C vitaledgechannel -n basic -c '{"Args":["GetEdgeCoinBalance","user1"]}'
   ```
   The output should display the details of the `EdgeCoin` for `user1`.

4. **Query All Event Logs** (optional custom function):
   If you add a function to retrieve all logged events, you can query it:
   ```bash
   docker exec -it peer0.org1.example.com peer chaincode query -C vitaledgechannel -n basic -c '{"Args":["GetAllEvents"]}'
   ```

---

### **2. View Blockchain Details**

#### **a. Inspect the World State**
The "world state" is the current state of all data on the blockchain, stored as key-value pairs in CouchDB or LevelDB (depending on your setup).

1. **Enable CouchDB for the Peer**:
   If not already done, update your `docker-compose.yml` file to enable CouchDB for `peer0.org1.example.com`.

2. **Access the CouchDB Web UI**:
   - Open a browser and navigate to `http://localhost:5984/_utils/`.
   - Login with `admin` credentials (default: `admin/password`).
   - Explore the database corresponding to your peer (e.g., `vitaledgechannel_ledger`).
   - Look for the `event1` or `user1` entries to confirm your chaincode's impact.

#### **b. Inspect the Ledger**
The ledger contains all the transactions and blocks.

1. **Query the Ledger State**:
   Use the following command to query the current ledger state for `peer0.org1.example.com`:
   ```bash
   docker exec -it peer0.org1.example.com peer channel fetch 0 /tmp/genesis_block.pb -o orderer.example.com:7050 -c vitaledgechannel --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem
   ```
   Decode the block:
   ```bash
   configtxlator proto_decode --input /tmp/genesis_block.pb --type common.Block --output genesis_block.json
   ```
   Review the `genesis_block.json` file to see the ledger's content.

2. **Fetch Recent Blocks**:
   Fetch recent blocks to see how transactions are grouped into blocks:
   ```bash
   docker exec -it peer0.org1.example.com peer channel fetch newest /tmp/newest_block.pb -o orderer.example.com:7050 -c vitaledgechannel --tls --cafile /etc/hyperledger/fabric-ca/ordererOrg/localhost-9054-ca-orderer.pem
   ```
   Decode and inspect the block:
   ```bash
   configtxlator proto_decode --input /tmp/newest_block.pb --type common.Block --output newest_block.json
   jq '.' newest_block.json
   ```

---

### **3. Understand What Happens Internally**
Here’s what happens in response to your chaincode operations:
- **Transactions**:
  - Every chaincode invocation (e.g., `CreateEdgeCoin`, `LogEvent`) creates a transaction.
  - Transactions are recorded in the blockchain ledger.
- **World State Updates**:
  - The key-value store (world state) is updated for any `PutState` calls in your chaincode.
- **Blocks**:
  - Transactions are grouped into blocks and added to the blockchain in sequential order.
  - Blocks are linked cryptographically using hashes, ensuring immutability.

---

### **4. Debugging and Logs**

#### **a. Peer Logs**
If a command fails, inspect the logs for `peer0.org1.example.com` to debug:
```bash
docker logs peer0.org1.example.com
```

#### **b. Chaincode Logs**
If there’s an error in your chaincode, look for logs:
```bash
docker logs dev-peer0.org1.example.com-basic_1.0-xyz123
```
(Replace `xyz123` with the actual container ID or name for your chaincode.)

---

### **5. Advanced: Write a Client Script**
Instead of manually invoking commands, write a Go or Node.js script to interact with the blockchain.

#### **Example: Invoke `LogEvent` via CLI**
Write a script in Go:
```go
package main

import (
    "fmt"
    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func main() {
    client, err := setupChannelClient()
    if err != nil {
        panic(err)
    }

    response, err := client.Execute(channel.Request{
        ChaincodeID: "basic",
        Fcn:         "LogEvent",
        Args:        [][]byte{[]byte("event1"), []byte("CreateEdgeCoin"), []byte("Initial token assignment for user1")},
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Transaction ID: %s\n", response.TransactionID)
}

func setupChannelClient() (*channel.Client, error) {
    // Channel client setup logic here
}
```

Compile and run the script to automate interaction with the chaincode.

---

### **Summary**
By testing the chaincode and inspecting the blockchain:
1. We can confirm our features are working (via CLI or scripts).
2. We can inspect the blockchain's internals (blocks, ledger, world state).
3. We can debug errors effectively using logs.
