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

package compute

import (
	"testing"

	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

// Some arbitrary MAC addresses.
const (
	macAddr1 = "d1:b2:18:25:70:ab"
	macAddr2 = "d1:b2:18:25:70:ac"
	macAddr3 = "d1:b2:18:25:70:ad"
	macAddr4 = "d1:b2:18:25:70:ae"
)

// An address structure as generated by OpenStack. e.g.
// https://docs.openstack.org/api-ref/compute/?expanded=show-server-details-detail#show-server-details
type networkAddress struct {
	Version int    `json:"version"`
	Addr    string `json:"addr"`
	Type    string `json:"OS-EXT-IPS:type"`
	MacAddr string `json:"OS-EXT-IPS:mac_addr"`
}

func serverWithAddresses(addresses map[string][]networkAddress) *ServerExt {
	var server ServerExt

	server.Addresses = make(map[string]interface{})
	for network, addressList := range addresses {
		server.Addresses[network] = addressList
	}

	return &server
}

func TestNetworkStatus_Addresses(t *testing.T) {
	tests := []struct {
		name      string
		addresses map[string][]networkAddress
		want      []corev1.NodeAddress
	}{
		{
			name: "Single network single address",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					},
				},
			},
			want: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "192.168.0.1",
				},
			},
		},
		{
			name: "Fixed and floating addresses",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr2,
					},
				},
			},
			want: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "192.168.0.1",
				}, {
					Type:    corev1.NodeExternalIP,
					Address: "10.0.0.1",
				},
			},
		},
		{
			name: "Ignore IPv6",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 6,
						Addr:    "fe80::f816:3eff:fe56:3174",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 6,
						Addr:    "fe80::f816:3eff:fe56:3175",
						Type:    "floating",
						MacAddr: macAddr2,
					}, {
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr3,
					},
				},
			},
			want: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "192.168.0.1",
				},
			},
		},
		{
			name: "Multiple networks",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr2,
					},
				},
				"extraNet1": {
					{
						Version: 4,
						Addr:    "192.168.1.1",
						Type:    "fixed",
						MacAddr: macAddr3,
					},
				},
				"extraNet2": {
					{
						Version: 4,
						Addr:    "192.168.2.1",
						Type:    "fixed",
						MacAddr: macAddr4,
					},
				},
			},
			want: []corev1.NodeAddress{
				{
					Type:    corev1.NodeInternalIP,
					Address: "192.168.1.1",
				}, {
					Type:    corev1.NodeInternalIP,
					Address: "192.168.2.1",
				}, {
					Type:    corev1.NodeInternalIP,
					Address: "192.168.0.1",
				}, {
					Type:    corev1.NodeExternalIP,
					Address: "10.0.0.1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			is := &InstanceStatus{
				server: serverWithAddresses(tt.addresses),
				logger: logr.Discard(),
			}
			instanceNS, err := is.NetworkStatus()
			g.Expect(err).NotTo(HaveOccurred())

			got := instanceNS.Addresses()
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestInstanceNetworkStatus(t *testing.T) {
	tests := []struct {
		name           string
		addresses      map[string][]networkAddress
		networkName    string
		wantIP         string
		wantFloatingIP string
	}{
		{
			name: "Single network single address",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "",
		},
		{
			name: "Fixed and floating addresses",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr2,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "10.0.0.1",
		},
		{
			name: "Ignore IPv6",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 6,
						Addr:    "fe80::f816:3eff:fe56:3174",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 6,
						Addr:    "fe80::f816:3eff:fe56:3175",
						Type:    "floating",
						MacAddr: macAddr2,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr3,
					}, {
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr4,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "10.0.0.1",
		},
		{
			name: "Ignore unknown address type",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.2",
						Type:    "not-valid",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "192.168.0.3",
						Type:    "unknown",
						MacAddr: macAddr2,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr3,
					}, {
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr4,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "10.0.0.1",
		},
		{
			name: "Multiple networks",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr2,
					},
				},
				"extraNet1": {
					{
						Version: 4,
						Addr:    "192.168.1.1",
						Type:    "fixed",
						MacAddr: macAddr3,
					},
				},
				"extraNet2": {
					{
						Version: 4,
						Addr:    "192.168.2.1",
						Type:    "fixed",
						MacAddr: macAddr4,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "10.0.0.1",
		},
		{
			name: "First IP",
			addresses: map[string][]networkAddress{
				"primary": {
					{
						Version: 4,
						Addr:    "192.168.0.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					}, {
						Version: 4,
						Addr:    "10.0.0.1",
						Type:    "floating",
						MacAddr: macAddr2,
					}, {
						Version: 4,
						Addr:    "192.168.0.2",
						Type:    "fixed",
						MacAddr: macAddr3,
					}, {
						Version: 4,
						Addr:    "10.0.0.2",
						Type:    "floating",
						MacAddr: macAddr4,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "192.168.0.1",
			wantFloatingIP: "10.0.0.1",
		},
		{
			name: "Network not found",
			addresses: map[string][]networkAddress{
				"extraNet1": {
					{
						Version: 4,
						Addr:    "192.168.1.1",
						Type:    "fixed",
						MacAddr: macAddr1,
					},
					{
						Version: 4,
						Addr:    "10.0.1.1",
						Type:    "floating",
						MacAddr: macAddr2,
					},
				},
			},
			networkName:    "primary",
			wantIP:         "",
			wantFloatingIP: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			is := &InstanceStatus{
				server: serverWithAddresses(tt.addresses),
				logger: logr.Discard(),
			}
			ns, err := is.NetworkStatus()
			g.Expect(err).NotTo(HaveOccurred())

			ip := ns.IP(tt.networkName)
			g.Expect(ip).To(Equal(tt.wantIP))

			floatingIP := ns.FloatingIP(tt.networkName)
			g.Expect(floatingIP).To(Equal(tt.wantFloatingIP))
		})
	}
}
