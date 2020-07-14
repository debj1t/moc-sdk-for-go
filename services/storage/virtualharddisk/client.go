// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache v2.0 License.

package virtualharddisk

import (
	"context"
	"github.com/microsoft/moc-proto/pkg/auth"
	"github.com/microsoft/moc-proto/pkg/errors"
	"github.com/microsoft/moc-sdk-for-go/services/storage"
)

// Service interface
type Service interface {
	Get(context.Context, string, string, string) (*[]storage.VirtualHardDisk, error)
	CreateOrUpdate(context.Context, string, string, string, *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error)
	Delete(context.Context, string, string, string) error
}

// Client structure
type VirtualHardDiskClient struct {
	storage.BaseClient
	internal Service
}

// NewClient method returns new client
func NewVirtualHardDiskClient(cloudFQDN string, authorizer auth.Authorizer) (*VirtualHardDiskClient, error) {
	c, err := newVirtualHardDiskClient(cloudFQDN, authorizer)
	if err != nil {
		return nil, err
	}

	return &VirtualHardDiskClient{internal: c}, nil
}

// Get methods invokes the client Get method
func (c *VirtualHardDiskClient) Get(ctx context.Context, group, container, name string) (*[]storage.VirtualHardDisk, error) {
	return c.internal.Get(ctx, group, container, name)
}

// CreateOrUpdate methods invokes create or update on the client
func (c *VirtualHardDiskClient) CreateOrUpdate(ctx context.Context, group, container, name string, storage *storage.VirtualHardDisk) (*storage.VirtualHardDisk, error) {
	return c.internal.CreateOrUpdate(ctx, group, container, name, storage)
}

// Delete methods invokes delete of the storage resource
func (c *VirtualHardDiskClient) Delete(ctx context.Context, group, container, name string) error {
	return c.internal.Delete(ctx, group, container, name)
}

// Resize methods invokes delete of the storage resource
func (c *VirtualHardDiskClient) Resize(ctx context.Context, group, container, name string, newSize int64) error {
	vhds, err := c.Get(ctx, group, container, name)
	if err != nil {
		return err
	}

	if len(*vhds) == 0 {
		return errors.Wrapf(errors.NotFound, "%s", name)
	}

	vhd := (*vhds)[0]
	vhd.DiskSizeBytes = &newSize

	_, err = c.CreateOrUpdate(ctx, group, container, name, &vhd)

	return err
}
