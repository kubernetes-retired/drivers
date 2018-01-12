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

package main

import (
	"fmt"
	"os"

	"github.com/kubernetes-csi/drivers/pkg/glusterfs"
	"github.com/spf13/cobra"
)

var (
	endpoint string
	mode     string
	nodeID   string
)

func main() {
	cmd := &cobra.Command{
		Use:   "GlusterFS",
		Short: "GlusterFS CSI driver",
		Run: func(cmd *cobra.Command, args []string) {
			handle()
		},
	}

	cmd.PersistentFlags().StringVar(&mode, "mode", "", "operating mode: 'node' or 'controller'")
	cmd.MarkPersistentFlagRequired("mode")

	cmd.PersistentFlags().StringVar(&nodeID, "nodeid", "", "node id")
	cmd.MarkPersistentFlagRequired("nodeid")

	cmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "CSI endpoint")
	cmd.MarkPersistentFlagRequired("endpoint")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func handle() {
	d := glusterfs.NewDriver(nodeID, endpoint)
	d.Run()
}
