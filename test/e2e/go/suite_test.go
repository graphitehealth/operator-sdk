// Copyright 2020 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2e_go_test

import (
	"fmt"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/graphitehealth/operator-sdk/internal/testutils"
)

// TestE2EGo ensures the Go projects built with the SDK tool by using its binary.
func TestE2EGo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Operator SDK E2E Go Suite testing in short mode")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2EGo Suite")
}

var (
	tc testutils.TestContext
)

// BeforeSuite run before any specs are run to perform the required actions for all e2e Go tests.
var _ = BeforeSuite(func() {
	var err error

	By("creating a new test context")
	tc, err = testutils.NewTestContext(testutils.BinaryName, "GO111MODULE=on")
	Expect(err).NotTo(HaveOccurred())

	tc.Domain = "example.com"
	tc.Group = "cache"
	tc.Version = "v1alpha1"
	tc.Kind = "Memcached"
	tc.Resources = "memcacheds"
	tc.ProjectName = "memcached-operator"
	tc.Kubectl.Namespace = fmt.Sprintf("%s-system", tc.ProjectName)
	tc.Kubectl.ServiceAccount = fmt.Sprintf("%s-controller-manager", tc.ProjectName)

	By("copying sample to a temporary e2e directory")
	Expect(exec.Command("cp", "-r", "../../../testdata/go/v4/memcached-operator", tc.Dir).Run()).To(Succeed())

	By("preparing the prerequisites on cluster")
	tc.InstallPrerequisites()

	By("by adding scorecard custom patch file")
	err = tc.AddScorecardCustomPatchFile()
	Expect(err).NotTo(HaveOccurred())

	By("using dev image for scorecard-test")
	err = tc.ReplaceScorecardImagesForDev()
	Expect(err).NotTo(HaveOccurred())

	By("building the project image")
	err = tc.Make("docker-build", "IMG="+tc.ImageName)
	Expect(err).NotTo(HaveOccurred())

	onKind, err := tc.IsRunningOnKind()
	Expect(err).NotTo(HaveOccurred())
	if onKind {
		By("loading the required images into Kind cluster")
		Expect(tc.LoadImageToKindCluster()).To(Succeed())
		Expect(tc.LoadImageToKindClusterWithName("quay.io/operator-framework/scorecard-test:dev")).To(Succeed())
		Expect(tc.LoadImageToKindClusterWithName("quay.io/operator-framework/custom-scorecard-tests:dev")).To(Succeed())
	}

	By("generating bundle")
	Expect(tc.GenerateBundle()).To(Succeed())

	By("installing cert manager bundle")
	Expect(tc.InstallCertManager()).To(Succeed())
})

// AfterSuite run after all the specs have run, regardless of whether any tests have failed to ensures that
// all be cleaned up
var _ = AfterSuite(func() {
	By("uninstall cert manager bundle")
	tc.UninstallCertManager()

	By("uninstalling prerequisites")
	tc.UninstallPrerequisites()

	By("destroying container image and work dir")
	tc.Destroy()
})
