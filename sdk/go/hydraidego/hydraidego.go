// Package hydraidego
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
package hydraidego

import (
	"context"
	"fmt"
	"github.com/hydraide/hydraide/generated/hydraidepbgo"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/client"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	errorMessageConnectionError   = "connection error"
	errorMessageCtxTimeout        = "context timeout exceeded"
	errorMessageCtxClosedByClient = "context closed by client"
	errorMessageInvalidArgument   = "invalid argument"
	errorMessageNotFound          = "sanctuary not found"
	errorMessageUnknown           = "unknown error"
)

type Hydraidego interface {
	Heartbeat(ctx context.Context) error
	RegisterSwamp(ctx context.Context, request *RegisterSwampRequest) []error
}

type RegisterSwampRequest struct {
	// SwampPattern defines the pattern to register in HydrAIDE.
	// You can use wildcards (*) for dynamic parts.
	//
	// If the pattern includes a wildcard, HydrAIDE registers it on all servers,
	// since it cannot predict where the actual Swamp will reside after resolution.
	//
	// If the pattern has no wildcard, HydrAIDE uses its internal logic to determine
	// which server should handle the Swamp, and registers it only there.
	//
	// Example (no wildcard): Sanctuary("users").Realm("logs").Swamp("johndoe")
	// → Registered only on one server.
	//
	// Example (with wildcard): Sanctuary("users").Realm("logs").Swamp("*")
	// → Registered on all servers to ensure universal match.
	SwampPattern name.Name

	// CloseAfterIdle defines the idle time (inactivity period) after which
	// the Swamp is automatically closed and flushed from memory.
	//
	// When this timeout expires, HydrAIDE will:
	// - flush all changes to disk (if persistent),
	// - unload the Swamp from RAM,
	// - release any temporary resources.
	//
	// This helps keep memory lean and ensures disk durability when needed.
	CloseAfterIdle time.Duration

	// IsInMemorySwamp controls whether the Swamp should exist only in memory.
	//
	// If true → Swamp data is volatile and will be lost when closed.
	// If false → Swamp data is also persisted to disk.
	//
	// In-memory Swamps are ideal for:
	// - transient data between services,
	// - ephemeral socket messages,
	// - disappearing chat messages,
	// or any short-lived data flow.
	//
	// ⚠️ Warning: For in-memory Swamps, if CloseAfterIdle triggers,
	// all data is permanently lost.
	IsInMemorySwamp bool

	// FilesystemSettings provides persistence-related configuration.
	//
	// This is ignored if IsInMemorySwamp is true.
	// If persistence is enabled (IsInMemorySwamp = false), these settings control:
	// - how often data is flushed to disk,
	// - how large each chunk file can grow.
	//
	// If nil, the server will use its default settings.
	FilesystemSettings *SwampFilesystemSettings
}

type SwampFilesystemSettings struct {

	// WriteInterval defines how often (in seconds) HydrAIDE should write
	// new, modified, or deleted Treasures from memory to disk.
	//
	// If the Swamp is closed before this interval expires, it will still flush all data.
	// This setting optimizes for SSD wear vs. durability:
	// - Short intervals = safer but more writes.
	// - Longer intervals = fewer writes, but higher risk if crash occurs.
	//
	// Minimum allowed value is 1 second.
	WriteInterval time.Duration

	// MaxFileSize defines the maximum compressed chunk size on disk.
	//
	// Once this size is reached, a new chunk is created for further writes.
	//
	// This prevents large file rewrites, which can damage SSDs over time.
	// Smaller sizes → more files, better endurance.
	// Larger sizes → fewer files, better read performance for rarely-changed data.
	//
	// ⚠️ Always ensure MaxFileSize is larger than the filesystem block size.
	// HydrAIDE automatically compresses data, so this refers to the compressed size.
	MaxFileSize int
}

type hydraidego struct {
	client client.Client
}

func New(client client.Client) Hydraidego {
	return &hydraidego{
		client: client,
	}
}

// Heartbeat checks if all HydrAIDE servers are reachable.
// If any server is unreachable, it returns an aggregated error.
// If all are reachable, it returns nil.
//
// This method can be used to monitor the health of your HydrAIDE cluster.
// However, note that HydrAIDE clients have automatic reconnection logic,
// so a temporary network issue may not surface unless it persists.
func (h *hydraidego) Heartbeat(ctx context.Context) error {

	// Retrieve all unique gRPC service clients from the internal client pool.
	serviceClients := h.client.GetUniqueServiceClients()

	// Collect any errors encountered during heartbeat checks.
	allErrors := make([]string, 0)

	// Iterate through each server and perform a heartbeat ping.
	for _, serviceClient := range serviceClients {
		_, err := serviceClient.Heartbeat(ctx, &hydraidepbgo.HeartbeatRequest{
			Ping: "ping",
		})

		// If an error occurred, add it to the collection.
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("error: %v", err))
		}
	}

	// If any servers failed to respond, return a formatted error containing all issues.
	if len(allErrors) > 0 {
		return fmt.Errorf("one or many servers are not reachable: %v", allErrors)
	}

	// All servers responded successfully — return nil to indicate success.
	return nil
}

// RegisterSwamp registers a Swamp pattern across the appropriate HydrAIDE servers.
//
// This method is required before using a Swamp. It tells HydrAIDE how to handle
// memory, persistence, and routing for that pattern.
//
//   - If the SwampPattern contains any wildcard (e.g. Sanctuary("*"), Realm("*"), or Swamp("*")),
//     the pattern is registered on **all** servers.
//   - If the pattern is exact (no wildcard at any level), it is registered **only on the responsible server**,
//     based on HydrAIDE's internal name-to-folder mapping.
//
// ⚠️ While wildcarding the Sanctuary is technically possible, it is not recommended,
// as Sanctuary represents a high-level logical domain and should remain stable.
//
// Returns a list of errors, one for each server where registration failed.
// If registration is fully successful, it returns nil.
func (h *hydraidego) RegisterSwamp(ctx context.Context, request *RegisterSwampRequest) []error {

	// Container to collect any errors during registration.
	allErrors := make([]error, 0)

	// Validate that SwampPattern is provided.
	if request.SwampPattern == nil {
		allErrors = append(allErrors, fmt.Errorf("SwampPattern is required"))
		return allErrors
	}

	// List of servers where the Swamp pattern will be registered.
	selectedServers := make([]hydraidepbgo.HydraideServiceClient, 0)

	// Wildcard patterns must be registered on all servers,
	// because we don’t know in advance which server will handle each resolved Swamp.
	if request.SwampPattern.IsWildcardPattern() {
		selectedServers = h.client.GetUniqueServiceClients()
	} else {
		// For non-wildcard patterns, we determine the responsible server
		// using HydrAIDE’s name-based routing logic.
		selectedServers = append(selectedServers, h.client.GetServiceClient(request.SwampPattern))
	}

	// Iterate through the selected servers and register the Swamp on each.
	for _, serviceClient := range selectedServers {

		// Construct the RegisterSwampRequest payload for the gRPC call.
		rsr := &hydraidepbgo.RegisterSwampRequest{
			SwampPattern:    request.SwampPattern.Get(),
			CloseAfterIdle:  int64(request.CloseAfterIdle.Seconds()),
			IsInMemorySwamp: request.IsInMemorySwamp,
		}

		// If the Swamp is persistent (not in-memory), apply filesystem settings.
		if !request.IsInMemorySwamp && request.FilesystemSettings != nil {
			wi := int64(request.FilesystemSettings.WriteInterval.Seconds())
			mfs := int64(request.FilesystemSettings.MaxFileSize)
			rsr.WriteInterval = &wi
			rsr.MaxFileSize = &mfs
		}

		// Attempt to register the Swamp pattern on the current server.
		_, err := serviceClient.RegisterSwamp(ctx, rsr)

		// Handle any errors returned from the gRPC call.
		if err != nil {
			if s, ok := status.FromError(err); ok {
				switch s.Code() {
				case codes.Unavailable:
					allErrors = append(allErrors, NewError(ErrCodeConnectionError, errorMessageConnectionError))
				case codes.DeadlineExceeded:
					allErrors = append(allErrors, NewError(ErrCodeCtxTimeout, errorMessageCtxTimeout))
				case codes.Canceled:
					allErrors = append(allErrors, NewError(ErrCodeCtxClosedByClient, errorMessageCtxClosedByClient))
				case codes.InvalidArgument:
					allErrors = append(allErrors, NewError(ErrCodeInvalidArgument, fmt.Sprintf("%s: %v", errorMessageInvalidArgument, s.Message())))
				case codes.NotFound:
					allErrors = append(allErrors, NewError(ErrCodeNotFound, fmt.Sprintf("%s: %v", errorMessageNotFound, s.Message())))
				default:
					allErrors = append(allErrors, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err)))
				}
			} else {
				allErrors = append(allErrors, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err)))
			}
		}
	}

	// If any server failed, return the list of errors.
	if len(allErrors) > 0 {
		return allErrors
	}

	// All servers responded successfully – registration complete.
	return nil
}

// Lock acquires a distributed business-level lock for a specific domain/key.
//
// This is not tied to a single Swamp or Treasure — it’s a **cross-cutting domain lock**.
// You can use it to serialize logic across services or modules that operate
// on the same logical entity (e.g. user ID, order ID, transaction flow).
//
// 🧠 Ideal for scenarios like:
// - Credit transfers between users
// - Order/payment processing pipelines
// - Any sequence of operations where **no other process** should interfere
//
// ⚠️ Locking is **not required** for general Swamp access, reads, or standard writes.
// Use it **only** when your logic depends on critical, exclusive execution.
// Example: You want to deduct 10 credits from UserA and add it to UserB —
// and no other process should modify either user’s balance until this is done.
//
// ⚠️ This is a blocking lock — your flow will **wait** until the lock becomes available.
// The lock is acquired only when no other process holds it.
//
// ➕ The `ttl` ensures the system is self-healing:
// If a client crashes or forgets to unlock, the lock is **automatically released** after the TTL expires.
//
// ⏳ Important context behavior:
//   - If another client holds the lock, your request will block until it's released.
//   - If you set a context timeout or deadline, **make sure it's long enough** for the other process
//     to finish and call `Unlock()` — otherwise you may get a context timeout before acquiring the lock.
//
// ⚠️ The lock is issued **only on the first server**, to ensure consistency across distributed setups.
//
// Parameters:
//   - key:     Unique string representing the business domain to lock (e.g. "user:1234:credit")
//   - ttl:     Time-to-live for the lock. If not unlocked manually, it's auto-released after this duration.
//
// Returns:
// - lockID:   A unique identifier for the acquired lock — must be passed to `Unlock()`.
// - err:      Error if the lock could not be acquired, or if the context expired.
func (h *hydraidego) Lock(ctx context.Context, key string, ttl time.Duration) (lockID string, err error) {

	// Get available servers
	serverClients := h.client.GetUniqueServiceClients()

	// Always acquire business-level locks from the first server for consistency
	response, err := serverClients[0].Lock(ctx, &hydraidepbgo.LockRequest{
		Key: key,
		TTL: ttl.Milliseconds(),
	})

	// Handle network and gRPC-specific errors
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Unavailable:
				return "", NewError(ErrCodeConnectionError, errorMessageConnectionError)
			case codes.DeadlineExceeded:
				return "", NewError(ErrCodeCtxTimeout, errorMessageCtxTimeout)
			case codes.Canceled:
				return "", NewError(ErrCodeCtxClosedByClient, errorMessageCtxClosedByClient)
			default:
				return "", NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return "", NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	// Defensive check in case server returns an empty lockID
	if response.GetLockID() == "" {
		return "", NewError(ErrCodeNotFound, "lock ID not found")
	}

	// Successfully acquired the lock
	return response.GetLockID(), nil
}

// Unlock releases a previously acquired business-level lock using the lock ID.
//
// This is the counterpart of `Lock()`, and must be called once the critical section ends.
// The lock is matched by both key and lock ID, ensuring safety even in multi-client flows.
//
// ⚠️ Unlock always targets the first server — consistency is maintained at the entry point.
//
// Parameters:
// - key:     Same key used during locking (e.g. "user:1234:credit")
// - lockID:  The unique lock identifier returned by Lock()
//
// Returns:
// - err:     If the lock was not found, already released, or an error occurred during release.
func (h *hydraidego) Unlock(ctx context.Context, key string, lockID string) error {

	// Get available servers
	serverClients := h.client.GetUniqueServiceClients()

	_, err := serverClients[0].Unlock(ctx, &hydraidepbgo.UnlockRequest{
		Key:    key,
		LockID: lockID,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Unavailable:
				return NewError(ErrCodeConnectionError, errorMessageConnectionError)
			case codes.DeadlineExceeded:
				return NewError(ErrCodeCtxTimeout, errorMessageCtxTimeout)
			case codes.Canceled:
				return NewError(ErrCodeCtxClosedByClient, errorMessageCtxClosedByClient)
			case codes.NotFound:
				return NewError(ErrCodeNotFound, "key, or lock ID not found")
			default:
				return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	// Lock released successfully
	return nil

}
