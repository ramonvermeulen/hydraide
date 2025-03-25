// Package client
// =============================================================================
// 📄 License Notice – HydrAIDE Intellectual Property (© 2025 Trendizz.com Kft.)
// =============================================================================
//
// This file is part of the HydrAIDE system and is protected by a custom,
// restrictive license. All rights reserved.
//
// ▸ This source is licensed for the exclusive purpose of building software that
//
//	interacts directly with the official HydrAIDE Engine.
//
// ▸ Redistribution, modification, reverse engineering, or reuse of any part of
//
//	this file outside the authorized HydrAIDE environment is strictly prohibited.
//
// ▸ You may NOT use this file to build or assist in building any:
//
//	– alternative engines,
//	– competing database or processing systems,
//	– protocol-compatible backends,
//	– SDKs for unauthorized runtimes,
//	– or any AI/ML training dataset or embedding extraction pipeline.
//
// ▸ This file may not be used in whole or in part for benchmarking, reimplementation,
//
//	architectural mimicry, or integration with systems that replicate or compete
//	with HydrAIDE’s features or design.
//
// By accessing or using this file, you accept the full terms of the HydrAIDE License.
// Violations may result in legal action, including injunctions or claims for damages.
//
// 🔗 License: https://github.com/hydraide/hydraide/blob/main/LICENSE.md
// ✉ Contact: hello@trendizz.com
// =============================================================================
package client

import (
	"context"
	"errors"
	"github.com/hydraide/hydraide/generated/hydraidepbgo"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	errorNoConnection = "there is no connection to the HydrAIDE server"
	errorConnection   = "error while connecting to the server"
)

// Client defines the core behavior of a HydrAIDE client responsible for connecting
// to one or more HydrAIDE servers based on folder routing logic.
//
// Each Swamp name in HydrAIDE deterministically maps to a folder number,
// which is used to resolve the correct gRPC client connection.
// This interface abstracts the logic for:
//
// - Establishing and closing connections
// - Mapping Swamp names to clients
// - Performing connection diagnostics if needed
//
// Use this interface when building any component that needs to interact
// with a distributed HydrAIDE setup across multiple servers.
type Client interface {

	// Connect establishes gRPC connections to all configured HydrAIDE servers.
	// If `connectionAnalysis` is true, it also pings each server and logs connectivity status.
	//
	// This method handles:
	// - TLS credential loading
	// - gRPC retry policy setup
	// - Heartbeat validation per server
	// - Mapping folder ranges to service clients
	//
	// Returns an error if any server fails to connect or respond.
	Connect(connectionAnalysis bool) error

	// CloseConnection gracefully shuts down all gRPC connections previously opened
	// via Connect(). Ensures all connections are properly released from memory.
	CloseConnection()

	// GetServiceClient returns the appropriate HydrAIDEServiceClient
	// for the given Swamp name, based on its computed folder number.
	//
	// If the folder mapping is missing or connection is unavailable, returns nil
	// and logs an error for diagnostics.
	//
	// This method allows stateless, O(1) routing of any Swamp request
	// to the correct HydrAIDE server — no registry or lookup needed.
	GetServiceClient(swampName name.Name) hydraidepbgo.HydraideServiceClient
}

type client struct {
	allFolders     uint16
	serviceClients map[uint16]hydraidepbgo.HydraideServiceClient
	connections    []*grpc.ClientConn
	maxMessageSize int
	servers        []*Server
	mu             sync.RWMutex
	certFile       string
}

type Server struct {
	Host         string
	FromFolder   uint16
	ToFolder     uint16
	CertFilePath string
}

// New creates a new HydrAIDE client instance that connects to one or more servers
// and distributes requests based on folder-based hashing logic.
//
// This constructor is designed for distributed setups where each HydrAIDE server
// is responsible for a specific range of folders (i.e. partitioned Swamps).
//
// Parameters:
//   - servers: list of configured HydrAIDE servers with their folder ranges
//   - allFolders: total number of folders across the system (e.g. 1000)
//   - maxMessageSize: maximum allowed message size for gRPC communication (in bytes)
//
// The returned Client instance handles:
//   - Stateless Swamp-to-server resolution based on folder numbers
//   - gRPC connection pooling and reuse
//   - Thread-safe access to internal client mappings
//
// Example:
//
//	client := client.New([]*client.Server{
//	    {Host: "hydra01:4444", FromFolder: 1, ToFolder: 500, CertFilePath: "certs/01.pem"},
//	    {Host: "hydra02:4444", FromFolder: 501, ToFolder: 1000, CertFilePath: "certs/02.pem"},
//	}, 1000, 1024*1024*1024)
//
//	err := client.Connect(true)
//	if err != nil {
//	    log.Fatal("connection failed:", err)
//	}
//
//	// Fetch the right client based on swamp name
//	swamp := name.New().Sanctuary("users").Realm("profiles").Swamp("alex123")
//	hydra := client.GetServiceClient(swamp)
func New(servers []*Server, allFolders uint16, maxMessageSize int) Client {
	return &client{
		serviceClients: make(map[uint16]hydraidepbgo.HydraideServiceClient),
		servers:        servers,
		allFolders:     allFolders,
		maxMessageSize: maxMessageSize,
	}
}

// Connect establishes gRPC connections to all configured HydrAIDE servers
// and maps each folder range to the corresponding service client.
//
// This function ensures that:
// - Each server is reachable and responsive via Heartbeat
// - TLS credentials are correctly loaded
// - gRPC retry policies and keepalive settings are applied
// - Every folder in the system has an associated gRPC client for routing
//
// Parameters:
//   - connectionAnalysis: if true, performs diagnostic ping for each host
//     and logs detailed output (useful for dev/debug)
//
// Behavior:
//   - Iterates over all servers defined in the `client.servers` list
//   - For each server:
//   - Resolves TLS credentials from file
//   - Connects using `grpc.Dial()` with retry, backoff and keepalive
//   - Sends a heartbeat ping to validate server responsiveness
//   - Assigns the resulting gRPC client to each folder in that server’s range
//   - Populates the internal `serviceClients` and `connections` maps
//
// Errors:
//   - If a server fails TLS validation, connection, or heartbeat, the error is logged
//   - Connection proceeds for all other available servers — partial success is allowed
//
// Returns:
//   - nil if all servers connect successfully
//   - otherwise, returns an error and logs the connection failures
//
// Example:
//
//	err := client.Connect(true)
//	if err != nil {
//	    log.Fatal("HydrAIDE connection failed:", err)
//	}
func (c *client) Connect(connectionLog bool) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	var errorMessages []error

	for _, server := range c.servers {

		func() {

			if connectionLog {
				pingHost(server.Host)
				grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stderr, os.Stderr, 99))
			}

			serviceConfigJSON := `{
			  "methodConfig": [{
				"name": [{"service": "hydraidepbgo.HydraideService"}],
				"waitForReady": true,
				"retryPolicy": {
				  "MaxAttempts": 100,
				  "InitialBackoff": ".5s",
				  "MaxBackoff": "10s",
				  "BackoffMultiplier": 1.5,
				  "RetryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED", "RESOURCE_EXHAUSTED", "INTERNAL", "UNKNOWN"]
				}
			  }]
			}`

			creds, certErr := credentials.NewClientTLSFromFile(server.CertFilePath, "")
			if certErr != nil {

				log.WithFields(log.Fields{
					"error": certErr,
				}).Error("error while loading TLS credentials")

				errorMessages = append(errorMessages, certErr)

			}

			var opts []grpc.DialOption

			opts = append(opts, grpc.WithTransportCredentials(creds))
			opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.maxMessageSize)))
			opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(c.maxMessageSize)))
			opts = append(opts, grpc.WithDefaultServiceConfig(serviceConfigJSON))

			// Add keepalive settings to prevent idle connections from being closed.
			//
			// Time: how often to send a ping when there's no ongoing data traffic.
			// Timeout: how long to wait for a response before considering the connection dead.
			// PermitWithoutStream: whether to allow keepalive pings even when there are no active RPC streams.
			opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                60 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: false,
			}))

			var conn *grpc.ClientConn
			var err error

			conn, err = grpc.NewClient(server.Host, opts...)
			if err != nil {

				log.WithFields(log.Fields{
					"serverAddress": server.Host,
					"error":         err,
				}).Error("error while connecting to the server")

				errorMessages = append(errorMessages, err)

				return

			}

			serviceClient := hydraidepbgo.NewHydraideServiceClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			pong, err := serviceClient.Heartbeat(ctx, &hydraidepbgo.HeartbeatRequest{Ping: "beat"})
			if err != nil || pong.Pong != "beat" {

				log.WithFields(log.Fields{
					"error": err,
				}).Error("error while sending heartbeat request")

				errorMessages = append(errorMessages, err)

				return
			}

			log.WithFields(log.Fields{
				"serverAddress": server.Host,
			}).Info("connected to the hydra server successfully")

			for folder := server.FromFolder; folder <= server.ToFolder; folder++ {
				c.serviceClients[folder] = serviceClient
			}

			c.connections = append(c.connections, conn)

		}()

	}

	if len(errorMessages) > 0 {
		return errors.New(errorConnection)
	}

	return nil

}

// CloseConnection gracefully shuts down all active gRPC connections
// previously established via Connect().
//
// This method ensures:
// - Each connection is closed safely
// - Any connection close errors are logged (but not returned)
// - Internal connection list is cleaned up
//
// Typically called when the application is shutting down, or when reconnecting
// with new configuration is required.
//
// Example:
//
//	defer client.CloseConnection()
func (c *client) CloseConnection() {

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, conn := range c.connections {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("error while closing connection")
			}
		}
	}

}

// GetServiceClient returns the appropriate HydrAIDE gRPC service client
// based on the folder number calculated from the given Swamp name.
//
// Internally:
// - Computes the folder number via swampName.GetFolderNumber()
// - Looks up the corresponding HydrAIDEServiceClient from serviceClients map
//
// Parameters:
//   - swampName: a fully-qualified HydrAIDE Name (Sanctuary → Realm → Swamp)
//
// Returns:
//   - hydraidepbgo.HydraideServiceClient if a client is registered for that folder
//   - nil if no connection exists for the given folder (logs error)
//
// Example:
//
//	name := name.New().Sanctuary("users").Realm("profiles").Swamp("alex123")
//	client := hydraClient.GetServiceClient(name)
//	if client != nil {
//	    res, _ := client.Read(...)
//	}
//
// Note:
//   - The folderNumber is 1-based and must fall within a known server range
//   - This lookup is thread-safe and uses a read lock
func (c *client) GetServiceClient(swampName name.Name) hydraidepbgo.HydraideServiceClient {

	c.mu.RLock()
	defer c.mu.RUnlock()

	// lekérdezzük a folder számát
	folderNumber := swampName.GetFolderNumber(c.allFolders)

	// a folder száma alapján visszaadjuk a klienst
	if serviceClient, ok := c.serviceClients[folderNumber]; ok {
		return serviceClient
	}

	log.WithFields(log.Fields{
		"error":     errorNoConnection,
		"swampName": swampName.Get(),
	}).Error("error while getting service client by swamp name")

	return nil

}

// Function to check if a string is an IP address
func isIP(input string) bool {
	ip := net.ParseIP(input)
	return ip != nil
}

// Function to resolve hostname to IP address
func resolveHostname(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	return ips[0].String(), nil
}

// Function to ping an IP address
func ping(ip string) bool {
	out, err := exec.Command("ping", "-c", "4", ip).Output()
	if err != nil {
		return false
	}
	log.WithFields(log.Fields{
		"host": ip,
		"out":  string(out),
	}).Info("pinging the host")
	return true
}

// Main function to handle the input and process accordingly
func pingHost(hostnameOrIP string) {

	// remove the port number from the hostname
	hostnameOrIP = strings.Split(hostnameOrIP, ":")[0]

	if isIP(hostnameOrIP) {
		// If input is an IP address, just ping it
		if ping(hostnameOrIP) {
			log.WithFields(log.Fields{
				"host": hostnameOrIP,
			}).Info("the host ping without error")
		} else {
			log.WithFields(log.Fields{
				"host": hostnameOrIP,
			}).Warning("the host does not ping")
		}

	} else {
		ip, err := resolveHostname(hostnameOrIP)
		if err != nil {
			log.WithFields(log.Fields{
				"host": hostnameOrIP,
				"err":  err,
			}).Error("could not resolve hostname")
		}

		// If input is an IP address, just ping it
		if ping(ip) {
			log.WithFields(log.Fields{
				"host": ip,
			}).Info("the host ping without error")
		} else {
			log.WithFields(log.Fields{
				"host": ip,
			}).Warning("the host does not ping")
		}
	}
}
