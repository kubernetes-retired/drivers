/*
Copyright 2017 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package glusterfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type driver struct {
	csiDriver *csicommon.CSIDriver
	endpoint  string

	ids *csicommon.DefaultIdentityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

const (
	driverName = "GlusterFS"
)

var (
	version = csi.Version{
		Minor: 1,
	}
)

func GetSupportedVersions() []*csi.Version {
	return []*csi.Version{&version}
}

func NewDriver(nodeID, endpoint string) *driver {
	glog.Infof("Driver: %v version: %v", driverName, csicommon.GetVersionString(&version))

	d := &driver{}

	d.endpoint = endpoint

	csiDriver := csicommon.NewCSIDriver(driverName, &version, GetSupportedVersions(), nodeID)
	csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME})
	csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME})
	csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER})

	d.csiDriver = csiDriver

	return d
}

func NewNodeServer(d *csicommon.csiDriver) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
	}
}

func (d *driver) Run() {
	if mode == "node" {
		d.ns = NewNodeServer(d.csiDriver)
		csicommon.RunNodePublishServer(endpoint, d.ns)
	} else if mode == "controller" {
		d.cs = NewControllerServer(d.csiDriver)
		csicommon.RunControllerPublishServer(endpoint, d.cs)
	} else {
		glog.Errorf("Unknown driver mode '%v'", mode)
	}
}
