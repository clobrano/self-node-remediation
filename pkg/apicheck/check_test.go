package apicheck

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
)

func TestPopPeerIps(t *testing.T) {
	tests := []struct {
		name                  string
		inputPeers            []corev1.PodIP
		inputCount            int
		expectedOutPeers      []corev1.PodIP
		expectedExcludedPeers []corev1.PodIP
	}{
		{
			name:                  "Take some peers (Happy path)",
			inputPeers:            []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			inputCount:            2,
			expectedOutPeers:      []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}},
			expectedExcludedPeers: []corev1.PodIP{{IP: "9.10.11.12"}},
		},
		{
			name:                  "No Peers",
			inputPeers:            []corev1.PodIP{},
			inputCount:            1,
			expectedOutPeers:      []corev1.PodIP{},
			expectedExcludedPeers: []corev1.PodIP{},
		},

		{
			name:                  "More count than peers",
			inputPeers:            []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			inputCount:            5,
			expectedOutPeers:      []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			expectedExcludedPeers: []corev1.PodIP{},
		},
		{
			name:                  "Some IPs are empty",
			inputPeers:            []corev1.PodIP{{IP: "1.2.3.4"}, {IP: ""}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			inputCount:            2,
			expectedOutPeers:      []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}},
			expectedExcludedPeers: []corev1.PodIP{{IP: "9.10.11.12"}},
		},
		{
			name:                  "Take all peers",
			inputPeers:            []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			inputCount:            3,
			expectedOutPeers:      []corev1.PodIP{{IP: "1.2.3.4"}, {IP: "5.6.7.8"}, {IP: "9.10.11.12"}},
			expectedExcludedPeers: []corev1.PodIP{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ApiConnectivityCheckConfig{
				Log: logr.Logger{},
			}
			c := ApiConnectivityCheck{
				config: &config,
			}

			// Copy the input peers as the function can modify it and
			// we want to compare the result later.
			inputPeers := make([]corev1.PodIP, len(tt.inputPeers))
			copy(inputPeers, tt.inputPeers)

			// Call the UUT
			outPeers := c.popPeerIPs(&inputPeers, tt.inputCount)

			// Check the results
			if !reflect.DeepEqual(outPeers, tt.expectedOutPeers) {
				t.Errorf("selected peers expected: %v, got: %v", tt.expectedOutPeers, outPeers)
			}
			if !reflect.DeepEqual(inputPeers, tt.expectedExcludedPeers) {
				t.Errorf("modified input peers expected: %v, got: %v", tt.expectedExcludedPeers, inputPeers)
			}

		})

	}
}
