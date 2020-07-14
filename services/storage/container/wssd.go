// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package container

import (
	"context"
	"fmt"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"
	wssdcloudstorage "github.com/microsoft/moc-proto/rpc/cloudagent/storage"
	wssdcloudcommon "github.com/microsoft/moc-proto/rpc/common"
	wssdcloudclient "github.com/microsoft/moc-sdk-for-go/pkg/client"
	"github.com/microsoft/moc-sdk-for-go/services/storage"
)

type client struct {
	wssdcloudstorage.ContainerAgentClient
}

// newClient - creates a client session with the backend wssdcloud agent
func newContainerClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdcloudclient.GetStorageContainerClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, location, name string) (*[]storage.Container, error) {
	request, err := getContainerRequest(wssdcloudcommon.Operation_GET, location, name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	return getContainersFromResponse(response, location), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, location, name string, container *storage.Container) (*storage.Container, error) {
	request, err := getContainerRequest(wssdcloudcommon.Operation_POST, location, name, container)
	if err != nil {
		return nil, err
	}
	response, err := c.ContainerAgentClient.Invoke(ctx, request)
	if err != nil {
		return nil, err
	}
	containers := getContainersFromResponse(response, location)

	if len(*containers) == 0 {
		return nil, fmt.Errorf("[Container][Create] Unexpected error: Creating a storage interface returned no result")
	}

	return &((*containers)[0]), nil
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, location, name string) error {
	container, err := c.Get(ctx, location, name)
	if err != nil {
		return err
	}
	if len(*container) == 0 {
		return fmt.Errorf("Virtual Network [%s] not found", name)
	}

	request, err := getContainerRequest(wssdcloudcommon.Operation_DELETE, location, name, &(*container)[0])
	if err != nil {
		return err
	}
	_, err = c.ContainerAgentClient.Invoke(ctx, request)

	return err

}

func getContainerRequest(opType wssdcloudcommon.Operation, location, name string, storage *storage.Container) (*wssdcloudstorage.ContainerRequest, error) {
	request := &wssdcloudstorage.ContainerRequest{
		OperationType: opType,
		Containers:    []*wssdcloudstorage.Container{},
	}

	var err error

	wssdcontainer := &wssdcloudstorage.Container{
		Name:         name,
		LocationName: location,
	}

	if len(location) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Location not specified")
	}

	if storage != nil {
		wssdcontainer, err = getWssdContainer(storage, location)
		if err != nil {
			return nil, err
		}
	}
	request.Containers = append(request.Containers, wssdcontainer)

	return request, nil
}

func getContainersFromResponse(response *wssdcloudstorage.ContainerResponse, location string) *[]storage.Container {
	virtualHardDisks := []storage.Container{}
	for _, container := range response.GetContainers() {
		virtualHardDisks = append(virtualHardDisks, *(getContainer(container, location)))
	}

	return &virtualHardDisks
}
