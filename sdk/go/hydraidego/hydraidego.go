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
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/hydraide/hydraide/generated/hydraidepbgo"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/client"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"sync"
	"time"
)

const (
	errorMessageConnectionError     = "connection error"
	errorMessageCtxTimeout          = "context timeout exceeded"
	errorMessageCtxClosedByClient   = "context closed by client"
	errorMessageInvalidArgument     = "invalid argument"
	errorMessageNotFound            = "sanctuary not found"
	errorMessageUnknown             = "unknown error"
	errorMessageSwampNameNotCorrect = "swamp name is not correct"
	errorMessageSwampNotFound       = "swamp not found"
	errorMessageInternalError       = "internal error"
	errorMessageKeyAlreadyExists    = "key already exists"
)

const (
	tagHydrAIDE  = "hydraide"
	tagKey       = "key"
	tagValue     = "value"
	tagOmitempty = "omitempty"
	tagCreatedAt = "createdAt"
	tagCreatedBy = "createdBy"
	tagUpdatedAt = "updatedAt"
	tagUpdatedBy = "updatedBy"
	tagExpireAt  = "expireAt"
)

type Hydraidego interface {
	Heartbeat(ctx context.Context) error
	RegisterSwamp(ctx context.Context, request *RegisterSwampRequest) []error
	DeRegisterSwamp(ctx context.Context, swampName name.Name) []error
	Lock(ctx context.Context, key string, ttl time.Duration) (lockID string, err error)
	Unlock(ctx context.Context, key string, lockID string) error
	IsSwampExist(ctx context.Context, swampName name.Name) (bool, error)
	IsKeyExists(ctx context.Context, swampName name.Name, key string) (bool, error)
	CatalogCreate(ctx context.Context, swampName name.Name, model any) error
	CatalogCreateMany(ctx context.Context, swampName name.Name, models []any, iterator CreateManyIteratorFunc) error
	CatalogCreateManyToMany(ctx context.Context, request []*CatalogManyToManyRequest, iterator CatalogCreateManyToManyIteratorFunc) error
	CatalogRead(ctx context.Context, swampName name.Name, key string, model any) error
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

// DeRegisterSwamp removes a previously registered Swamp pattern from the relevant HydrAIDE server(s).
//
// 🧠 This is the **counterpart of RegisterSwamp()**, and follows the same routing logic:
//   - If the pattern includes any wildcard (e.g. Sanctuary("*"), Realm("*"), or Swamp("*")),
//     the deregistration is propagated to **all** servers.
//   - If the pattern is fully qualified (no wildcards), the request is routed to the exact responsible server
//     using HydrAIDE’s O(1) folder mapping logic.
//
// 🔥 Important notes:
//
// This function **does not delete any data** or existing Swamps — it only removes the **pattern registration**
// from the internal registry. This affects how future pattern-based operations behave.
//
// ✅ **When should you use this?**
//   - When you're deprecating a pattern **permanently**, e.g. restructuring your domain logic.
//   - When you're **migrating** from one pattern to another, and want to avoid potential pattern conflicts.
//   - When your team changes the logic of how logs, sessions, credits, etc. are stored,
//     and you want to cleanly retire the old pattern.
//
// ⚠️ **When should you NOT use this?**
//   - If a Swamp is just temporarily inactive or empty — it will unload itself automatically.
//     There is no need to deregister unless you're redesigning structure.
//
// 🛠️ Typical migration flow:
// 1. Migrate existing data to a new Swamp pattern
// 2. Delete the old Swamp's Treasures (using Delete or DeleteAll)
// 3. Finally, call `DeRegisterSwamp()` to remove the pattern itself
//
// ❗ If you skip step 2, the Swamp files may remain on disk even if the pattern is gone.
//
// Returns:
// - A list of errors if deregistration fails on any server
// - Nil if deregistration completes successfully across all relevant servers
func (h *hydraidego) DeRegisterSwamp(ctx context.Context, swampName name.Name) []error {

	// Container to collect any errors during deregistration.
	allErrors := make([]error, 0)

	// Validate that SwampPattern is provided.
	if swampName == nil {
		allErrors = append(allErrors, fmt.Errorf("SwampPattern is required"))
		return allErrors
	}

	// List of servers where the Swamp pattern will be registered.
	selectedServers := make([]hydraidepbgo.HydraideServiceClient, 0)

	// Wildcard patterns must be registered on all servers,
	// because we don’t know in advance which server will handle each resolved Swamp.
	if swampName.IsWildcardPattern() {
		selectedServers = h.client.GetUniqueServiceClients()
	} else {
		// For non-wildcard patterns, we determine the responsible server
		// using HydrAIDE’s name-based routing logic.
		selectedServers = append(selectedServers, h.client.GetServiceClient(swampName))
	}

	// Iterate through the selected servers and register the Swamp on each.
	for _, serviceClient := range selectedServers {

		// Construct the RegisterSwampRequest payload for the gRPC call.
		rsr := &hydraidepbgo.DeRegisterSwampRequest{
			SwampPattern: swampName.Get(),
		}

		// Attempt to Deregister the Swamp pattern on the current server.
		_, err := serviceClient.DeRegisterSwamp(ctx, rsr)

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

// IsSwampExist checks whether a specific Swamp currently exists in the HydrAIDE system.
//
// ⚠️ This is a **direct existence check** – it does NOT accept wildcards or patterns.
// You must provide a fully resolved Swamp name (Sanctuary + Realm + Swamp).
//
// ✅ When to use this:
// - When you want to check if a Swamp was previously created by another process
// - When a Swamp may have been deleted automatically (e.g., became empty)
// - When you want to determine Swamp presence **without hydrating or loading data**
// - As part of fast lookups, hydration conditionals, or visibility toggles
//
// 🔍 **Real-world example**:
// Suppose you're generating AI analysis per domain and storing them in separate Swamps:
//
//	Sanctuary("domains").Realm("ai").Swamp("trendizz.com")
//	Sanctuary("domains").Realm("ai").Swamp("hydraide.io")
//
// When rendering a UI list of domains, you don’t want to load full AI data.
// Instead, use `IsSwampExist()` to check if an AI analysis exists for each domain,
// and show a ✅ or ❌ icon accordingly — without incurring I/O or memory cost.
//
// ⚙️ Behavior:
// - If the Swamp exists → returns (true, nil)
// - If it never existed or was auto-deleted → returns (false, nil)
// - If a server error occurs → returns (false, error)
//
// 🚀 This check is extremely fast: O(1) routing + metadata lookup.
// ➕ It does **not hydrate or load** the Swamp into memory — it only checks for existence on disk.
//
//	If the Swamp is already open, it stays open. If not, it stays closed.
//	This allows for high-frequency checks without affecting memory or system state.
//
// ⚠️ Requires that the Swamp pattern for the given name was previously registered.
func (h *hydraidego) IsSwampExist(ctx context.Context, swampName name.Name) (bool, error) {

	response, err := h.client.GetServiceClient(swampName).IsSwampExist(ctx, &hydraidepbgo.IsSwampExistRequest{
		SwampName: swampName.Get(),
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Unavailable:
				return false, NewError(ErrCodeConnectionError, errorMessageConnectionError)
			case codes.DeadlineExceeded:
				return false, NewError(ErrCodeCtxTimeout, errorMessageCtxTimeout)
			case codes.Canceled:
				return false, NewError(ErrCodeCtxClosedByClient, errorMessageCtxClosedByClient)
			case codes.InvalidArgument:
				return false, NewError(ErrCodeNotFound, fmt.Sprintf("%s: %v", errorMessageSwampNameNotCorrect, s.Message()))
			case codes.FailedPrecondition:
				return false, NewError(ErrCodeSwampNotFound, fmt.Sprintf("%s: %v", errorMessageSwampNotFound, s.Message()))
			default:
				return false, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return false, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	if !response.IsExist {
		return false, nil
	}

	return true, nil

}

// IsKeyExists checks whether a specific key exists inside a given Swamp.
//
// 🔍 This is a **memory-aware** check — the Swamp is always hydrated (loaded into memory) as part of this operation.
// If the Swamp was not yet loaded, it will be loaded now and remain in memory based on its registered settings
// (e.g. CloseAfterIdle, persistence, etc.).
//
// ✅ When to use this:
// - When you want to check if a given key has been previously inserted into a Swamp
// - When you're implementing **unique key checks**, deduplication, or conditional inserts
// - When performance is critical and the Swamp is expected to be open or heavily reused
//
// ⚠️ Difference from `IsSwampExist()`:
// `IsSwampExist()` checks for Swamp presence on disk **without hydration**
// `IsKeyExists()` loads the Swamp and searches for the exact key
//
// 🧠 Real-world example:
// In Trendizz.com’s crawler, we keep domain-specific Swamps (e.g. `.hu`, `.de`, `.fr`) open in memory.
// Each Swamp contains a list of already-seen domains.
// Before crawling a new domain, we call `IsKeyExists()` to check if it's already indexed.
// This lets us skip unnecessary work and ensures we don't reprocess the same domain twice.
//
// 🔁 Return values:
// - `(true, nil)` → Swamp and key both exist
// - `(false, nil)` → Swamp exists, but key does not
// - `(false, ErrCodeSwampNotFound)` → Swamp does not exist
// - `(false, <other error>)` → Some database/server issue occurred
//
// ⚠️ Always use **fully qualified Swamp names** – no wildcards allowed.
// If the Swamp was not registered, or was deleted due to being empty, this will return `ErrCodeSwampNotFound`.
func (h *hydraidego) IsKeyExists(ctx context.Context, swampName name.Name, key string) (bool, error) {

	response, err := h.client.GetServiceClient(swampName).IsKeyExist(ctx, &hydraidepbgo.IsKeyExistRequest{
		SwampName: swampName.Get(),
		Key:       key,
	})

	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.Unavailable:
				return false, NewError(ErrCodeConnectionError, errorMessageConnectionError)
			case codes.DeadlineExceeded:
				return false, NewError(ErrCodeCtxTimeout, errorMessageCtxTimeout)
			case codes.Canceled:
				return false, NewError(ErrCodeCtxClosedByClient, errorMessageCtxClosedByClient)
			case codes.InvalidArgument:
				return false, NewError(ErrCodeNotFound, fmt.Sprintf("%s: %v", errorMessageInvalidArgument, s.Message()))
			case codes.FailedPrecondition:
				return false, NewError(ErrCodeSwampNotFound, fmt.Sprintf("%s: %v", errorMessageNotFound, s.Message()))
			default:
				return false, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return false, NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	if !response.IsExist {
		return false, nil
	}

	return true, nil

}

// CatalogCreate inserts a new Treasure into a Swamp using a tagged Go struct as the input model.
//
// 🧠 Purpose:
// This is a "catalog-style" insert — ideal for Swamps that hold a large number of similar entities,
// like users, transactions, credit logs, etc. Think of it as writing a row into a virtual table.
//
// ✅ Behavior:
// - Creates the Swamp if it does not exist yet
// - Inserts the provided key-value pair **only if the key does not already exist**
// - Will return an error if the key already exists
//
// 📦 Model requirements:
// - The `model` must be a **pointer to a struct**
// - The struct must contain a `hydraide:"key"` field with a non-empty string
// - Optionally, it may contain:
//
//   - `hydraide:"value"` → the main value of the Treasure.
//
//     ✅ Supported types:
//
//   - string, bool
//
//   - int8, int16, int32, int64
//
//   - uint8, uint16, uint32, uint64
//
//   - float32, float64
//
//   - struct or pointer to struct (automatically GOB-encoded)
//
//     🔬 Best practice:
//     Always use the **smallest suitable numeric type**.
//     For example: prefer `uint8` or `int16` over `int`.
//     HydrAIDE stores values in raw binary form — so smaller types directly reduce
//     memory usage and disk space.
//
//   - `hydraide:"expireAt"`   → expiration logic (time.Time)
//
//   - `hydraide:"createdBy"`  → who created it (string)
//
//   - `hydraide:"createdAt"`  → when it was created (time.Time)
//
//   - `hydraide:"updatedBy"`  → optional metadata
//
//   - `hydraide:"updatedAt"`  → optional metadata
//
// ✨ Example use case 1:
// You store user records in a Swamp:
//
//	Sanctuary("system").Realm("users").Swamp("all")
//
// Each call to `CatalogCreate()` adds a new user — uniquely identified by `UserUUID`.
//
// ✨ Example use case 2 (real-world):
// In Trendizz.com’s domain crawler, we store known domains in Swamps per TLD.
// Instead of first checking if a domain exists, we call `CatalogCreate()` directly.
// If the domain already exists → we receive `ErrCodeAlreadyExists`.
// If it doesn’t → it is inserted in one step.
// This saves a read roundtrip and simplifies the control flow.
//
// 🔁 Return values:
// - `nil` → success, insert completed
// - `ErrCodeAlreadyExists` → key already exists in the Swamp
// - `ErrCodeInvalidModel` → struct is invalid (e.g. not a pointer, missing tags)
// - Other database-level error codes if something went wrong
// Example: CreditLog model used with CatalogCreate()
// This model stores credit-related changes per user in a Swamp.
// Each record is identified by UserUUID and optionally enriched with metadata.
func (h *hydraidego) CatalogCreate(ctx context.Context, swampName name.Name, model any) error {

	kvPair, err := convertModelToKeyValuePair(model)
	if err != nil {
		return NewError(ErrCodeInvalidModel, err.Error())
	}

	// egyetlen adatot hozunk létre a hydrában overwrit nélkül.
	// A swamp mindenképpen létrejön, ha még nem létezett, de a kulcs csak akkor, ha még nem létezett
	setResponse, err := h.client.GetServiceClient(swampName).Set(ctx, &hydraidepbgo.SetRequest{
		Swamps: []*hydraidepbgo.SwampRequest{
			{
				SwampName: swampName.Get(),
				KeyValues: []*hydraidepbgo.KeyValuePair{
					kvPair,
				},
				CreateIfNotExist: true,
				Overwrite:        false,
			},
		},
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
			case codes.Internal:
				return NewError(ErrCodeInternalDatabaseError, fmt.Sprintf("%s: %v", errorMessageInternalError, s.Message()))
			default:
				return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	// ha a kulcs már létezik...
	for _, swamp := range setResponse.GetSwamps() {
		for _, kv := range swamp.GetKeysAndStatuses() {
			if kv.GetStatus() == hydraidepbgo.Status_NOTHING_CHANGED {
				return NewError(ErrCodeAlreadyExists, errorMessageKeyAlreadyExists)
			}
		}
	}

	return nil

}

type CreateManyIteratorFunc func(key string, err error) error

// CatalogCreateMany inserts multiple Treasures into a single Swamp using catalog-style logic.
//
// 🧠 Purpose:
// Use this when you want to batch-insert a list of records into a Swamp,
// each with its own unique key — for example, uploading multiple users, products,
// or log entries at once.
//
// ✅ Behavior:
// - Creates the Swamp if it does not exist yet
// - Converts each input model into a KeyValuePair using `convertModelToKeyValuePair()`
// - Inserts all items in a single SetRequest
// - Fails **only** if the gRPC call fails or if a model is invalid
//
// 📦 Model Requirements:
// Each element in `models` must be a pointer to a struct,
// and follow the same field tagging rules as in `CatalogCreate()`:
//   - `hydraide:"key"`     → required non-empty string
//   - `hydraide:"value"`   → optional value (primitive, struct, pointer)
//   - `hydraide:"createdAt"`, `expireAt`, etc. → optional metadata
//
// 🔁 Iterator (optional):
// You may provide an `iterator` function to handle per-record responses.
// It will be called for each inserted item with:
//
//	key string  → the unique Treasure key
//	err error   → nil if inserted successfully,
//	              `ErrCodeAlreadyExists` if the key was already present
//
// This allows you to track insert success/failure **per item**, without manually parsing the response.
//
// If `iterator` is `nil`, the function will insert all models silently, and return only global errors.
//
// ✨ Example use:
//
//	var users []any = []any{&User1, &User2, &User3}
//	err := client.CatalogCreateMany(ctx, name.Swamp("users", "all", "2025"), users, func(key string, err error) error {
//	    if err != nil {
//	        log.Printf("❌ failed to insert %s: %v", key, err)
//	    } else {
//	        log.Printf("✅ inserted: %s", key)
//	    }
//	    return nil
//	})
//
// 🧯 Error Handling:
// - If any model is invalid → `ErrCodeInvalidModel`
// - If the entire gRPC Set call fails → appropriate connection or database error
// - If a key already exists → passed back through the iterator as `ErrCodeAlreadyExists`
// - If no iterator is provided, duplicates are silently skipped
//
// 🔁 Return:
//   - `nil` if the operation succeeded and/or the iterator handled everything
//   - Any error returned by the iterator will abort processing and be returned immediately
//   - If the underlying gRPC Set request fails (e.g. connection error, database failure),
//     the function returns a global error (e.g. ErrCodeConnectionError, ErrCodeInternalDatabaseError, etc.)
func (h *hydraidego) CatalogCreateMany(ctx context.Context, swampName name.Name, models []any, iterator CreateManyIteratorFunc) error {

	kvPairs := make([]*hydraidepbgo.KeyValuePair, 0, len(models))

	for _, model := range models {
		kvPair, err := convertModelToKeyValuePair(model)
		if err != nil {
			return NewError(ErrCodeInvalidModel, err.Error())
		}
		kvPairs = append(kvPairs, kvPair)
	}

	setResponse, err := h.client.GetServiceClient(swampName).Set(ctx, &hydraidepbgo.SetRequest{
		Swamps: []*hydraidepbgo.SwampRequest{
			{
				SwampName:        swampName.Get(),
				KeyValues:        kvPairs,
				CreateIfNotExist: true,
				Overwrite:        false,
			},
		},
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
			case codes.Internal:
				return NewError(ErrCodeInternalDatabaseError, fmt.Sprintf("%s: %v", errorMessageInternalError, s.Message()))
			default:
				return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}
	}

	// process the response and start the iterator if the iterator is not nil
	if iterator != nil {
		for _, swamp := range setResponse.GetSwamps() {
			for _, kv := range swamp.GetKeysAndStatuses() {
				if kv.GetStatus() == hydraidepbgo.Status_NOTHING_CHANGED {
					if iterErr := iterator(kv.GetKey(), NewError(ErrCodeAlreadyExists, errorMessageKeyAlreadyExists)); iterErr != nil {
						return iterErr
					}
				} else {
					if iterErr := iterator(kv.GetKey(), nil); iterErr != nil {
						return iterErr
					}
				}
			}
		}
	}

	return nil

}

type CatalogCreateManyToManyIteratorFunc func(swampName name.Name, key string, err error) error

type CatalogManyToManyRequest struct {
	SwampName name.Name
	Models    []any
}

// CatalogCreateManyToMany inserts batches of catalog-style models across multiple Swamps,
// and intelligently groups requests per server to optimize communication.
//
// 🧠 Use this function when:
// - You want to insert data into multiple Swamps (e.g. one per user, per region, per domain)
// - The Swamps may be distributed across multiple HydrAIDE servers
// - You want to minimize the number of gRPC calls and batch multiple writes per server
//
// ✅ Behavior:
// - Groups all SwampRequests by their destination server (based on Swamp name hashing)
// - Sends **one SetRequest per server**, bundling all Swamps and KeyValuePairs
// - Converts each model using `convertModelToKeyValuePair()`
// - Automatically creates Swamps if they don't exist
// - Does **not overwrite existing keys**
//
// 🔁 Iterator function (optional):
//
//	If provided, it will be called for every inserted key, with:
//	- swampName: the Swamp where the key was written
//	- key: the actual key
//	- err: nil if success, or ErrCodeAlreadyExists if key already existed
//
//	You can use this to log, retry, or track insert results.
//
// Example:
//
//	requests := []*CatalogManyToManyRequest{
//	    {
//	        SwampName: name.New().Sanctuary("domains").Realm("ai").Swamp("hu"),
//	        Models: []any{...},
//	    },
//	    {
//	        SwampName: name.New().Sanctuary("domains").Realm("ai").Swamp("de"),
//	        Models: []any{...},
//	    },
//	    {
//	        SwampName: name.New().Sanctuary("domains").Realm("ai").Swamp("fr"),
//	        Models: []any{...},
//	    },
//	}
//
//	err := client.CatalogCreateManyToMany(ctx, requests, func(swamp name.Name, key string, err error) error {
//	    if err != nil {
//	        log.Printf("❌ failed to insert %s into %s: %v", key, swamp.Get(), err)
//	    } else {
//	        log.Printf("✅ inserted %s into %s", key, swamp.Get())
//	    }
//	    return nil
//	})
//
// 🔥 Ideal for:
// - Crawler results
// - Batch imports
// - Indexing pipelines
// - Data normalization jobs
//
// 🧯 Errors:
// - Any invalid model → `ErrCodeInvalidModel`
// - gRPC/connection errors → mapped to consistent SDK error codes
// - Iterator errors → if the callback returns a non-nil error, processing stops immediately
func (h *hydraidego) CatalogCreateManyToMany(ctx context.Context, request []*CatalogManyToManyRequest, iterator CatalogCreateManyToManyIteratorFunc) error {

	type requestGroup struct {
		client        hydraidepbgo.HydraideServiceClient
		swampRequests []*hydraidepbgo.SwampRequest
	}

	serverRequests := make(map[string]*requestGroup)

	for _, req := range request {

		// lekérdezzük a szewrver adatait a swamp neve alapján
		clientAndHost := h.client.GetServiceClientAndHost(req.SwampName)

		if _, ok := serverRequests[clientAndHost.Host]; !ok {
			serverRequests[clientAndHost.Host] = &requestGroup{
				client: clientAndHost.GrpcClient,
			}
		}

		kvPairs := make([]*hydraidepbgo.KeyValuePair, 0, len(req.Models))

		for _, model := range req.Models {
			kvPair, err := convertModelToKeyValuePair(model)
			if err != nil {
				return NewError(ErrCodeInvalidModel, err.Error())
			}
			kvPairs = append(kvPairs, kvPair)
		}

		serverRequests[clientAndHost.Host].swampRequests = append(serverRequests[clientAndHost.Host].swampRequests, &hydraidepbgo.SwampRequest{
			SwampName:        req.SwampName.Get(),
			KeyValues:        kvPairs,
			CreateIfNotExist: true,
			Overwrite:        false,
		})

	}

	for _, reqGroup := range serverRequests {

		setResponse, err := reqGroup.client.Set(ctx, &hydraidepbgo.SetRequest{
			Swamps: reqGroup.swampRequests,
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
				case codes.InvalidArgument:
					return NewError(ErrCodeInvalidArgument, fmt.Sprintf("%s: %v", errorMessageInvalidArgument, s.Message()))
				case codes.Internal:
					return NewError(ErrCodeInternalDatabaseError, fmt.Sprintf("%s: %v", errorMessageInternalError, s.Message()))
				default:
					return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
				}
			} else {
				return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		}

		// process the response and start the iterator if the iterator is not nil
		if iterator != nil {
			for _, swamp := range setResponse.GetSwamps() {
				swampNameObj := name.Load(swamp.GetSwampName())
				// végigmegyünk a kulcsokon és azok státuszán
				for _, kv := range swamp.GetKeysAndStatuses() {
					if kv.GetStatus() == hydraidepbgo.Status_NOTHING_CHANGED {
						if iterErr := iterator(swampNameObj, kv.GetKey(), NewError(ErrCodeAlreadyExists, errorMessageKeyAlreadyExists)); iterErr != nil {
							return iterErr
						}
					} else {
						if iterErr := iterator(swampNameObj, kv.GetKey(), nil); iterErr != nil {
							return iterErr
						}
					}
				}
			}
		}

	}

	return nil

}

// CatalogRead retrieves a single Treasure by key from the specified Swamp,
// and unmarshals the result into the provided Go model.
//
// 🧠 Use this function when:
// - You want to read a single key from a Swamp
// - You want the result directly mapped into a typed struct
// - You need reliable error codes (e.g. key not found, invalid model, etc.)
//
// ✅ Behavior:
// - Sends a GetRequest to the Hydra server responsible for the Swamp
// - Extracts the first returned Treasure (if exists)
// - Automatically unmarshals the value into the given struct using field tags
//   - Required: `hydraide:"key"`
//   - Optional: `hydraide:"value"`, `expireAt`, `createdBy`, etc.
//
// - Supports GOB-decoded slices, maps, pointers, and all primitive types
//
// 📌 Notes:
// - The model parameter must be a pointer to a struct
// - Only one Treasure is expected — if none found, returns ErrCodeNotFound
//
// 🔥 Ideal for:
// - Real-time lookups
// - Detail views
// - Conditional logic (e.g. check if user already exists)
//
// 🧯 Errors:
// - Key not found → `ErrCodeNotFound`
// - Invalid model or conversion error → `ErrCodeInvalidModel`
// - Swamp not found → `ErrCodeSwampNotFound`
// - Timeout / context / network issues → appropriate SDK error codes
func (h *hydraidego) CatalogRead(ctx context.Context, swampName name.Name, key string, model any) error {

	// a swampot és a kulcsot beállítjuk
	swamps := []*hydraidepbgo.GetSwamp{
		{
			SwampName: swampName.Get(),
			Keys:      []string{key},
		},
	}

	// lekérdezzüók az egyetlen kulcsot a hydrából
	response, err := h.client.GetServiceClient(swampName).Get(ctx, &hydraidepbgo.GetRequest{
		Swamps: swamps,
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
			case codes.FailedPrecondition:
				return NewError(ErrCodeSwampNotFound, fmt.Sprintf("%s: %v", errorMessageSwampNotFound, s.Message()))
			case codes.Internal:
				return NewError(ErrCodeInternalDatabaseError, fmt.Sprintf("%s: %v", errorMessageInternalError, s.Message()))
			default:
				return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
			}
		} else {
			return NewError(ErrCodeUnknown, fmt.Sprintf("%s: %v", errorMessageUnknown, err))
		}

	}

	for _, swamp := range response.GetSwamps() {
		for _, treasure := range swamp.GetTreasures() {
			if treasure.IsExist == false {
				return NewError(ErrCodeNotFound, "key not found")
			}
			if convErr := convertProtoTreasureToModel(treasure, model); convErr != nil {
				return NewError(ErrCodeInvalidModel, convErr.Error())
			}
			return nil
		}
	}

	// the treasure does not exist
	return NewError(ErrCodeNotFound, "key not found")

}

// convertModelToKeyValuePair converts a Go struct (passed as pointer) into a HydrAIDE-compatible KeyValuePair message.
//
// 🧠 This is an **internal serialization helper** used by the Go SDK to translate user-defined models
// into the binary format that HydrAIDE expects when inserting or updating Treasures.
//
// ✅ Supported field tags:
// - `hydraide:"key"`       → Marks the string field to use as the Treasure key (must be non-empty).
// - `hydraide:"value"`     → Marks the value field (can be any supported primitive or complex type).
// - `hydraide:"expireAt"`  → Optional `time.Time`, marks the logical expiry time of the Treasure.
// - `hydraide:"createdAt"` / `createdBy` / `updatedAt` / `updatedBy` → Optional metadata fields.
// - `hydraide:"omitempty"` → Skips the field during encoding if it's zero, nil, or empty.
//
// ✅ Supported value types:
// - Primitives: string, bool, int, uint, float (various widths)
// - time.Time (as int64 UNIX timestamp)
// - Slices and maps (serialized as GOB-encoded binary blobs)
// - Structs and pointers (also GOB-encoded)
// - `nil` / empty values are optionally excluded if marked with `omitempty`
//
// ⚠️ Requirements:
// - The input **must be a pointer to a struct**, otherwise the function returns an error.
// - The struct **must contain a field marked as `hydrun:"key"`** with a non-empty string.
// - The value can be a primitive or complex field marked with `hydrun:"value"`.
// - If no value is provided, the resulting KeyValuePair will include a `VoidVal=true` marker.
//
// 🧬 Why this matters:
// HydrAIDE works with protocol-level binary messages.
// Every Treasure must be sent as a KeyValuePair with a valid key and (optionally) a value.
// This function bridges Go structs and HydrAIDE’s native format, abstracting encoding logic.
//
// ✨ This is how arbitrary business models (e.g. `UserProfile`, `InvoiceItem`) are safely,
// efficiently and correctly transformed into Treasure representations.
//
// 📌 If you're building a new SDK (e.g. for Python, Rust, Node.js), your implementation
// should follow the same principles:
// - Tag-driven key/value separation
// - Support for void values and expiration
// - Metadata injection
// - Optional field skipping (e.g. omitempty)
// - Consistent type coercion for known value types
func convertModelToKeyValuePair(model any) (*hydraidepbgo.KeyValuePair, error) {

	// Get the reflection value of the input model
	v := reflect.ValueOf(model)

	// 🧪 Validate the input: it must be a pointer to a struct.
	// This is required because we'll be using reflection to iterate over the fields
	// and extract tags and values dynamically. Non-pointer or non-struct inputs are invalid.
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("input must be a pointer to a struct")
	}

	// This flag tracks whether any value has been set.
	// If no value is provided (only key or metadata), we'll later set VoidVal = true.
	valueVoid := true

	// Initialize the KeyValuePair that will hold the final encoded output
	kvPair := &hydraidepbgo.KeyValuePair{}

	// Get the actual struct (dereferenced value) and its type
	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)

		// Check if the field has a `hydraide:"omitempty"` tag,
		// and skip it if the value is considered "empty" (zero, nil, blank, etc.)
		if tag, ok := field.Tag.Lookup(tagHydrAIDE); ok && tag == tagOmitempty {
			value := v.Field(i)

			// Evaluate "emptiness" based on Go's zero-value semantics per type
			// Strings → must not be empty
			// Pointers → must not be nil
			// Numbers → must not be zero
			// Slices/Maps → must not be nil or empty
			// time.Time → must not be zero (uninitialized)

			if (value.Kind() == reflect.String && value.String() == "") ||
				(value.Kind() == reflect.Ptr && value.IsNil()) ||
				(value.Kind() == reflect.Int8 && value.Int() == 0) ||
				(value.Kind() == reflect.Int16 && value.Int() == 0) ||
				(value.Kind() == reflect.Int32 && value.Int() == 0) ||
				(value.Kind() == reflect.Int64 && value.Int() == 0) ||
				(value.Kind() == reflect.Int && value.Int() == 0) ||
				(value.Kind() == reflect.Uint8 && value.Uint() == 0) ||
				(value.Kind() == reflect.Uint16 && value.Uint() == 0) ||
				(value.Kind() == reflect.Uint32 && value.Uint() == 0) ||
				(value.Kind() == reflect.Uint64 && value.Uint() == 0) ||
				(value.Kind() == reflect.Uint && value.Uint() == 0) ||
				(value.Kind() == reflect.Float32 && value.Float() == 0) ||
				(value.Kind() == reflect.Float64 && value.Float() == 0) ||
				(value.Kind() == reflect.Slice && (value.IsNil() || value.Len() == 0)) ||
				(value.Kind() == reflect.Map && (value.IsNil() || value.Len() == 0)) ||
				(value.Type() == reflect.TypeOf(time.Time{}) && value.Interface().(time.Time).IsZero()) {

				// If the field is empty, skip further processing and continue to the next field
				continue
			}
		}

		// Check if the current field is marked as the `key` field (via `hydraide:"key"` tag)
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagKey {

			value := v.Field(i)

			// Validate that the field is a non-empty string — required for all HydrAIDE Treasures.
			// Keys must always be explicit and unique within a Swamp.
			if value.Kind() == reflect.String && value.String() != "" {
				// Found the key — assign it to the KeyValuePair
				kvPair.Key = value.String()
				valueVoid = false
				continue
			}

			// If the key field is missing or empty, this is an invalid model
			return nil, errors.New("key field must be a non-empty string")
		}

		// Check if the current field is tagged as the `value` field (via `hydraide:"value"`)
		// This field holds the actual value of the Treasure.
		// We detect its type using reflection and populate the corresponding proto field in KeyValuePair.
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagValue {

			value := v.Field(i)

			switch value.Kind() {

			// 🧵 Simple primitives (string, bool, numbers)
			case reflect.String:
				stringVal := value.String()
				kvPair.StringVal = &stringVal
				valueVoid = false
				continue

			case reflect.Bool:
				// HydrAIDE uses a custom Boolean enum to allow storing `false` values explicitly
				boolVal := hydraidepbgo.Boolean_FALSE
				if value.Bool() {
					boolVal = hydraidepbgo.Boolean_TRUE
				}
				kvPair.BoolVal = &boolVal
				valueVoid = false
				continue

			// 🧮 Unsigned integers
			case reflect.Uint8, reflect.Uint16, reflect.Uint32:
				intVal := uint32(value.Uint())
				switch value.Kind() {
				case reflect.Uint8:
					kvPair.Uint8Val = &intVal
				case reflect.Uint16:
					kvPair.Uint16Val = &intVal
				case reflect.Uint32:
					kvPair.Uint32Val = &intVal
				}
				valueVoid = false
				continue

			case reflect.Uint64:
				intVal := value.Uint()
				kvPair.Uint64Val = &intVal
				valueVoid = false
				continue

			// 🔢 Signed integers
			case reflect.Int8, reflect.Int16, reflect.Int32:
				intVal := int32(value.Int())
				switch value.Kind() {
				case reflect.Int8:
					kvPair.Int8Val = &intVal
				case reflect.Int16:
					kvPair.Int16Val = &intVal
				case reflect.Int32:
					kvPair.Int32Val = &intVal
				}
				valueVoid = false
				continue

			case reflect.Int, reflect.Int64:
				intVal := value.Int()
				kvPair.Int64Val = &intVal
				valueVoid = false
				continue

			// 🔬 Floating point numbers
			case reflect.Float32:
				floatVal := float32(value.Float())
				kvPair.Float32Val = &floatVal
				valueVoid = false
				continue

			case reflect.Float64:
				floatVal := value.Float()
				kvPair.Float64Val = &floatVal
				valueVoid = false
				continue

			// 🧱 Complex binary types – slices, maps, pointers, structs (excluding time)
			case reflect.Slice:

				// Special case for []byte → raw binary value
				if value.Type().Elem().Kind() == reflect.Uint8 {
					kvPair.BytesVal = value.Bytes()
					valueVoid = false
				} else {
					// All other slices are GOB-encoded
					registerGobTypeIfNeeded(value.Interface())
					var buf bytes.Buffer
					encoder := gob.NewEncoder(&buf)
					if err := encoder.Encode(value.Interface()); err != nil {
						return nil, fmt.Errorf("could not GOB-encode slice: %w", err)
					}
					kvPair.BytesVal = buf.Bytes()
					valueVoid = false
				}
				continue

			case reflect.Map:
				registerGobTypeIfNeeded(value.Interface())
				var buf bytes.Buffer
				encoder := gob.NewEncoder(&buf)
				if err := encoder.Encode(value.Interface()); err != nil {
					return nil, fmt.Errorf("could not GOB-encode map: %w", err)
				}
				kvPair.BytesVal = buf.Bytes()
				valueVoid = false
				continue

			case reflect.Ptr:

				if value.IsNil() {
					// Ignore nil pointers
					continue
				}
				registerGobTypeIfNeeded(value.Interface())
				var buf bytes.Buffer
				encoder := gob.NewEncoder(&buf)
				if err := encoder.Encode(value.Interface()); err != nil {
					return nil, fmt.Errorf("could not GOB-encode pointer value: %w", err)
				}
				kvPair.BytesVal = buf.Bytes()
				valueVoid = false
				continue

			// 🕒 Special case for time.Time → store as int64 (Unix timestamp)
			case reflect.Struct:
				if value.Type() == reflect.TypeOf(time.Time{}) {
					timeValue := value.Interface().(time.Time)
					if !timeValue.IsZero() {
						intVal := timeValue.UTC().Unix()
						kvPair.Int64Val = &intVal
						valueVoid = false
						continue
					}
				}

			// ❌ Any other unsupported type is rejected explicitly
			default:
				return nil, errors.New(fmt.Sprintf("unsupported value type: %s", value.Kind().String()))
			}
		}

		// Process the `expireAt` field (tagged with `hydraide:"expireAt"`).
		// This defines the logical expiration time of the Treasure.
		// Once the given timestamp is reached, HydrAIDE will treat the record as expired.
		// - Must be of type `time.Time`
		// - Must not be the zero time
		// - Automatically converted to a `timestamppb.Timestamp` for protobuf
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagExpireAt {
			value := v.Field(i)
			if value.Kind() != reflect.Struct || value.Type() != reflect.TypeOf(time.Time{}) {
				return nil, errors.New("expireAt field must be a time.Time")
			}
			expireAt := value.Interface().(time.Time).UTC()
			if expireAt.IsZero() {
				return nil, errors.New("expireAt field must be a non-zero time.Time")
			}
			kvPair.ExpiredAt = timestamppb.New(expireAt)
			valueVoid = false
			continue
		}

		// Process the `createdBy` field (tagged with `hydraide:"createdBy"`).
		// Optional metadata indicating who or what created the Treasure.
		// - Must be of type `string`
		// - Empty values are ignored
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagCreatedBy {
			value := v.Field(i)
			if value.Kind() != reflect.String {
				return nil, errors.New("createdBy field must be a string")
			}
			if value.String() != "" {
				createdBy := value.String()
				kvPair.CreatedBy = &createdBy
				valueVoid = false
			}
			continue
		}

		// Process the `createdAt` field (tagged with `hydraide:"createdAt"`).
		// Optional metadata representing when the Treasure was created.
		// - Must be of type `time.Time`
		// - Must not be zero
		// - Converted to protobuf-compatible timestamp
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagCreatedAt {
			value := v.Field(i)
			if value.Kind() != reflect.Struct || value.Type() != reflect.TypeOf(time.Time{}) {
				return nil, errors.New("createdAt field must be a time.Time")
			}
			createdAt := value.Interface().(time.Time).UTC()
			if createdAt.IsZero() {
				return nil, errors.New("createdAt field must be a non-zero time.Time")
			}
			kvPair.CreatedAt = timestamppb.New(createdAt)
			valueVoid = false
			continue
		}

		// Process the `updatedBy` field (tagged with `hydraide:"updatedBy"`).
		// Optional metadata indicating who or what last updated the Treasure.
		// - Must be of type `string`
		// - Ignored if empty
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagUpdatedBy {
			value := v.Field(i)
			if value.Kind() != reflect.String {
				return nil, errors.New("updatedBy field must be a string")
			}
			if value.String() != "" {
				updatedBy := value.String()
				kvPair.UpdatedBy = &updatedBy
				valueVoid = false
			}
			continue
		}

		// Process the `updatedAt` field (tagged with `hydraide:"updatedAt"`).
		// Optional metadata representing the last modification time of the Treasure.
		// - Must be of type `time.Time`
		// - Must be non-zero
		// - Automatically converted to a `timestamppb.Timestamp` for protobuf transmission
		if key, ok := field.Tag.Lookup(tagHydrAIDE); ok && key == tagUpdatedAt {
			value := v.Field(i)
			if value.Kind() != reflect.Struct || value.Type() != reflect.TypeOf(time.Time{}) {
				return nil, errors.New("updatedAt field must be a time.Time")
			}
			updatedAt := value.Interface().(time.Time).UTC()
			if updatedAt.IsZero() {
				return nil, errors.New("updatedAt field must be a non-zero time.Time")
			}
			kvPair.UpdatedAt = timestamppb.New(updatedAt)
			valueVoid = false
			continue
		}

	}

	// Final validation: the key must be present and non-empty.
	// This is a hard requirement — all Treasures in HydrAIDE must have a key.
	if kvPair.Key == "" {
		return nil, errors.New("key field not found")
	}

	// If no value was set during processing, mark the KeyValuePair as void.
	// This tells HydrAIDE that the record has no explicit value (e.g. it's a flag, or purely metadata).
	if valueVoid {
		kvPair.VoidVal = &valueVoid
	}

	// Return the fully constructed KeyValuePair for insertion into the system.
	return kvPair, nil

}

// convertProtoTreasureToModel maps a hydraidepbgo.Treasure protobuf object back into a Go struct.
//
// The target model must be a pointer to a struct. Fields are matched using `hydraide` struct tags:
// - `key`: assigns Treasure.Key to the struct's key field.
// - `value`: maps the appropriate typed value from Treasure into the struct's value field.
// - `expireAt`, `createdBy`, `createdAt`, `updatedBy`, `updatedAt`: optional metadata fields.
//
// Supported value conversions include:
// - Primitive types: string, bool, intX, uintX, floatX
// - time.Time (from int64 UNIX timestamp)
// - []byte (raw bytes)
// - All other slices, maps, and pointers (GOB-encoded in BytesVal)
//
// If the field type does not match the Treasure value type, it is silently skipped.
// If decoding fails (e.g. from GOB), an error is returned.
func convertProtoTreasureToModel(treasure *hydraidepbgo.Treasure, model any) error {

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("input must be a pointer to a struct")
	}

	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagKey {
			v.Elem().Field(i).SetString(treasure.GetKey())
			continue
		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagValue {

			if treasure.StringVal != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.String:
					v.Elem().Field(i).SetString(treasure.GetStringVal())
					continue
				default:
					continue
				}
			}

			if treasure.Uint8Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Uint8:
					v.Elem().Field(i).SetUint(uint64(treasure.GetUint8Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Uint16Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Uint16:
					v.Elem().Field(i).SetUint(uint64(treasure.GetUint16Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Uint32Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Uint32:
					v.Elem().Field(i).SetUint(uint64(treasure.GetUint32Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Uint64Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Uint64:
					v.Elem().Field(i).SetUint(treasure.GetUint64Val())
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Int8Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Int8:
					v.Elem().Field(i).SetInt(int64(treasure.GetInt8Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Int16Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Int16:
					v.Elem().Field(i).SetInt(int64(treasure.GetInt16Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Int32Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Int32:
					v.Elem().Field(i).SetInt(int64(treasure.GetInt32Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Int64Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Int64:

					v.Elem().Field(i).SetInt(treasure.GetInt64Val())
					continue

				case reflect.Struct:

					// ha time.Time típusú mezőről van szó
					if v.Elem().Field(i).Type() == reflect.TypeOf(time.Time{}) {
						// konvertáljuk vissza time.Time-ra az int64 UNIX timestampet
						timestamp := time.Unix(treasure.GetInt64Val(), 0).UTC()
						v.Elem().Field(i).Set(reflect.ValueOf(timestamp))
					}
					continue

				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Float32Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Float32:
					v.Elem().Field(i).SetFloat(float64(treasure.GetFloat32Val()))
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.Float64Val != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Float64:
					v.Elem().Field(i).SetFloat(treasure.GetFloat64Val())
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.BoolVal != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Bool:
					v.Elem().Field(i).SetBool(treasure.GetBoolVal() == hydraidepbgo.Boolean_TRUE)
					continue
				default:
					// skip the field because the value type is not the same as the model field type
					continue
				}
			}

			if treasure.BytesVal != nil {
				switch v.Elem().Field(i).Kind() {
				case reflect.Slice:
					if v.Elem().Field(i).Type().Elem().Kind() == reflect.Uint8 {
						v.Elem().Field(i).SetBytes(treasure.GetBytesVal())
					} else {

						decoder := gob.NewDecoder(bytes.NewReader(treasure.GetBytesVal()))
						decoded := reflect.New(v.Elem().Field(i).Type()).Interface()

						if err := decoder.Decode(decoded); err != nil {
							return fmt.Errorf("failed to decode gob into slice field %s: %w", t.Name(), err)
						}

						v.Elem().Field(i).Set(reflect.ValueOf(decoded).Elem())
					}

				case reflect.Map, reflect.Ptr:

					decoder := gob.NewDecoder(bytes.NewReader(treasure.GetBytesVal()))
					decoded := reflect.New(v.Elem().Field(i).Type()).Interface()

					if err := decoder.Decode(decoded); err != nil {
						return fmt.Errorf("failed to decode gob into map/ptr field %s: %w", t.Name(), err)
					}

					v.Elem().Field(i).Set(reflect.ValueOf(decoded).Elem())

				default:
					continue
				}
			}

			continue

		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagExpireAt {
			if treasure.ExpiredAt != nil {
				v.Elem().Field(i).Set(reflect.ValueOf(treasure.ExpiredAt.AsTime()))
			}
			continue
		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagCreatedBy {
			if treasure.CreatedBy != nil {
				v.Elem().Field(i).SetString(*treasure.CreatedBy)
			}
			continue
		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagCreatedAt {
			if treasure.CreatedAt != nil {
				v.Elem().Field(i).Set(reflect.ValueOf(treasure.CreatedAt.AsTime()))
			}
			continue
		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagUpdatedBy {
			if treasure.UpdatedBy != nil {
				v.Elem().Field(i).SetString(*treasure.UpdatedBy)
			}
			continue
		}

		if key, ok := t.Field(i).Tag.Lookup(tagHydrAIDE); ok && key == tagUpdatedAt {
			if treasure.UpdatedAt != nil {
				v.Elem().Field(i).Set(reflect.ValueOf(treasure.UpdatedAt.AsTime()))
			}
			continue
		}

	}

	return nil

}

var (
	// registeredTypes keeps track of all types registered with gob to prevent duplicate registrations.
	registeredTypes = make(map[reflect.Type]struct{})
	mutex           sync.Mutex
)

// registerGobTypeIfNeeded safely registers a type with the gob encoder,
// making sure the same type is not registered multiple times.
//
// This function handles both pointer and base types by recursively registering
// the underlying struct when a pointer is provided.
//
// If the type has already been registered, the function does nothing.
// If gob.Register panics (e.g. due to naming conflicts), it recovers gracefully.
func registerGobTypeIfNeeded(val interface{}) {
	t := reflect.TypeOf(val)

	// If it's a pointer, first register the underlying base type (e.g., *StructX -> StructX)
	if t.Kind() == reflect.Ptr {
		registerGobTypeIfNeeded(reflect.Zero(t.Elem()).Interface())
		// Avoid registering the pointer type separately to reduce conflict risk.
		// Remove this return if you want to explicitly register both base and pointer types.
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// Only register if not already registered
	if _, ok := registeredTypes[t]; !ok {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered in registerGobTypeIfNeeded: %v\n", r)
			}
		}()

		gob.Register(val)
		registeredTypes[t] = struct{}{}
	}
}
