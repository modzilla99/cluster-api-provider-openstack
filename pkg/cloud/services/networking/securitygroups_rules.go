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

package networking

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha6"
)

var defaultRules = []infrav1.SecurityGroupRule{
	{
		Direction:      "egress",
		Description:    "Full open",
		EtherType:      "IPv4",
		PortRangeMin:   0,
		PortRangeMax:   0,
		Protocol:       "",
		RemoteIPPrefix: "",
	},
	{
		Direction:      "egress",
		Description:    "Full open",
		EtherType:      "IPv6",
		PortRangeMin:   0,
		PortRangeMax:   0,
		Protocol:       "",
		RemoteIPPrefix: "",
	},
}

// Permit traffic for etcd, kubelet.
func getSGControlPlaneCommon(remoteGroupIDSelf, secWorkerGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "Etcd",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  2379,
			PortRangeMax:  2380,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			// kubeadm says this is needed
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			// This is needed to support metrics-server deployments
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit traffic for calico.
func getSGControlPlaneCalico(remoteGroupIDSelf, secWorkerGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "BGP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  179,
			PortRangeMax:  179,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "BGP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  179,
			PortRangeMax:  179,
			Protocol:      "tcp",
			RemoteGroupID: secWorkerGroupID,
		},
		{
			Description:   "IP-in-IP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			Protocol:      "ipip",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "IP-in-IP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			Protocol:      "ipip",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit traffic for cilium.
func getSGControlPlaneCilium(remoteGroupIDSelf, secWorkerGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "HealthChecks (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  4240,
			PortRangeMax:  4240,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "HealthChecks (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  4240,
			PortRangeMax:  4240,
			Protocol:      "tcp",
			RemoteGroupID: secWorkerGroupID,
		},
		{
			Description:   "VXLAN (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8472,
			PortRangeMax:  8472,
			Protocol:      "udp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "VXLAN (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8472,
			PortRangeMax:  8472,
			Protocol:      "udp",
			RemoteGroupID: secWorkerGroupID,
		},
		{
			Description:   "ICMP HealthCheck (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8,
			PortRangeMax:  0,
			Protocol:      "icmp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "ICMP HealthCheck (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8,
			PortRangeMax:  0,
			Protocol:      "icmp",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit traffic for kubelet.
func getSGWorkerCommon(remoteGroupIDSelf, secControlPlaneGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			// This is needed to support metrics-server deployments
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "Kubelet API",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  10250,
			PortRangeMax:  10250,
			Protocol:      "tcp",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Permit traffic for calico.
func getSGWorkerCalico(remoteGroupIDSelf, secControlPlaneGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "BGP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  179,
			PortRangeMax:  179,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "BGP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  179,
			PortRangeMax:  179,
			Protocol:      "tcp",
			RemoteGroupID: secControlPlaneGroupID,
		},
		{
			Description:   "IP-in-IP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			Protocol:      "ipip",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "IP-in-IP (calico)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			Protocol:      "ipip",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Permit traffic for cilium.
func getSGWorkerCilium(remoteGroupIDSelf, secControlPlaneGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "HealthChecks (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  4240,
			PortRangeMax:  4240,
			Protocol:      "tcp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "HealthChecks (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  4240,
			PortRangeMax:  4240,
			Protocol:      "tcp",
			RemoteGroupID: secControlPlaneGroupID,
		},
		{
			Description:   "VXLAN (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8472,
			PortRangeMax:  8472,
			Protocol:      "udp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "VXLAN (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8472,
			PortRangeMax:  8472,
			Protocol:      "udp",
			RemoteGroupID: secControlPlaneGroupID,
		},
		{
			Description:   "ICMP HealthCheck (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8,
			PortRangeMax:  0,
			Protocol:      "icmp",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "ICMP HealthCheck (cilium)",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  8,
			PortRangeMax:  0,
			Protocol:      "icmp",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

// Permit traffic for ssh control plane.
func GetSGControlPlaneSSH(secBastionGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "SSH",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  22,
			PortRangeMax:  22,
			Protocol:      "tcp",
			RemoteGroupID: secBastionGroupID,
		},
	}
}

// Permit traffic for ssh worker.
func GetSGWorkerSSH(secBastionGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "SSH",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  22,
			PortRangeMax:  22,
			Protocol:      "tcp",
			RemoteGroupID: secBastionGroupID,
		},
	}
}

// Allow all traffic, including from outside the cluster, to access the API.
func GetSGControlPlaneHTTPS() []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:  "Kubernetes API",
			Direction:    "ingress",
			EtherType:    "IPv4",
			PortRangeMin: 6443,
			PortRangeMax: 6443,
			Protocol:     "tcp",
		},
	}
}

// Allow all traffic, including from outside the cluster, to access node port services.
func GetSGWorkerNodePort() []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:  "Node Port Services",
			Direction:    "ingress",
			EtherType:    "IPv4",
			PortRangeMin: 30000,
			PortRangeMax: 32767,
			Protocol:     "tcp",
		},
	}
}

// Permit all ingress from the cluster security groups.
func GetSGControlPlaneAllowAll(remoteGroupIDSelf, secWorkerGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: secWorkerGroupID,
		},
	}
}

// Permit all ingress from the cluster security groups.
func GetSGWorkerAllowAll(remoteGroupIDSelf, secControlPlaneGroupID string) []infrav1.SecurityGroupRule {
	return []infrav1.SecurityGroupRule{
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: remoteGroupIDSelf,
		},
		{
			Description:   "In-cluster Ingress",
			Direction:     "ingress",
			EtherType:     "IPv4",
			PortRangeMin:  0,
			PortRangeMax:  0,
			Protocol:      "",
			RemoteGroupID: secControlPlaneGroupID,
		},
	}
}

func GetSGControlPlaneGeneral(remoteGroupIDSelf, secWorkerGroupID string) []infrav1.SecurityGroupRule {
	controlPlaneRules := []infrav1.SecurityGroupRule{}
	controlPlaneRules = append(controlPlaneRules, getSGControlPlaneCommon(remoteGroupIDSelf, secWorkerGroupID)...)
	controlPlaneRules = append(controlPlaneRules, getSGControlPlaneCalico(remoteGroupIDSelf, secWorkerGroupID)...)
	controlPlaneRules = append(controlPlaneRules, getSGControlPlaneCilium(remoteGroupIDSelf, secWorkerGroupID)...)
	return controlPlaneRules
}

func GetSGWorkerGeneral(remoteGroupIDSelf, secControlPlaneGroupID string) []infrav1.SecurityGroupRule {
	workerRules := []infrav1.SecurityGroupRule{}
	workerRules = append(workerRules, getSGWorkerCommon(remoteGroupIDSelf, secControlPlaneGroupID)...)
	workerRules = append(workerRules, getSGWorkerCalico(remoteGroupIDSelf, secControlPlaneGroupID)...)
	workerRules = append(workerRules, getSGWorkerCilium(remoteGroupIDSelf, secControlPlaneGroupID)...)
	return workerRules
}
