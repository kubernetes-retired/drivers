/*
Copyright 2017 The Kubernetes Authors.

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

package image

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"

	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

type image struct {
	driver *csicommon.CSIDriver

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

type imageVolume struct {
	VolName string `json:"volName"`
	VolID   string `json:"volID"`
	VolSize int64  `json:"volSize"`
	VolPath string `json:"volPath"`
}

var imageVolumes map[string]imageVolume

var (
	imageDriver *image
	vendorVersion  = "0.2.0"
)

func init() {
	imageVolumes = map[string]imageVolume{}
}

func GetImageDriver() *image {
	return &image{}
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
	}
}

func NewNodeServer(d *csicommon.CSIDriver) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func (i *image) Run(driverName, nodeID, endpoint string) {
	glog.Infof("Driver: %v ", driverName)

	// Initialize default library driver
	i.driver = csicommon.NewCSIDriver(driverName, vendorVersion, nodeID)
	if i.driver == nil {
		glog.Fatalln("Failed to initialize CSI Driver.")
	}
	i.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY})

	// Create GRPC servers
	i.ids = NewIdentityServer(i.driver)
	i.ns = NewNodeServer(i.driver)
	i.cs = NewControllerServer(i.driver)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(endpoint, i.ids, i.cs, i.ns)
	s.Wait()
}
