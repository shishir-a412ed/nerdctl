/*
   Copyright The containerd Authors.

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
	"testing"

	"github.com/containerd/nerdctl/pkg/testutil"
)

func TestCompletion(t *testing.T) {
	testutil.DockerIncompatible(t)
	base := testutil.NewBase(t)
	const gbc = "--generate-bash-completion"
	// cmd is executed with base.Args={"--namespace=nerdctl-test"}
	base.Cmd("--cgroup-manager", gbc).AssertOutContains("cgroupfs\n")
	base.Cmd("--snapshotter", gbc).AssertOutContains("native\n")
	base.Cmd(gbc).AssertOutContains("run\n")
	base.Cmd(gbc, "--snapshotter", gbc).AssertOutContains("native\n")
	base.Cmd("run", "-", gbc).AssertOutContains("--network\n")
	base.Cmd("run", "-", gbc).AssertNoOut("--namespace\n")      // --namespace is a global flag, not "run" flag
	base.Cmd("run", "-", gbc).AssertNoOut("--cgroup-manager\n") // --cgroup-manager is a global flag, not "run" flag
	base.Cmd("run", "-n", gbc).AssertOutContains("--network\n")
	base.Cmd("run", "-n", gbc).AssertNoOut("--namespace\n") // --namespace is a global flag, not "run" flag
	base.Cmd("run", "--ne", gbc).AssertOutContains("--network\n")
	base.Cmd("run", "--net", gbc).AssertOutContains("bridge\n")
	base.Cmd("run", "--net", gbc).AssertOutContains("host\n")
	base.Cmd("run", "-it", "--net", gbc).AssertOutContains("bridge\n")
	base.Cmd("run", "-it", "--rm", "--net", gbc).AssertOutContains("bridge\n")
	base.Cmd("run", "--restart", gbc).AssertOutContains("always\n")
	base.Cmd("network", "inspect", gbc).AssertOutContains("bridge\n")
	base.Cmd("network", "rm", gbc).AssertNoOut("bridge\n") // bridge is unremovable
	base.Cmd("network", "rm", gbc).AssertNoOut("host\n")   // host is unremovable
	base.Cmd("run", "--cap-add", gbc).AssertOutContains("sys_admin\n")
	base.Cmd("run", "--cap-add", gbc).AssertNoOut("CAP_SYS_ADMIN\n") // invalid form

	// Tests with an image
	base.Cmd("pull", testutil.AlpineImage).AssertOK()
	base.Cmd("run", "-i", gbc).AssertOutContains(testutil.AlpineImage)
	base.Cmd("run", "-it", gbc).AssertOutContains(testutil.AlpineImage)
	base.Cmd("run", "-it", "--rm", gbc).AssertOutContains(testutil.AlpineImage)

	// Tests with an network
	testNetworkName := "nerdctl-test-completion"
	defer base.Cmd("network", "rm", testNetworkName).Run()
	base.Cmd("network", "create", testNetworkName).AssertOK()
	base.Cmd("network", "rm", gbc).AssertOutContains(testNetworkName)
	base.Cmd("run", "--net", gbc).AssertOutContains(testNetworkName)

	// Tests with raw base (without Args={"--namespace=nerdctl-test"})
	rawBase := testutil.NewBase(t)
	rawBase.Args = nil // unset "--namespace=nerdctl-test"
	rawBase.Cmd("--cgroup-manager", gbc).AssertOutContains("cgroupfs\n")
	rawBase.Cmd(gbc).AssertOutContains("run\n")
	// mind {"--namespace=nerdctl-test"} vs {"--namespace", "nerdctl-test"}
	rawBase.Cmd("--namespace", testutil.Namespace, gbc).AssertOutContains("run\n")
	rawBase.Cmd("--namespace", testutil.Namespace, "run", "-i", gbc).AssertOutContains(testutil.AlpineImage)
}
