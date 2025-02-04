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

package genutil

import (
	"github.com/operator-framework/operator-registry/pkg/lib/bundle"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/graphitehealth/operator-sdk/internal/generate/collector"
)

// GetManifestObjects returns all objects to be written to a manifests directory from collector.Manifests.
func GetManifestObjects(c *collector.Manifests, extraSAs []string) (objs []client.Object) {
	// All CRDs passed in should be written.
	for i := range c.V1CustomResourceDefinitions {
		objs = append(objs, &c.V1CustomResourceDefinitions[i])
	}
	for i := range c.V1beta1CustomResourceDefinitions {
		objs = append(objs, &c.V1beta1CustomResourceDefinitions[i])
	}

	// All Services passed in should be written.
	for i := range c.Services {
		objs = append(objs, &c.Services[i])
	}

	// Add all other supported kinds
	for i := range c.Others {
		obj := &c.Others[i]
		if supported, _ := bundle.IsSupported(obj.GroupVersionKind().Kind); supported {
			objs = append(objs, obj)
		}
	}

	// RBAC objects (including ServiceAccounts) that are not a part of the CSV should be written.
	_, _, rbacObjs := c.SplitCSVPermissionsObjects(extraSAs)
	objs = append(objs, rbacObjs...)

	removeNamespace(objs)
	return objs
}

// removeNamespace removes the namespace field of resources intended to be inserted into
// an OLM manifests directory.
//
// This is required to pass OLM validations which require that namespaced resources do
// not include explicit namespace settings. OLM automatically installs namespaced
// resources in the same namespace that the operator is installed in, which is determined
// at runtime, not bundle/packagemanifests creation time.
func removeNamespace(objs []client.Object) {
	for _, obj := range objs {
		obj.SetNamespace("")
	}
}
