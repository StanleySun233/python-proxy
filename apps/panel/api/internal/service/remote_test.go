package service

import (
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
)

func TestResolveRemoteAccessPathIDUsesProtocolDefault(t *testing.T) {
	defaults := []domain.RemoteAccessDefault{
		{Protocol: domain.RemoteProtocolSSH, AccessPathID: "ssh-path"},
		{Protocol: domain.RemoteProtocolRDP, AccessPathID: "rdp-path"},
	}
	if got := resolveRemoteAccessPathID("", domain.RemoteProtocolSSH, defaults); got != "ssh-path" {
		t.Fatalf("SSH path = %q", got)
	}
	if got := resolveRemoteAccessPathID("explicit", domain.RemoteProtocolSSH, defaults); got != "explicit" {
		t.Fatalf("explicit path = %q", got)
	}
}

func TestValidateRemoteDefaultPathRequiresProtocolMatch(t *testing.T) {
	path := domain.NodeAccessPath{Enabled: true, Mode: domain.PathModeTCP, Protocol: domain.AccessProtocolTCP, ServiceType: domain.AccessServiceTCPAccess, RemoteProtocol: domain.RemoteProtocolRDP}
	if got := validateRemoteDefaultPath(path, domain.RemoteProtocolSSH); got != "remote_default_protocol_mismatch" {
		t.Fatalf("validation = %q", got)
	}
	path.RemoteProtocol = domain.RemoteProtocolSSH
	if got := validateRemoteDefaultPath(path, domain.RemoteProtocolSSH); got != "" {
		t.Fatalf("validation = %q", got)
	}
}
