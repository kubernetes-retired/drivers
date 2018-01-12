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

package glusterfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"k8s.io/kubernetes/pkg/volume"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	gigabyte = int64(1024 * 1024 * 1024)
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func GetVersionString(ver *csi.Version) string {
	return fmt.Sprintf("%d.%d.%d", ver.Major, ver.Minor, ver.Patch)
}

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// Volume Name
	volName := req.GetName()
	if len(volName) == 0 {
		volName = uuid.NewUUID().String()
	}

	// Volume Size - Default is 1 GiB
	volSizeBytes := gigabyte
	if req.GetCapacityRange() != nil {
		volSizeBytes = int64(req.GetCapacityRange().GetRequiredBytes())
	}
	volSizeGB := int(volume.RoundUpSize(volSizeBytes, gigabyte))

	// Volume Parameters
	volFoo := req.GetParameters()["foo"]
	if len(volFoo) == 0 {
		volFoo = "default"
	}

	// Volume Create
	_, err := CreateGlusterFSVolume()
	if err != nil {
		glog.Errorf("CreateVolume failed: %v", err)
		return nil, err
	}
	glog.V(1).Infof("Succesfully created volume '%v'", volName)

	resp := &csi.CreateVolumeResponse{
		VolumeInfo: &csi.VolumeInfo{
			Id: resID,
			Attributes: map[string]string{
				"foo": volFoo,
			},
		},
	}
	return resp, nil
}

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
}

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
}

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	for _, cap := range req.VolumeCapabilities {
		if cap.GetAccessMode().GetMode() != csi.VolumeCapability_AccessMode_MULTI_NODE_WRITER {
			return &csi.ValidateVolumeCapabilitiesResponse{false, ""}, nil
		}
	}
	return &csi.ValidateVolumeCapabilitiesResponse{true, ""}, nil
}
