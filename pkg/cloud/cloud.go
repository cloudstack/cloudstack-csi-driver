//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Package cloud contains CloudStack related
// functions.
package cloud

import (
	"context"
	"errors"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"k8s.io/klog/v2"
)

// Interface is the CloudStack client interface.

//nolint:interfacebloat
type Interface interface {
	GetNodeInfo(ctx context.Context, vmName string) (*VM, error)
	GetVMByID(ctx context.Context, vmID string) (*VM, error)

	ListZonesID(ctx context.Context) ([]string, error)

	GetVolumeByID(ctx context.Context, volumeID string) (*Volume, error)
	GetVolumeByName(ctx context.Context, name string) (*Volume, error)
	CreateVolume(ctx context.Context, diskOfferingID, zoneID, name string, sizeInGB int64) (string, error)
	DeleteVolume(ctx context.Context, id string) error
	AttachVolume(ctx context.Context, volumeID, vmID string) (string, error)
	DetachVolume(ctx context.Context, volumeID string) error
	ExpandVolume(ctx context.Context, volumeID string, newSizeInGB int64) error

	CreateVolumeFromSnapshot(ctx context.Context, zoneID, name, projectID, snapshotID string, sizeInGB int64) (*Volume, error)
	GetSnapshotByID(ctx context.Context, snapshotID string) (*Snapshot, error)
	GetSnapshotByName(ctx context.Context, name string) (*Snapshot, error)
	CreateSnapshot(ctx context.Context, volumeID, name string) (*Snapshot, error)
	DeleteSnapshot(ctx context.Context, snapshotID string) error
	ListSnapshots(ctx context.Context, volumeID, snapshotID string) ([]*Snapshot, error)
}

// Volume represents a CloudStack volume.
type Volume struct {
	ID   string
	Name string

	// Size in Bytes
	Size int64

	DiskOfferingID string
	DomainID       string
	ProjectID      string
	ZoneID         string

	VirtualMachineID string
	DeviceID         string
}

// Volume represents a CloudStack snapshot.
type Snapshot struct {
	ID   string
	Name string
	Size int64

	DomainID  string
	ProjectID string
	ZoneID    string

	VolumeID  string
	CreatedAt string
}

// VM represents a CloudStack Virtual Machine.
type VM struct {
	ID     string
	ZoneID string
}

// Specific errors.
var (
	ErrNotFound       = errors.New("not found")
	ErrTooManyResults = errors.New("too many results")
	ErrAlreadyExists  = errors.New("already exists")
)

// client is the implementation of Interface.
type client struct {
	*cloudstack.CloudStackClient
}

// New creates a new cloud connector, given its configuration.
func New(config *Config) Interface {
	csClient := cloudstack.NewAsyncClient(config.APIURL, config.APIKey, config.SecretKey, config.VerifySSL)

	// Set the project id to every request.
	// This is possible because we just could work in one project with the previous implementation.
	if config.ProjectID != "" {
		csClient.DefaultOptions(cloudstack.WithProject(config.ProjectID))
		klog.Background().V(2).Info("Set projectID to cloud connector", "projectID", config.ProjectID)
	}

	return &client{csClient}
}
