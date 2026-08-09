package service

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
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

func TestOpenRemoteTCPUpstreamUsesNextEntrypointAfterRejection(t *testing.T) {
	first := remoteTCPTestListener(t, "failed", nil)
	defer first.Close()
	frames := make(chan tcpAccessAuthFrame, 1)
	second := remoteTCPTestListener(t, "connected", frames)
	defer second.Close()
	firstPort := first.Addr().(*net.TCPAddr).Port
	secondPort := second.Addr().(*net.TCPAddr).Port
	record := remoteSessionRecord{
		ProxyToken: "token",
		TargetHost: "target.test",
		TargetPort: 22,
		TCPAttempts: []remoteTCPAttempt{
			{Host: "127.0.0.1", Port: firstPort},
			{Host: "127.0.0.1", Port: secondPort, ChainCandidates: []tcpAccessChainCandidate{{NextNodeID: "b"}, {NextNodeID: "f"}}},
		},
	}
	conn, _, err := openRemoteTCPUpstream(record)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	frame := <-frames
	if len(frame.ChainCandidates) != 2 || frame.ChainCandidates[1].NextNodeID != "f" {
		t.Fatalf("frame = %+v", frame)
	}
}

func remoteTCPTestListener(t *testing.T, status string, frames chan<- tcpAccessAuthFrame) net.Listener {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		var frame tcpAccessAuthFrame
		if decodeErr := json.NewDecoder(conn).Decode(&frame); decodeErr != nil {
			return
		}
		if frames != nil {
			frames <- frame
		}
		_, _ = fmt.Fprintf(conn, "{\"status\":%s}\n", strconv.Quote(status))
	}()
	return listener
}

func TestRemoteChainCandidatesEnumeratesStandbyRelayPaths(t *testing.T) {
	path := domain.NodeAccessPath{TopologyGroups: []domain.AccessPathTopologyGroup{
		{Candidates: []string{"a", "e"}},
		{Candidates: []string{"b", "f"}},
		{Candidates: []string{"c"}},
	}}
	want := []tcpAccessChainCandidate{
		{NextNodeID: "b", RemainingHopNodeIDs: []string{"c"}},
		{NextNodeID: "f", RemainingHopNodeIDs: []string{"c"}},
	}
	if got := remoteChainCandidates(path); !reflect.DeepEqual(got, want) {
		t.Fatalf("candidates = %+v", got)
	}
}

func TestRemoteTCPAttemptsPreserveHealthyEntrypointOrder(t *testing.T) {
	path := domain.NodeAccessPath{
		ListenPort: 2990,
		Entrypoints: []domain.AccessEntrypoint{
			{NodeID: "a", Host: "a.test", Status: "degraded"},
			{NodeID: "e", Host: "e.test", Status: "healthy"},
			{NodeID: "g", Host: "g.test", Status: "healthy"},
		},
		TopologyGroups: []domain.AccessPathTopologyGroup{{Candidates: []string{"a", "e", "g"}}, {Candidates: []string{"c"}}},
	}
	want := []remoteTCPAttempt{
		{Host: "e.test", Port: 2990, ChainCandidates: []tcpAccessChainCandidate{{NextNodeID: "c"}}},
		{Host: "g.test", Port: 2990, ChainCandidates: []tcpAccessChainCandidate{{NextNodeID: "c"}}},
	}
	got, err := remoteTCPAttempts(path)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("attempts = %+v err=%v", got, err)
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
