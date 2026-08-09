package proxyservice

import (
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
)

func TestValidateRemoteProtocolPathRequiresTCPAccess(t *testing.T) {
	tests := []struct {
		name        string
		remote      string
		mode        string
		protocol    string
		serviceType string
		want        string
	}{
		{name: "ssh tcp", remote: domain.RemoteProtocolSSH, mode: domain.PathModeTCP, protocol: domain.AccessProtocolTCP, serviceType: domain.AccessServiceTCPAccess},
		{name: "rdp tcp", remote: domain.RemoteProtocolRDP, mode: domain.PathModeTCP, protocol: domain.AccessProtocolTCP, serviceType: domain.AccessServiceTCPAccess},
		{name: "ssh forward", remote: domain.RemoteProtocolSSH, mode: domain.PathModeForward, protocol: domain.AccessProtocolHTTP, serviceType: domain.AccessServiceHTTPForwardProxy, want: "invalid_remote_protocol_path"},
		{name: "unknown", remote: "vnc", mode: domain.PathModeTCP, protocol: domain.AccessProtocolTCP, serviceType: domain.AccessServiceTCPAccess, want: "invalid_remote_protocol"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validateRemoteProtocolPath(test.remote, test.mode, test.protocol, test.serviceType); got != test.want {
				t.Fatalf("validateRemoteProtocolPath() = %q, want %q", got, test.want)
			}
		})
	}
}
