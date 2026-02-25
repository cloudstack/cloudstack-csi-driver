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

package cloud

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

func (c *client) GetSnapshotByID(ctx context.Context, snapshotID string) (*Snapshot, error) {
	logger := klog.FromContext(ctx)
	p := c.Snapshot.NewListSnapshotsParams()
	if snapshotID != "" {
		p.SetId(snapshotID)
	}
	logger.V(2).Info("CloudStack API call", "command", "ListSnapshots", "params", map[string]string{
		"id": snapshotID,
	})
	l, err := c.Snapshot.ListSnapshots(p)
	if err != nil {
		return nil, err
	}
	if l.Count == 0 {
		return nil, ErrNotFound
	}
	if l.Count > 1 {
		return nil, ErrTooManyResults
	}
	snapshot := l.Snapshots[0]
	s := Snapshot{
		ID:        snapshot.Id,
		Name:      snapshot.Name,
		DomainID:  snapshot.Domainid,
		ProjectID: snapshot.Projectid,
		ZoneID:    snapshot.Zoneid,
		VolumeID:  snapshot.Volumeid,
	}

	return &s, nil
}

func (c *client) CreateSnapshot(ctx context.Context, volumeID, name string) (*Snapshot, error) {
	logger := klog.FromContext(ctx)
	p := c.Snapshot.NewCreateSnapshotParams(volumeID)
	if name != "" {
		p.SetName(name)
	}
	logger.V(2).Info("CloudStack API call", "command", "CreateSnapshot", "params", map[string]string{
		"volumeid": volumeID,
		"name":     name,
	})

	snapshot, err := c.Snapshot.CreateSnapshot(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error %v", err)
	}

	snap := Snapshot{
		ID:        snapshot.Id,
		Name:      snapshot.Name,
		Size:      snapshot.Virtualsize,
		DomainID:  snapshot.Domainid,
		ProjectID: snapshot.Projectid,
		ZoneID:    snapshot.Zoneid,
		VolumeID:  snapshot.Volumeid,
		CreatedAt: snapshot.Created,
	}

	return &snap, nil
}

func (c *client) DeleteSnapshot(_ context.Context, snapshotID string) error {
	p := c.Snapshot.NewDeleteSnapshotParams(snapshotID)
	_, err := c.Snapshot.DeleteSnapshot(p)
	if err != nil && strings.Contains(err.Error(), "4350") {
		// CloudStack error InvalidParameterValueException
		return ErrNotFound
	}

	return err
}

func (c *client) GetSnapshotByName(ctx context.Context, name string) (*Snapshot, error) {
	logger := klog.FromContext(ctx)
	if name == "" {
		return nil, ErrNotFound
	}
	p := c.Snapshot.NewListSnapshotsParams()
	p.SetName(name)
	logger.V(2).Info("CloudStack API call", "command", "ListSnapshots", "params", map[string]string{
		"name": name,
	})
	l, err := c.Snapshot.ListSnapshots(p)
	if err != nil {
		return nil, err
	}
	if l.Count == 0 {
		return nil, ErrNotFound
	}
	if l.Count > 1 {
		return nil, ErrTooManyResults
	}
	snapshot := l.Snapshots[0]
	s := Snapshot{
		ID:        snapshot.Id,
		Name:      snapshot.Name,
		DomainID:  snapshot.Domainid,
		ProjectID: snapshot.Projectid,
		ZoneID:    snapshot.Zoneid,
		VolumeID:  snapshot.Volumeid,
		CreatedAt: snapshot.Created,
	}

	return &s, nil
}

func (c *client) ListSnapshots(ctx context.Context, volumeID, snapshotID string) ([]*Snapshot, error) {
	logger := klog.FromContext(ctx)
	p := c.Snapshot.NewListSnapshotsParams()
	if snapshotID != "" {
		p.SetId(snapshotID)
	}
	if volumeID != "" {
		p.SetVolumeid(volumeID)
	}
	logger.V(2).Info("CloudStack API call", "command", "ListSnapshots", "params", map[string]string{
		"id":       snapshotID,
		"volumeid": volumeID,
	})
	l, err := c.Snapshot.ListSnapshots(p)
	if err != nil {
		return nil, err
	}
	if l.Count == 0 {
		return []*Snapshot{}, nil
	}
	result := make([]*Snapshot, 0, l.Count)
	for _, snapshot := range l.Snapshots {
		s := &Snapshot{
			ID:        snapshot.Id,
			Name:      snapshot.Name,
			Size:      snapshot.Virtualsize,
			DomainID:  snapshot.Domainid,
			ProjectID: snapshot.Projectid,
			ZoneID:    snapshot.Zoneid,
			VolumeID:  snapshot.Volumeid,
			CreatedAt: snapshot.Created,
		}
		result = append(result, s)
	}

	return result, nil
}
