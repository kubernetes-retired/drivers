/*
Copyright 2018 The Kubernetes Authors.

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

package hostpath

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubernetes-csi/csi-test/pkg/sanity"
	"github.com/stretchr/testify/require"
)

func TestHostpathDriver(t *testing.T) {
	tmp, err := ioutil.TempDir("", "hostpath")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	driver := GetHostPathDriver()
	endpoint := "unix://" + tmp + "/hostpath-driver.sock"
	s, err := driver.Start("hostpath-driver", "test-node", endpoint)
	defer s.ForceStop()

	// Now call the test suite.
	config := sanity.Config{
		TargetPath:  tmp + "/target-path",
		StagingPath: tmp + "/staging-path",
		Address:     endpoint,
	}
	sanity.Test(t, &config)
}
