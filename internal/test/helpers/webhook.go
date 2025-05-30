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

package helpers

import (
	"net"
	"os"
	"strconv"
	"time"

	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/klog/v2"
	utilyaml "sigs.k8s.io/cluster-api/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

const (
	mutatingWebhookKind   = "MutatingWebhookConfiguration"
	validatingWebhookKind = "ValidatingWebhookConfiguration"
	mutatingwebhook       = "mutating-webhook-configuration"
	validatingwebhook     = "validating-webhook-configuration"
)

// Mutate the name of each webhook, because kubebuilder generates the same name for all controllers.
// In normal usage, kustomize will prefix the controller name, which we have to do manually here.
func appendWebhookConfiguration(configyamlFile []byte, tag string) ([]*v1.MutatingWebhookConfiguration, []*v1.ValidatingWebhookConfiguration, error) {
	var mutatingWebhooks []*v1.MutatingWebhookConfiguration
	var validatingWebhooks []*v1.ValidatingWebhookConfiguration
	objs, err := utilyaml.ToUnstructured(configyamlFile)
	if err != nil {
		klog.Fatalf("failed to parse yaml")
	}
	// look for resources of kind MutatingWebhookConfiguration
	for i := range objs {
		o := objs[i]
		if o.GetKind() == mutatingWebhookKind {
			// update the name in metadata
			if o.GetName() == mutatingwebhook {
				var m v1.MutatingWebhookConfiguration
				o.SetName(mutatingwebhook + "-" + tag)
				if err := scheme.Convert(&o, &m, nil); err != nil {
					return nil, nil, err
				}
				mutatingWebhooks = append(mutatingWebhooks, &m)
			}
		}
		if o.GetKind() == validatingWebhookKind {
			// update the name in metadata
			if o.GetName() == validatingwebhook {
				var v v1.ValidatingWebhookConfiguration
				o.SetName(validatingwebhook + "-" + tag)
				if err := scheme.Convert(&o, &v, nil); err != nil {
					return nil, nil, err
				}
				validatingWebhooks = append(validatingWebhooks, &v)
			}
		}
	}
	return mutatingWebhooks, validatingWebhooks, err
}

// InitializeWebhookInEnvironment initializes WebhookInstallOptions for the provided environment.
func InitializeWebhookInEnvironment(e *envtest.Environment, configPath string) {
	if configPath == "" {
		klog.Fatalf("webhook configuration path is empty")
	}
	configyamlFile, err := os.ReadFile(configPath) //nolint:gosec
	if err != nil {
		klog.Fatalf("Failed to read core webhook configuration file: %v ", err)
	}
	if err != nil {
		klog.Fatalf("failed to parse yaml")
	}
	// append the webhook with suffix to avoid clashing webhooks. repeated for every webhook
	mutatingWebhooks, validatingWebhooks, err := appendWebhookConfiguration(configyamlFile, "config")
	if err != nil {
		klog.Fatalf("Failed to append core controller webhook config: %v", err)
	}

	e.WebhookInstallOptions = envtest.WebhookInstallOptions{
		MaxTime:            20 * time.Second,
		PollInterval:       time.Second,
		ValidatingWebhooks: validatingWebhooks,
		MutatingWebhooks:   mutatingWebhooks,
	}
}

// WaitForWebhooks waits for TestEnvironment's webhooks to accept connections.
func (t *TestEnvironment) WaitForWebhooks() {
	port := t.env.WebhookInstallOptions.LocalServingPort

	klog.Infof("Waiting for webhook port %d to be open prior to running tests", port)
	timeout := 1 * time.Second
	for {
		time.Sleep(1 * time.Second)
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), timeout)
		if err != nil {
			klog.Infof("Webhook port is not ready, will retry in %v: %s", timeout, err)
			continue
		}
		if err := conn.Close(); err != nil {
			klog.Info("Connection to webhook port could not be closed. Continuing with tests...")
		}
		klog.Info("Webhook port is now open. Continuing with tests...")
		return
	}
}
