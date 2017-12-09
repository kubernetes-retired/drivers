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
	"github.com/golang/glog"
	gcli "github.com/heketi/heketi/client/api/go-client"
	gapi "github.com/heketi/heketi/pkg/glusterfs/api"
)

func CreateGlusterFSVolume() (r *v1.GlusterfsVolumeSource, size int, err error) {
	var clusterIDs []string
	if p.url == "" {
		glog.Errorf("REST server endpoint is empty")
		return nil, 0, fmt.Errorf("failed to create glusterfs REST client, REST URL is empty")
	}
	cli := gcli.NewClient(p.url, p.user, p.secretValue)
	if cli == nil {
		glog.Errorf("failed to create glusterfs rest client")
		return nil, 0, fmt.Errorf("failed to create glusterfs REST client, REST server authentication failed")
	}
	if p.provisionerConfig.clusterID != "" {
		clusterIDs = dstrings.Split(p.clusterID, ",")
		glog.V(4).Infof("provided clusterIDs: %v", clusterIDs)
	}
	gid64 := int64(gid)
	volumeReq := &gapi.VolumeCreateRequest{
		Size:                 sz, 
		Clusters:             clusterIDs, 
		Gid:                  gid64, 
		Durability:           p.volumeType, 
		GlusterVolumeOptions: p.volumeOptions
	}
	volume, err := cli.VolumeCreate(volumeReq)
	if err != nil {
		glog.Errorf("Error creating volume: %v", err)
		return nil, err
	}
	glog.V(1).Infof("volume with size: %d and name: %s created", volume.Size, volume.Name)
	dynamicHostIps, err := getClusterNodes(cli, volume.Cluster)
	if err != nil {
		glog.Errorf("error [%v] when getting cluster nodes for volume %s", err, volume)
		return nil, 0, fmt.Errorf("error [%v] when getting cluster nodes for volume %s", err, volume)
	}

	// The 'endpointname' is created in form of 'glusterfs-dynamic-<claimname>'.
	// createEndpointService() checks for this 'endpoint' existence in PVC's namespace and
	// If not found, it create an endpoint and svc using the IPs we dynamically picked at time
	// of volume creation.
	epServiceName := dynamicEpSvcPrefix + p.options.PVC.Name
	epNamespace := p.options.PVC.Namespace
	endpoint, service, err := p.createEndpointService(epNamespace, epServiceName, dynamicHostIps, p.options.PVC.Name)
	if err != nil {
		glog.Errorf("failed to create endpoint/service: %v", err)
		deleteErr := cli.VolumeDelete(volume.Id)
		if deleteErr != nil {
			glog.Errorf("error when deleting the volume :%v , manual deletion required", deleteErr)
		}
		return nil, 0, fmt.Errorf("failed to create endpoint/service %v", err)
	}
	glog.V(3).Infof("dynamic ep %v and svc : %v ", endpoint, service)
	return &v1.GlusterfsVolumeSource{
		EndpointsName: endpoint.Name,
		Path:          volume.Name,
		ReadOnly:      false,
	}, sz, nil
}
