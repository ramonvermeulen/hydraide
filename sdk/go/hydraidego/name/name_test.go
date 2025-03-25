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
package name

import (
	"crypto/rand"
	"math/big"
	"path/filepath"
	"testing"
)

const (
	testDataCount = 100000
	allServers    = 100
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

func TestGetServerNumberDistribution(t *testing.T) {
	serverCounts := make(map[uint16]int)

	for i := 0; i < testDataCount; i++ {
		// Generate random Sanctuary, Realm, and Swamp names
		sanctuary := randomString(8)
		realm := randomString(8)
		swamp := randomString(8)

		// Construct a new Name
		n := New().Sanctuary(sanctuary).Realm(realm).Swamp(swamp)

		// Determine the server number
		server := n.GetServerNumber(allServers)

		// Increment the counter for the corresponding server
		serverCounts[server]++
	}

	// Log the number of entries per server
	for server, count := range serverCounts {
		t.Logf("Server %d: %d entries", server, count)
	}

	// Check if the distribution is reasonably even
	expectedPerServer := testDataCount / allServers
	threshold := expectedPerServer / 10 // Allow 10% deviation

	for server, count := range serverCounts {
		if count < expectedPerServer-threshold || count > expectedPerServer+threshold {
			t.Errorf("Server %d received too few or too many entries: %d", server, count)
		}
	}
}

// BenchmarkName_Load measures the performance of reconstructing a Name
// from a canonical path string (e.g. "users/johndoe/info").
// goos: linux
// goarch: amd64
// pkg: github.com/trendizz/hydra-spine/hydra/name
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkName_LoadFromCanonicalForm
// BenchmarkName_Load-32    	 3671059	       320.5 ns/op
func BenchmarkName_Load(b *testing.B) {
	canonicalForm := filepath.Join("users", "johndoe", "info")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Load(canonicalForm)
	}
}

// BenchmarkName_Get measures the performance of the Get() method,
// which returns the canonical path representation of a Name instance
// in the format: "sanctuary/realm/swamp".
//
// This method is lightweight and frequently used internally by the SDK
// for tasks like logging, routing, and debugging.
//
// Example output: "users/johndoe/info"
//
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraide/sdk/go/hydraidego/name
// cpu: AMD Ryzen Threadripper 2950X 16-Core Processor
// BenchmarkName_Get
// BenchmarkName_Get-32    	1000000000	         0.5228 ns/op
// PASS
func BenchmarkName_Get(b *testing.B) {
	nameObj := New().Sanctuary("users").Realm("johndoe").Swamp("info")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nameObj.Get()
	}
}

// BenchmarkName_Add measures the performance of constructing a full Name
// by chaining Sanctuary(), Realm(), and Swamp() calls.
//
// This benchmark reflects how efficiently the SDK can create hierarchical
// name representations, which are the foundation of most operations in HydrAIDE.
//
// The resulting structure is used internally for data addressing and routing.
//
// Example: "users/johndoe/info"
//
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraide/sdk/go/hydraidego/name
// cpu: AMD Ryzen Threadripper 2950X 16-Core Processor
// BenchmarkName_Add
// BenchmarkName_Add-32    	25031070	        41.09 ns/op
// PASS
func BenchmarkName_Add(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New().Sanctuary("users").Realm("johndoe").Swamp("info")
	}
}

// BenchmarkGetServerNumber measures the performance of computing the
// server assignment for a given Name using the GetServerNumber() method.
//
// This method uses xxhash to deterministically assign a Name to one of N servers,
// enabling stateless, distributed routing logic inside the SDK.
//
// The result is a 1-based server number (e.g. 1–1000), based on the full path:
// "BenchmarkSanctuary/BenchmarkRealm/BenchmarkSwamp".
//
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraide/sdk/go/hydraidego/name
// cpu: AMD Ryzen Threadripper 2950X 16-Core Processor
// BenchmarkGetServerNumber
// BenchmarkGetServerNumber-32    	76946376	        15.19 ns/op
// PASS
func BenchmarkGetServerNumber(b *testing.B) {
	sanctuary := "BenchmarkSanctuary"
	realm := "BenchmarkRealm"
	swamp := "BenchmarkSwamp"
	n := New().Sanctuary(sanctuary).Realm(realm).Swamp(swamp)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = n.GetServerNumber(allServers)
	}
}
