package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// EventRequest defines the payload for LogEvent
type EventRequest struct {
	EventID      string `json:"event_id"`
	EventType    string `json:"event_type"`
	EventDetails string `json:"event_details"`
	User         string `json:"user"`
	Timestamp    string `json:"timestamp"`
}

// TokenRequest defines the payload for IncrementToken
type TokenRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	Amount int    `json:"amount"`
}

func main() {
	r := mux.NewRouter()

	// Define REST endpoints
	r.HandleFunc("/api/log-event", LogEventHandler).Methods("POST")
	r.HandleFunc("/api/increment-token", IncrementTokenHandler).Methods("POST")

	log.Println("Starting REST API server on :8082")
	log.Fatal(http.ListenAndServe(":8082", r))
}

// LogEventHandler handles POST requests to log an event
func LogEventHandler(w http.ResponseWriter, r *http.Request) {
	peer := r.URL.Query().Get("peer") // Get the peer from the query parameter

	var peerEndpoint string
	log.Println("peer received is:", peer)
	switch peer {
	case "clinicians.xmed.ai":
		peerEndpoint = "peer0.clinicians.xmed.ai:7051"
	case "scientists.xnome.net":
		peerEndpoint = "peer0.scientists.xnome.net:9051"
	default:
		// http.Error(w, "Invalid or missing peer parameter", http.StatusBadRequest)
		// return
		peerEndpoint = "peer0.clinicians.xmed.ai:7051"
	}

	// Use the `peerEndpoint` to connect to the correct peer.
	fmt.Println("Selected Peer Endpoint:", peerEndpoint)

	var req EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call LogEvent on the blockchain
	err := invokeBlockchain("LogEvent", req.EventID, req.EventType, req.EventDetails, req.User, req.Timestamp)
	if err != nil {
		http.Error(w, "Failed to log event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event logged successfully"))
}

// IncrementTokenHandler handles POST requests to increment a token
func IncrementTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call IncrementToken on the blockchain
	err := invokeBlockchain("IncrementToken", req.UserID, req.Token, req.Amount)
	if err != nil {
		http.Error(w, "Failed to increment token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token incremented successfully"))
}

// Function to dynamically get the private key file from the directory
func getPrivateKeyPath(keystoreDir string) (string, error) {
	// List all files in the keystore directory
	files, err := ioutil.ReadDir(keystoreDir)
	if err != nil {
		return "", fmt.Errorf("failed to read keystore directory: %w", err)
	}

	// Ensure there's at least one file
	if len(files) == 0 {
		return "", fmt.Errorf("no files found in keystore directory: %s", keystoreDir)
	}

	// Return the full path of the first file (assuming only one key file exists)
	return filepath.Join(keystoreDir, files[0].Name()), nil
}

func invokeBlockchain(function string, args ...interface{}) error {

	// Load the connection profile
	// connectionProfile := "../path/to/connection.json" // Adjust this path
	peerEndpoint := "peer0.clinicians.xmed.ai:7051" // Update with your actual peer endpoint
	tlsCertPath := "/app/org1/tls/ca.crt"

	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return fmt.Errorf("failed to read TLS cert: %w", err)
	}
	log.Println("Successfully read TLS cert")

	// Create TLS credentials
	creds := credentials.NewClientTLSFromCert(parseCert(cert), "")
	log.Println("Successfully created credentials")

	// Create gRPC connection
	grpcConn, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %w", err)
	}
	defer grpcConn.Close()
	log.Println("Successfully created grps connection")

	// Load certificates and private key for the identity
	certPath := "/app/org1/msp/signcerts/cert.pem"
	// privateKeyPath := "/app/certs/private-key.pem"
	// certPath := "../path/to/cert.pem"              // Update to actual cert path
	// privateKeyPath := "../path/to/private_key.pem" // Update to actual private key path

	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read identity cert: %w", err)
	}
	log.Println("Successfully read identity cert from:", certPath)

	// Directory containing the private key
	keystoreDir := "/app/org1/msp/keystore"

	// Dynamically determine the private key file path
	privateKeyPath, err := getPrivateKeyPath(keystoreDir)
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	// Other existing logic
	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}
	// privateKeyPEM, err := os.ReadFile(privateKeyPath)
	// if err != nil {
	// 	return fmt.Errorf("failed to read private key: %w", err)
	// }
	log.Println("Successfully read private key from:", privateKeyPath)

	// Parse the X.509 certificate
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	log.Println("Successfully parsed certificate PEM")

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse X.509 certificate: %w", err)
	}
	log.Println("Successfully parsed X.509 certificate")

	// Parse the private key
	block, _ = pem.Decode(privateKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to parse private key PEM")
	}
	log.Println("Successfully parsed private key PEM")

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	log.Println("Successfully parsed private key")

	// Create identity and signer
	id, err := identity.NewX509Identity("Org1MSP", certificate)
	if err != nil {
		return fmt.Errorf("failed to create identity: %w", err)
	}
	log.Println("Successfully created identity - Org1MSP")

	// Assert that the parsed private key is of type *ecdsa.PrivateKey
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return fmt.Errorf("unsupported key type: %T", privateKey)
	}

	// Pass the ecdsa.PrivateKey directly to NewPrivateKeySign
	signer, err := identity.NewPrivateKeySign(ecdsaPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to create signer: %w", err)
	}

	// signer, err := identity.NewPrivateKeySign(privateKey.(interface {
	// 	Sign(rand io.Reader, digest []byte, opts interface{}) ([]byte, error)
	// }))
	// if err != nil {
	// 	return fmt.Errorf("failed to create signer: %w", err)
	// }

	log.Println("Successfully created signer")

	// Connect to the fabric gateway
	gw, err := client.Connect(
		id,
		client.WithClientConnection(grpcConn),
		client.WithSign(signer),
		client.WithEvaluateTimeout(10*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(20*time.Second),
		client.WithCommitStatusTimeout(30*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}
	defer gw.Close()
	log.Println("Successfully connected to gateway")

	// Get the network and contract
	network := gw.GetNetwork("vitaledgechannel")
	contract := network.GetContract("vitaledgechaincode")

	// Prepare transaction arguments
	argsJSON, _ := json.Marshal(args)
	log.Println("Prepared transaction arguments:", string(argsJSON))

	// Submit the transaction
	log.Println("Submitting transaction")
	_, err = contract.SubmitTransaction(function, string(argsJSON))
	return err
}

func parseCert(cert []byte) *x509.CertPool {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	return certPool
}

// // invokeBlockchain abstracts the blockchain invocation logic using fabric-gateway
// func invokeBlockchain(function string, args ...interface{}) error {
// 	// Load the connection profile
// 	connectionProfile := "../connection.json" // Adjust this path as needed
// 	connectionBytes, err := os.ReadFile(connectionProfile)
// 	if err != nil {
// 		return fmt.Errorf("failed to read connection profile: %w", err)
// 	}

// 	// Load certificates and private key for the identity
// 	certPath := "../path/to/cert.pem"              // Update this to the actual path to the cert
// 	privateKeyPath := "../path/to/private_key.pem" // Update this to the actual path to the private key

// 	cert, err := os.ReadFile(certPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read cert: %w", err)
// 	}
// 	privateKey, err := os.ReadFile(privateKeyPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read private key: %w", err)
// 	}

// 	// Parse the private key
// 	block, _ := pem.Decode(privateKey)
// 	if block == nil {
// 		return fmt.Errorf("failed to parse private key PEM")
// 	}
// 	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse private key: %w", err)
// 	}

// 	// Create the identity and signer
// 	id := client.NewX509Identity("Org1MSP", string(cert))
// 	signer, err := client.NewSigner(key)
// 	if err != nil {
// 		return fmt.Errorf("failed to create signer: %w", err)
// 	}

// 	// Set up the gateway configuration
// 	gatewayConfig := &client.Gateway{
// 		Identity: id,
// 		Signer:   signer,
// 	}

// 	// Establish the connection
// 	gw, err := client.Connect(gatewayConfig, connectionBytes)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to gateway: %w", err)
// 	}
// 	defer gw.Close()

// 	// Get the network and contract
// 	network := gw.GetNetwork("vitaledgechannel")
// 	contract := network.GetContract("vitaledgechaincode")

// 	// Prepare transaction arguments
// 	argsJSON, _ := json.Marshal(args)

// 	// Submit the transaction
// 	_, err = contract.SubmitTransaction(function, string(argsJSON))
// 	return err
// }

// // invokeBlockchain abstracts the blockchain invocation logic using fabric-gateway
// func invokeBlockchain(function string, args ...interface{}) error {
// 	// Load the connection profile
// 	connectionProfile := "../connection.json" // Adjust this path as needed
// 	connectionBytes, err := os.ReadFile(connectionProfile)
// 	if err != nil {
// 		return fmt.Errorf("failed to read connection profile: %w", err)
// 	}

// 	// Load certificates and private key for the identity
// 	certPath := "../path/to/cert.pem"              // Update this to the actual path to the cert
// 	privateKeyPath := "../path/to/private_key.pem" // Update this to the actual path to the private key

// 	cert, err := os.ReadFile(certPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read cert: %w", err)
// 	}
// 	privateKey, err := os.ReadFile(privateKeyPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read private key: %w", err)
// 	}

// 	// Parse the private key
// 	block, _ := pem.Decode(privateKey)
// 	if block == nil {
// 		return fmt.Errorf("failed to parse private key PEM")
// 	}
// 	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse private key: %w", err)
// 	}

// 	// Create the identity and signer
// 	// id := client.NewX509Identity("Org1MSP", string(cert))
// 	// signer, err := client.NewSigner(key)
// 	id := client.NewIdentity("Org1MSP", cert)
// 	signer, err := client.NewSigner(key)
// 	if err != nil {
// 		return fmt.Errorf("failed to create signer: %w", err)
// 	}

// 	// Connect to the fabric gateway
// 	gw, err := client.Connect(
// 		client.WithIdentity(id),
// 		client.WithSigner(signer),
// 		client.WithEndpoint("localhost:7051"), // Update with the actual peer endpoint
// 		client.WithTLSConfig(client.WithTLSCertsFromFile("../path/to/tls-cert.pem")),
// 	)
// 	client.WithBytesArguments()
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to gateway: %w", err)
// 	}
// 	defer gw.Close()

// 	// Get the network and contract
// 	network := gw.GetNetwork("vitaledgechannel")
// 	contract := network.GetContract("vitaledgechaincode")

// 	// Prepare transaction arguments
// 	argsJSON, _ := json.Marshal(args)

// 	// Submit the transaction
// 	_, err = contract.SubmitTransaction(function, string(argsJSON))
// 	return err
// }
