package proxyservice

import (
	"slices"
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestValidateChainGraphAcceptsConnectedCandidateGroups(t *testing.T) {
	nodes := []domain.Node{
		{ID: "node-a", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-e", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-b", Enabled: true},
		{ID: "node-f", Enabled: true},
		{ID: "node-c", ScopeKey: "scope-c", Enabled: true},
	}
	links := []domain.NodeLink{
		{SourceNodeID: "node-a", TargetNodeID: "node-b"},
		{SourceNodeID: "node-a", TargetNodeID: "node-f"},
		{SourceNodeID: "node-e", TargetNodeID: "node-b"},
		{SourceNodeID: "node-b", TargetNodeID: "node-c"},
		{SourceNodeID: "node-f", TargetNodeID: "node-c"},
	}
	groups := []proxy.ChainHopGroup{
		{Candidates: []string{"node-a", "node-e"}},
		{Candidates: []string{"node-b", "node-f"}},
		{Candidates: []string{"node-c"}},
	}

	if errors := validateChainGraph(groups, nodes, links, "scope-c"); len(errors) != 0 {
		t.Fatalf("validateChainGraph() errors = %v", errors)
	}
}

func TestValidateChainGraphRejectsInvalidCandidates(t *testing.T) {
	nodes := []domain.Node{
		{ID: "node-a", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-e", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-b", Enabled: true},
		{ID: "node-c", ScopeKey: "scope-c", Enabled: true},
		{ID: "node-wrong", ScopeKey: "scope-other", Enabled: true},
	}
	links := []domain.NodeLink{
		{SourceNodeID: "node-a", TargetNodeID: "node-b"},
		{SourceNodeID: "node-b", TargetNodeID: "node-c"},
	}
	tests := []struct {
		name   string
		groups []proxy.ChainHopGroup
		want   string
	}{
		{name: "empty group", groups: []proxy.ChainHopGroup{{Candidates: []string{"node-a"}}, {}}, want: "empty_hop_group:1"},
		{name: "duplicate", groups: []proxy.ChainHopGroup{{Candidates: []string{"node-a"}}, {Candidates: []string{"node-a"}}}, want: "duplicate_candidate:node-a"},
		{name: "isolated standby", groups: []proxy.ChainHopGroup{{Candidates: []string{"node-a", "node-e"}}, {Candidates: []string{"node-b"}}, {Candidates: []string{"node-c"}}}, want: "isolated_outbound_candidate:node-e:0"},
		{name: "wrong destination", groups: []proxy.ChainHopGroup{{Candidates: []string{"node-a"}}, {Candidates: []string{"node-wrong"}}}, want: "destination_scope_mismatch:node-wrong"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errors := validateChainGraph(test.groups, nodes, links, "scope-c")
			if !slices.Contains(errors, test.want) {
				t.Fatalf("validateChainGraph() errors = %v, want %q", errors, test.want)
			}
		})
	}
}

func TestResolveProbeCandidatePathUsesHealthyStandby(t *testing.T) {
	nodes := []domain.Node{
		{ID: "node-a", Name: "A", Enabled: true},
		{ID: "node-e", Name: "E", Enabled: true},
		{ID: "node-b", Name: "B", Enabled: true},
		{ID: "node-c", Name: "C", Enabled: true},
	}
	links := []domain.NodeLink{
		{SourceNodeID: "node-e", TargetNodeID: "node-b"},
		{SourceNodeID: "node-b", TargetNodeID: "node-c"},
	}
	transports := []domain.NodeTransport{
		{NodeID: "node-a", TransportType: domain.TransportTypePublicHTTPS, Status: domain.TransportStatusFailed},
		{NodeID: "node-e", TransportType: domain.TransportTypePublicHTTPS, Status: domain.TransportStatusAvailable},
		{NodeID: "node-b", ParentNodeID: "node-e", TransportType: domain.TransportTypeReverseWSParent, Status: domain.TransportStatusConnected},
		{NodeID: "node-c", ParentNodeID: "node-b", TransportType: domain.TransportTypeReverseWSParent, Status: domain.TransportStatusConnected},
	}
	groups := []proxy.ChainHopGroup{
		{Candidates: []string{"node-a", "node-e"}},
		{Candidates: []string{"node-b"}},
		{Candidates: []string{"node-c"}},
	}

	hops, blockingGroup, reason := resolveProbeCandidatePath(groups, nodes, links, transports)
	if blockingGroup != -1 || reason != "" {
		t.Fatalf("blockingGroup=%d reason=%q", blockingGroup, reason)
	}
	got := make([]string, 0, len(hops))
	for _, hop := range hops {
		got = append(got, hop.NodeID)
	}
	want := []string{"node-e", "node-b", "node-c"}
	if !slices.Equal(got, want) {
		t.Fatalf("resolved node IDs = %v, want %v", got, want)
	}
}
