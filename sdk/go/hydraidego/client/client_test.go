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
	"github.com/hydraide/hydraide/generated/hydraidepbgo"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
)

// fakeHydraideServiceClient is a dummy implementation for test use only.
type fakeHydraideServiceClient struct {
	hydraidepbgo.HydraideServiceClient
}

func (f *fakeHydraideServiceClient) Heartbeat(
	_ context.Context,
	_ *hydraidepbgo.HeartbeatRequest,
	_ ...grpc.CallOption,
) (*hydraidepbgo.HeartbeatResponse, error) {
	return &hydraidepbgo.HeartbeatResponse{Pong: "beat"}, nil
}

func TestClient_GetServiceClient(t *testing.T) {

	// Arrange
	mockClient := &fakeHydraideServiceClient{}
	c := &mockedClient{
		allFolders: 1000,
		serviceClients: map[uint64]hydraidepbgo.HydraideServiceClient{
			518: mockClient, // assume this is the computed folder
		},
	}

	swamp := name.New().Sanctuary("users").Realm("profiles").Swamp("john.doe")
	folder := swamp.GetIslandID(c.allFolders)

	// Act
	serviceClient := c.GetServiceClient(swamp)

	// Assert
	assert.NotNil(t, serviceClient)
	assert.Equal(t, mockClient, serviceClient)
	assert.Equal(t, uint64(518), folder, "Folder number mismatch – update map if needed")

}

// mockedClient implements client.Client but skips actual gRPC connection
type mockedClient struct {
	allFolders     uint64
	serviceClients map[uint64]hydraidepbgo.HydraideServiceClient
}

func (c *mockedClient) Connect(_ bool) error { return nil }
func (c *mockedClient) CloseConnection()     {}
func (c *mockedClient) GetServiceClient(swamp name.Name) hydraidepbgo.HydraideServiceClient {
	folder := swamp.GetIslandID(c.allFolders)
	return c.serviceClients[folder]
}
