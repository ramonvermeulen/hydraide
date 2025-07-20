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
