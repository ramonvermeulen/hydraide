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

type Client interface {
	Connect(connectionAnalysis bool) error
	CloseConnection()
	GetServiceClient(swampName name.Name) hydraidepbgo.HydraideServiceClient
	GetServiceClientAndHost(swampName name.Name) *ServiceClient
	GetUniqueServiceClients() []hydraidepbgo.HydraideServiceClient
}

type ServiceClient struct {
	GrpcClient hydraidepbgo.HydraideServiceClient
	Host       string
}

type client struct {
	allFolders     uint16
	serviceClients map[uint16]*ServiceClient
	uniqueServices []hydraidepbgo.HydraideServiceClient
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

// New creates a new HydrAIDE client instance that connects to one or more servers,
// and distributes requests based on folder-based hashing logic.
//
// This constructor is designed for distributed setups, where each HydrAIDE server
// is responsible for a specific range of folders (i.e. Swamp partitions).
//
// Parameters:
//   - servers: list of HydrAIDE servers to connect to (each with a folder range and TLS certificate)
//   - allFolders: total number of folders across the entire system (e.g. 1000)
//   - maxMessageSize: maximum allowed message size for gRPC requests/responses (in bytes)
//
// The returned Client instance handles:
//   - Stateless resolution of Swamp names to folders and corresponding servers
//   - Thread-safe access to internal gRPC connections (via RWMutex)
//   - Internal mapping from folder → *ServiceClient (which includes the gRPC client and Host info)
//
// Note:
// - Each gRPC connection is created lazily upon calling `Connect()`
// - The folder-based routing is consistent across all services, derived from Swamp name hash
//
// Example:
//
//	client := client.New([]*client.Server{
//	    {Host: "hydra01:4444", FromFolder: 1, ToFolder: 500, CertFilePath: "certs/01.pem"},
//	    {Host: "hydra02:4444", FromFolder: 501, ToFolder: 1000, CertFilePath: "certs/02.pem"},
//	}, 1000, 1024*1024*1024) // 1 GB max message size
//
//	err := client.Connect(true)
//	if err != nil {
//	    log.Fatal("connection failed:", err)
//	}
//
//	swamp := name.New().Sanctuary("users").Realm("profiles").Swamp("alex123")
//	service := client.GetServiceClient(swamp)
//	if service != nil {
//	    res, err := service.Read(...) // raw gRPC call
//	}
func New(servers []*Server, allFolders uint16, maxMessageSize int) Client {
	return &client{
		serviceClients: make(map[uint16]*ServiceClient),
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
//   - connectionAnalysis: if true, performs diagnostic ping for each Host
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
				c.serviceClients[folder] = &ServiceClient{
					GrpcClient: serviceClient,
					Host:       server.Host,
				}
			}

			c.connections = append(c.connections, conn)
			c.uniqueServices = append(c.uniqueServices, serviceClient)

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

// GetServiceClient returns the raw gRPC HydrAIDE service client for the given Swamp name.
//
// Internally:
// - Computes the folder number using the hash of the swamp name
// - Looks up the matching gRPC client from the internal serviceClients map
//
// Parameters:
//   - swampName: a fully-qualified HydrAIDE Name (Sanctuary → Realm → Swamp)
//
// Returns:
//   - hydraidepbgo.HydraideServiceClient (bound to the correct server)
//   - nil if no client is registered for the calculated folder
//
// Example:
//
//	name := name.New().Sanctuary("users").Realm("profiles").Swamp("alex123")
//	client := hydraClient.GetServiceClient(name)
//	if client != nil {
//	    res, _ := client.Read(...)
//	}
//
// Notes:
//   - The folder number is calculated from the swamp name hash, then routed to the correct server
//   - The lookup is thread-safe (uses a read lock)
//   - This method provides the low-level client only; use GetServiceClientWithMeta() if you need Host info or full routing metadata
func (c *client) GetServiceClient(swampName name.Name) hydraidepbgo.HydraideServiceClient {

	c.mu.RLock()
	defer c.mu.RUnlock()

	// lekérdezzük a folder számát
	folderNumber := swampName.GetFolderNumber(c.allFolders)

	// a folder száma alapján visszaadjuk a klienst
	if serviceClient, ok := c.serviceClients[folderNumber]; ok {
		return serviceClient.GrpcClient
	}

	log.WithFields(log.Fields{
		"error":     errorNoConnection,
		"swampName": swampName.Get(),
	}).Error("error while getting service client by swamp name")

	return nil

}

// GetServiceClientAndHost returns the full HydrAIDE service client wrapper for a given Swamp name.
//
// Unlike GetServiceClient(), which only returns the raw gRPC client,
// this method provides additional metadata, such as the Host identifier for the target server.
//
// Internally:
// - Computes the folder number by hashing the Swamp name
// - Resolves the target server from the folder-to-client map
//
// Returns:
// - *ServiceClient struct, which contains:
//   - `GrpcClient` → the actual gRPC HydrAIDEServiceClient
//   - `Host`       → the Host string of the resolved server (e.g. IP:port or logical name)
//
// - nil if no matching server is registered for the calculated folder
//
// Example:
//
//	swamp := name.New().Sanctuary("users").Realm("logs").Swamp("user123") \n client := hydraClient.GetServiceClientAndHost(swamp) \n if client != nil {\n    res, _ := client.GrpcClient.Read(...)\n    fmt.Println(\"Resolved Host:\", client.Host)\n}
//
// This is especially useful when:
// - Grouping Swamps by target server
// - Logging or debugging routing behavior
// - Performing multi-swamp operations where server affinity matters
//
// Note:
// - Thread-safe via internal read lock
// - The folder number is derived from the full Swamp name and `allFolders` total
func (c *client) GetServiceClientAndHost(swampName name.Name) *ServiceClient {

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

// GetUniqueServiceClients returns all unique HydrAIDE service clients
// only for internal use.
func (c *client) GetUniqueServiceClients() []hydraidepbgo.HydraideServiceClient {

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.uniqueServices

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
		"Host": ip,
		"out":  string(out),
	}).Info("pinging the Host")
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
				"Host": hostnameOrIP,
			}).Info("the Host ping without error")
		} else {
			log.WithFields(log.Fields{
				"Host": hostnameOrIP,
			}).Warning("the Host does not ping")
		}

	} else {
		ip, err := resolveHostname(hostnameOrIP)
		if err != nil {
			log.WithFields(log.Fields{
				"Host": hostnameOrIP,
				"err":  err,
			}).Error("could not resolve hostname")
		}

		// If input is an IP address, just ping it
		if ping(ip) {
			log.WithFields(log.Fields{
				"Host": ip,
			}).Info("the Host ping without error")
		} else {
			log.WithFields(log.Fields{
				"Host": ip,
			}).Warning("the Host does not ping")
		}
	}
}
