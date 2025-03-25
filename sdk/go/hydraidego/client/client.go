package client

import (
	"context"
	"github.com/hydraide/hydraide/generated/go/hydraidepb"
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
	maxMsgSize        = 1024 * 1024 * 1024 // 1GB
	errorNoConnection = "there is no connection to the hydra server"
	errorNoService    = "there is no service client"
)

type Client interface {
	Connect(connectionAnalysis bool) error
	CloseConnection()
	GetServiceClient() hydrapb.HydraServiceClient
}

type client struct {
	conn          *grpc.ClientConn
	serviceClient hydrapb.HydraServiceClient
	// serverAddress format like: 192.168.106.100:4444 or localhost:4444 (for testing)
	serverAddress string
	mu            sync.Mutex
	certFile      string
}

// New creates a new instance of the client
func New(serverAddress string, certFile string) Client {
	return &client{
		serverAddress: serverAddress,
		certFile:      certFile,
	}
}

// Connect serviceName like hydrapb.HydraService
func (c *client) Connect(connectionAnalysis bool) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if connectionAnalysis {
		// try to ping the host name
		pingHost(c.serverAddress)
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stdout, os.Stderr, os.Stderr, 99))
	}

	serviceConfigJSON := `{
		  "methodConfig": [{
		    "name": [{"service": "hydrapb.HydraService"}],
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

	//... egyéb kód
	creds, certErr := credentials.NewClientTLSFromFile(c.certFile, "")
	if certErr != nil {
		log.WithFields(log.Fields{
			"error": certErr,
		}).Error("error while loading TLS credentials")
		return certErr
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxMsgSize)))
	opts = append(opts, grpc.WithDefaultServiceConfig(serviceConfigJSON))
	// Keepalive opció hozzáadása: ez segít, hogy a kapcsolat ne szakadjon meg inaktivitás miatt.
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                60 * time.Second, // Milyen gyakran küldjünk ping-et, ha nincs adatforgalom.
		Timeout:             10 * time.Second, // Meddig várunk válaszra, mielőtt timeout-olnánk.
		PermitWithoutStream: false,            // Még akkor is küldjünk ping-et, ha nincs aktív stream.
	}))

	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.NewClient(c.serverAddress, opts...)
	if err != nil {
		log.WithFields(log.Fields{
			"serverAddress": c.serverAddress,
			"error":         err,
		}).Error("error while connecting to the server")
		return err
	}

	c.conn = conn
	c.serviceClient = hydrapb.NewHydraServiceClient(c.conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := c.serviceClient.Heartbeat(ctx, &hydrapb.HeartbeatRequest{Ping: "beat"})
	if err != nil || pong.Pong != "beat" {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error while sending heartbeat request")
		return err // return error
	}

	log.WithFields(log.Fields{
		"serverAddress": c.serverAddress,
	}).Info("connected to the hydra server successfully")

	return nil

}

func (c *client) CloseConnection() {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.WithFields(log.Fields{
				"serverAddress": c.serverAddress,
				"error":         err,
			}).Error("error while closing connection")
		}
	}

}

func (c *client) GetServiceClient() hydrapb.HydraServiceClient {

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		log.WithFields(log.Fields{
			"error": errorNoConnection,
		}).Error("error while getting service client")
		return nil
	}
	if c.serviceClient == nil {
		log.WithFields(log.Fields{
			"error": errorNoService,
		}).Error("error while getting service client")
		return nil
	}

	return c.serviceClient

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
