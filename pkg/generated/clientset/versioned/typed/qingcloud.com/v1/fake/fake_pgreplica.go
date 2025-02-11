/*
Copyright 2020 - 2021 Qingcloud Data Solutions, Inc.
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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	qingcloudcomv1 "github.com/qingcloud/postgres-operator/pkg/apis/qingcloud.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePgreplicas implements PgreplicaInterface
type FakePgreplicas struct {
	Fake *FakeQingcloudV1
	ns   string
}

var pgreplicasResource = schema.GroupVersionResource{Group: "qingcloud.com", Version: "v1", Resource: "pgreplicas"}

var pgreplicasKind = schema.GroupVersionKind{Group: "qingcloud.com", Version: "v1", Kind: "Pgreplica"}

// Get takes name of the pgreplica, and returns the corresponding pgreplica object, and an error if there is any.
func (c *FakePgreplicas) Get(ctx context.Context, name string, options v1.GetOptions) (result *qingcloudcomv1.Pgreplica, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(pgreplicasResource, c.ns, name), &qingcloudcomv1.Pgreplica{})

	if obj == nil {
		return nil, err
	}
	return obj.(*qingcloudcomv1.Pgreplica), err
}

// List takes label and field selectors, and returns the list of Pgreplicas that match those selectors.
func (c *FakePgreplicas) List(ctx context.Context, opts v1.ListOptions) (result *qingcloudcomv1.PgreplicaList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(pgreplicasResource, pgreplicasKind, c.ns, opts), &qingcloudcomv1.PgreplicaList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &qingcloudcomv1.PgreplicaList{ListMeta: obj.(*qingcloud.comcomv1.PgreplicaList).ListMeta}
	for _, item := range obj.(*qingcloudcomv1.PgreplicaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested pgreplicas.
func (c *FakePgreplicas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(pgreplicasResource, c.ns, opts))

}

// Create takes the representation of a pgreplica and creates it.  Returns the server's representation of the pgreplica, and an error, if there is any.
func (c *FakePgreplicas) Create(ctx context.Context, pgreplica *qingcloudcomv1.Pgreplica, opts v1.CreateOptions) (result *qingcloud.comcomv1.Pgreplica, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(pgreplicasResource, c.ns, pgreplica), &qingcloudcomv1.Pgreplica{})

	if obj == nil {
		return nil, err
	}
	return obj.(*qingcloudcomv1.Pgreplica), err
}

// Update takes the representation of a pgreplica and updates it. Returns the server's representation of the pgreplica, and an error, if there is any.
func (c *FakePgreplicas) Update(ctx context.Context, pgreplica *qingcloudcomv1.Pgreplica, opts v1.UpdateOptions) (result *qingcloud.comcomv1.Pgreplica, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(pgreplicasResource, c.ns, pgreplica), &qingcloudcomv1.Pgreplica{})

	if obj == nil {
		return nil, err
	}
	return obj.(*qingcloudcomv1.Pgreplica), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePgreplicas) UpdateStatus(ctx context.Context, pgreplica *qingcloudcomv1.Pgreplica, opts v1.UpdateOptions) (*qingcloud.comcomv1.Pgreplica, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(pgreplicasResource, "status", c.ns, pgreplica), &qingcloudcomv1.Pgreplica{})

	if obj == nil {
		return nil, err
	}
	return obj.(*qingcloudcomv1.Pgreplica), err
}

// Delete takes name of the pgreplica and deletes it. Returns an error if one occurs.
func (c *FakePgreplicas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(pgreplicasResource, c.ns, name), &qingcloudcomv1.Pgreplica{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePgreplicas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(pgreplicasResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &qingcloudcomv1.PgreplicaList{})
	return err
}

// Patch applies the patch and returns the patched pgreplica.
func (c *FakePgreplicas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *qingcloudcomv1.Pgreplica, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(pgreplicasResource, c.ns, name, pt, data, subresources...), &qingcloudcomv1.Pgreplica{})

	if obj == nil {
		return nil, err
	}
	return obj.(*qingcloudcomv1.Pgreplica), err
}
