/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha4

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereDeploymentZone to the Hub version (v1beta1).
func (src *VSphereDeploymentZone) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereDeploymentZone)

	if err := Convert_v1alpha4_VSphereDeploymentZone_To_v1beta1_VSphereDeploymentZone(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.VSphereDeploymentZone{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereDeploymentZone.
func (dst *VSphereDeploymentZone) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereDeploymentZone)

	if err := Convert_v1beta1_VSphereDeploymentZone_To_v1alpha4_VSphereDeploymentZone(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this VSphereDeploymentZoneList to the Hub version (v1beta1).
func (src *VSphereDeploymentZoneList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereDeploymentZoneList)
	return Convert_v1alpha4_VSphereDeploymentZoneList_To_v1beta1_VSphereDeploymentZoneList(src, dst, nil)
}

// ConvertFrom converts this VSphereDeploymentZoneList to the Hub version (v1beta1).
func (dst *VSphereDeploymentZoneList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereDeploymentZoneList)
	return Convert_v1beta1_VSphereDeploymentZoneList_To_v1alpha4_VSphereDeploymentZoneList(src, dst, nil)
}
