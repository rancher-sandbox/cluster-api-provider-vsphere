/*
Copyright 2022 The Kubernetes Authors.

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

package flavors

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	bootstrapv1 "sigs.k8s.io/cluster-api/api/bootstrap/kubeadm/v1beta2"
	controlplanev1 "sigs.k8s.io/cluster-api/api/controlplane/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/yaml"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/clusterclass"
	"sigs.k8s.io/cluster-api-provider-vsphere/packaging/flavorgen/flavors/env"
	"sigs.k8s.io/cluster-api-provider-vsphere/packaging/flavorgen/flavors/kubevip"
	"sigs.k8s.io/cluster-api-provider-vsphere/packaging/flavorgen/flavors/util"
)

func newClusterClass() clusterv1.ClusterClass {
	return clusterv1.ClusterClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1.GroupVersion.String(),
			Kind:       util.TypeToKind(&clusterv1.ClusterClass{}),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: env.ClusterClassNameVar,
		},
		Spec: clusterv1.ClusterClassSpec{
			Infrastructure: clusterv1.InfrastructureClass{
				TemplateRef: clusterv1.ClusterClassTemplateReference{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       util.TypeToKind(&infrav1.VSphereClusterTemplate{}),
					Name:       env.ClusterClassNameVar,
				},
			},
			ControlPlane: getControlPlaneClass(),
			Workers:      getWorkersClass(),
			Variables:    clusterclass.GetClusterClassVariables(true),
			Patches:      getClusterClassPatches(),
		},
	}
}

func newVMWareClusterClass() clusterv1.ClusterClass {
	return clusterv1.ClusterClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterv1.GroupVersion.String(),
			Kind:       util.TypeToKind(&clusterv1.ClusterClass{}),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: env.ClusterClassNameVar,
		},
		Spec: clusterv1.ClusterClassSpec{
			Infrastructure: clusterv1.InfrastructureClass{
				TemplateRef: clusterv1.ClusterClassTemplateReference{
					APIVersion: vmwarev1.GroupVersion.String(),
					Kind:       util.TypeToKind(&vmwarev1.VSphereClusterTemplate{}),
					Name:       env.ClusterClassNameVar,
				},
			},
			ControlPlane: getVMWareControlPlaneClass(),
			Workers:      getVMWareWorkersClass(),
			Variables:    clusterclass.GetClusterClassVariables(false),
			Patches:      getVMWareClusterClassPatches(),
		},
	}
}

func getControlPlaneClass() clusterv1.ControlPlaneClass {
	return clusterv1.ControlPlaneClass{
		TemplateRef: clusterv1.ClusterClassTemplateReference{
			Kind:       util.TypeToKind(&controlplanev1.KubeadmControlPlaneTemplate{}),
			Name:       fmt.Sprintf("%s-controlplane", env.ClusterClassNameVar),
			APIVersion: controlplanev1.GroupVersion.String(),
		},
		MachineInfrastructure: &clusterv1.ControlPlaneClassMachineInfrastructureTemplate{
			TemplateRef: clusterv1.ClusterClassTemplateReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       util.TypeToKind(&infrav1.VSphereMachineTemplate{}),
				Name:       fmt.Sprintf("%s-template", env.ClusterClassNameVar),
			},
		},
	}
}

func getVMWareControlPlaneClass() clusterv1.ControlPlaneClass {
	return clusterv1.ControlPlaneClass{
		TemplateRef: clusterv1.ClusterClassTemplateReference{
			Kind:       util.TypeToKind(&controlplanev1.KubeadmControlPlaneTemplate{}),
			Name:       fmt.Sprintf("%s-controlplane", env.ClusterClassNameVar),
			APIVersion: controlplanev1.GroupVersion.String(),
		},
		MachineInfrastructure: &clusterv1.ControlPlaneClassMachineInfrastructureTemplate{
			TemplateRef: clusterv1.ClusterClassTemplateReference{
				APIVersion: vmwarev1.GroupVersion.String(),
				Kind:       util.TypeToKind(&vmwarev1.VSphereMachineTemplate{}),
				Name:       fmt.Sprintf("%s-template", env.ClusterClassNameVar),
			},
		},
	}
}

func getWorkersClass() clusterv1.WorkersClass {
	return clusterv1.WorkersClass{
		MachineDeployments: []clusterv1.MachineDeploymentClass{
			{
				Class: fmt.Sprintf("%s-worker", env.ClusterClassNameVar),
				Bootstrap: clusterv1.MachineDeploymentClassBootstrapTemplate{
					TemplateRef: clusterv1.ClusterClassTemplateReference{
						APIVersion: bootstrapv1.GroupVersion.String(),
						Kind:       util.TypeToKind(&bootstrapv1.KubeadmConfigTemplate{}),
						Name:       fmt.Sprintf("%s-worker-bootstrap-template", env.ClusterClassNameVar),
					},
				},
				Infrastructure: clusterv1.MachineDeploymentClassInfrastructureTemplate{
					TemplateRef: clusterv1.ClusterClassTemplateReference{
						Kind:       util.TypeToKind(&infrav1.VSphereMachineTemplate{}),
						Name:       fmt.Sprintf("%s-worker-machinetemplate", env.ClusterClassNameVar),
						APIVersion: infrav1.GroupVersion.String(),
					},
				},
			},
		},
	}
}

func getVMWareWorkersClass() clusterv1.WorkersClass {
	return clusterv1.WorkersClass{
		MachineDeployments: []clusterv1.MachineDeploymentClass{
			{
				Class: fmt.Sprintf("%s-worker", env.ClusterClassNameVar),
				Bootstrap: clusterv1.MachineDeploymentClassBootstrapTemplate{
					TemplateRef: clusterv1.ClusterClassTemplateReference{
						APIVersion: bootstrapv1.GroupVersion.String(),
						Kind:       util.TypeToKind(&bootstrapv1.KubeadmConfigTemplate{}),
						Name:       fmt.Sprintf("%s-worker-bootstrap-template", env.ClusterClassNameVar),
					},
				},
				Infrastructure: clusterv1.MachineDeploymentClassInfrastructureTemplate{
					TemplateRef: clusterv1.ClusterClassTemplateReference{
						Kind:       util.TypeToKind(&vmwarev1.VSphereMachineTemplate{}),
						Name:       fmt.Sprintf("%s-worker-machinetemplate", env.ClusterClassNameVar),
						APIVersion: vmwarev1.GroupVersion.String(),
					},
				},
			},
		},
	}
}

func getClusterClassPatches() []clusterv1.ClusterClassPatch {
	return []clusterv1.ClusterClassPatch{
		createEmptyArraysPatch(),
		enableSSHPatch(),
		infraClusterPatch(),
		kubevip.TopologyPatch(),
	}
}

func getVMWareClusterClassPatches() []clusterv1.ClusterClassPatch {
	return []clusterv1.ClusterClassPatch{
		createEmptyArraysPatch(),
		enableSSHPatch(),
		vmWareInfraClusterPatch(),
		kubevip.TopologyPatch(),
	}
}

func getCredSecretNameTemplate() string {
	template := map[string]interface{}{
		"name": "{{ .credsSecretName }}",
		"kind": "Secret",
	}
	templateStr, _ := yaml.Marshal(template)
	return string(templateStr)
}

func getControlPlaneEndpointTemplate() string {
	template := map[string]interface{}{
		"host": "{{ .controlPlaneIpAddr }}",
		"port": "{{ .controlPlanePort }}",
	}
	templateStr, _ := yaml.Marshal(template)

	fixTemplateStr := string(templateStr)
	fixTemplateStr = strings.ReplaceAll(fixTemplateStr, "'{{ .controlPlanePort }}'", "{{ .controlPlanePort }}")
	return fixTemplateStr
}

func getEnableSSHIntoNodesTemplate() string {
	template := []map[string]interface{}{
		{
			"name": "capv",
			"sshAuthorizedKeys": []string{
				"{{ .sshKey }}",
			},
			"sudo": "ALL=(ALL) NOPASSWD:ALL",
		},
	}
	templateStr, _ := yaml.Marshal(template)
	return string(templateStr)
}

func newVSphereClusterTemplate() infrav1.VSphereClusterTemplate {
	return infrav1.VSphereClusterTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       util.TypeToKind(&infrav1.VSphereClusterTemplate{}),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      env.ClusterClassNameVar,
			Namespace: env.NamespaceVar,
		},
		Spec: infrav1.VSphereClusterTemplateSpec{
			Template: infrav1.VSphereClusterTemplateResource{
				Spec: infrav1.VSphereClusterSpec{},
			},
		},
	}
}

func newVMWareClusterTemplate() vmwarev1.VSphereClusterTemplate {
	return vmwarev1.VSphereClusterTemplate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vmwarev1.GroupVersion.String(),
			Kind:       util.TypeToKind(&vmwarev1.VSphereClusterTemplate{}),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      env.ClusterClassNameVar,
			Namespace: env.NamespaceVar,
		},
		Spec: vmwarev1.VSphereClusterTemplateSpec{
			Template: vmwarev1.VSphereClusterTemplateResource{
				Spec: vmwarev1.VSphereClusterSpec{},
			},
		},
	}
}

func newKubeadmControlPlaneTemplate(templateName string) controlplanev1.KubeadmControlPlaneTemplate {
	return controlplanev1.KubeadmControlPlaneTemplate{
		TypeMeta: metav1.TypeMeta{
			Kind:       util.TypeToKind(&controlplanev1.KubeadmControlPlaneTemplate{}),
			APIVersion: controlplanev1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      templateName,
			Namespace: env.NamespaceVar,
		},
		Spec: controlplanev1.KubeadmControlPlaneTemplateSpec{
			Template: controlplanev1.KubeadmControlPlaneTemplateResource{
				Spec: controlplanev1.KubeadmControlPlaneTemplateResourceSpec{
					KubeadmConfigSpec: defaultKubeadmInitSpec([]bootstrapv1.File{}),
				},
			},
		},
	}
}
