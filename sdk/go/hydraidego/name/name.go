// Package name
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
//
// Package name provides a deterministic naming convention for HydrAIDE,
// helping developers structure and resolve data into a three-level hierarchy:
// Sanctuary → Realm → Swamp.
//
// This structure enables:
// - O(1) access to data based on name
// - Stateless client-side routing in multi-server setups
// - Predictable folder and server mapping without orchestrators
//
// Each Name is constructed step-by-step (Sanctuary → Realm → Swamp),
// and the full path can be retrieved using Get(). Additionally,
// GetServerNumber(allServers) maps the current name to a consistent
// server index (1-based), using a fast and collision-resistant hash.
//
// Example usage:
//
//	name := New().Sanctuary("users").Realm("profiles").Swamp("alice123")
//	fmt.Println(name.Get()) // "users/profiles/alice123"
//	fmt.Println(name.GetServerNumber(1000)) // e.g. 774
//
// Use Load(path) to reconstruct a Name from an existing path string.
//
// This package is used across HydrAIDE SDKs to:
// - Determine data placement
// - Support distributed architectures
// - Enforce clean, intention-driven naming
// ----------------------------------------
// 📘 HydrAIDE Go SDK
// Full SDK documentation:
// https://github.com/hydraide/hydraide/blob/main/docs/sdk/go/README.md
// ----------------------------------------
package name

import (
	"github.com/cespare/xxhash/v2"
	"strings"
	"sync"
)

// Name defines a structured identifier used in HydrAIDE to deterministically
// map data into a distributed, folder-based architecture.
//
// Each Name represents a three-level hierarchy:
//
//	Sanctuary → Realm → Swamp
//
// This structure is essential for:
// - Organizing Swamps into logical domains
// - Generating predictable folder paths
// - Assigning each Swamp to a specific server without coordination
//
// The interface supports fluent chaining:
//
//	name := New().
//	    Sanctuary("users").
//	    Realm("profiles").
//	    Swamp("alice123")
//
//	name.Get()                 // "users/profiles/alice123"
//	name.GetServerNumber(100) // e.g. 42
//
// Usage of GetServerNumber ensures even distribution of data across N servers,
// enabling stateless multi-node architectures without external orchestrators.
//
// See also: Load(path string) to reconstruct a Name from a path.
type Name interface {
	Sanctuary(sanctuaryID string) Name
	Realm(realmName string) Name
	Swamp(swampName string) Name
	Get() string
	GetServerNumber(allServers int) uint16
}

type name struct {
	Path           string
	SanctuaryID    string
	RealmName      string
	SwampName      string
	ServerNumber   uint16
	hashPathMu     sync.Mutex
	folderNumberMu sync.Mutex
}

// New creates a new empty Name instance.
// Use this as the starting point for building hierarchical names
// by chaining Sanctuary(), Realm(), and Swamp().
func New() Name {
	return &name{}
}

// Sanctuary sets the top-level domain of the Name.
// Typically used to group major logical areas (e.g. "users", "products").
func (n *name) Sanctuary(sanctuaryID string) Name {
	return &name{
		SanctuaryID: sanctuaryID,
		Path:        sanctuaryID,
	}
}

// Realm sets the second-level scope under the Sanctuary.
// Often used to further categorize Swamps (e.g. "profiles", "settings").
func (n *name) Realm(realmName string) Name {
	return &name{
		SanctuaryID: n.SanctuaryID,
		RealmName:   realmName,
		Path:        n.Path + "/" + realmName,
	}
}

// Swamp sets the final segment of the Name — the Swamp itself.
// This represents the concrete storage unit where Treasures are kept.
// The full path becomes: sanctuary/realm/swamp.
func (n *name) Swamp(swampName string) Name {
	return &name{
		SanctuaryID: n.SanctuaryID,
		RealmName:   n.RealmName,
		SwampName:   swampName,
		Path:        n.Path + "/" + swampName,
	}
}

// Get returns the full hierarchical path of the Name in the format:
//
//	"sanctuary/realm/swamp"
//
// 🔒 Internal use only: This method is intended for SDK-level logic,
// such as logging, folder path generation, or internal diagnostics.
// SDK users should never need to call this directly.
func (n *name) Get() string {
	return n.Path
}

// GetServerNumber returns the 1-based index of the server responsible for this Name.
// It uses a fast, consistent xxhash hash over the combined Sanctuary, Realm, and Swamp
// to deterministically assign the Name to one of `allServers` available slots.
//
// 🔒 Internal use only: This function is used by the SDK to route
// the Name to the correct Hydra client instance in a distributed setup.
// It should not be called directly by application developers.
//
// Example (inside SDK logic):
//
//	client := router.Route(name.GetServerNumber(1000))
func (n *name) GetServerNumber(allServers int) uint16 {

	n.folderNumberMu.Lock()
	defer n.folderNumberMu.Unlock()

	if n.ServerNumber != 0 {
		return n.ServerNumber
	}

	hash := xxhash.Sum64([]byte(n.SanctuaryID + n.RealmName + n.SwampName))

	n.ServerNumber = uint16(hash%uint64(allServers)) + 1

	return n.ServerNumber

}

// Load reconstructs a Name from a given path string in the format:
//
//	"sanctuary/realm/swamp"
//
// It parses the path segments and returns a Name instance with all fields set.
//
// 🔒 Internal use only: This function is intended for SDK-level logic,
// such as reconstructing a Name from persisted references, file paths, or routing metadata.
// It should not be called by application developers directly.
func Load(path string) Name {
	splitPath := strings.Split(path, "/")
	sanctuaryID := splitPath[0]
	realmName := splitPath[1]
	swampName := splitPath[2]

	return &name{
		Path:        sanctuaryID + "/" + realmName + "/" + swampName,
		SanctuaryID: sanctuaryID,
		RealmName:   realmName,
		SwampName:   swampName,
	}
}
