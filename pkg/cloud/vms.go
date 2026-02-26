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

	"k8s.io/klog/v2"
)

func (c *client) GetVMByID(ctx context.Context, vmID string) (*VM, error) {
	logger := klog.FromContext(ctx)
	logger.V(2).Info("CloudStack API call", "command", "GetVirtualMachineByID", "params", map[string]string{
		"id": vmID,
	})

	vm, _, err := c.VirtualMachine.GetVirtualMachineByID(vmID)
	if err != nil {
		return nil, err
	}

	return &VM{
		ID:     vm.Id,
		ZoneID: vm.Zoneid,
	}, nil
}

func (c *client) getVMByName(ctx context.Context, name string) (*VM, error) {
	logger := klog.FromContext(ctx)
	logger.V(2).Info("CloudStack API call", "command", "GetVirtualMachineByName", "params", map[string]string{
		"name": name,
	})

	vm, _, err := c.VirtualMachine.GetVirtualMachineByName(name)
	if err != nil {
		return nil, err
	}

	return &VM{
		ID:     vm.Id,
		ZoneID: vm.Zoneid,
	}, nil
}
