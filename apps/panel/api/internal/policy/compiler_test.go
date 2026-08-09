package policy

import (
	"strings"
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestCompileRejectsUnsupportedNodeMatchType(t *testing.T) {
	_, err := Compile(nil, nil, nil, []proxy.RouteRule{
		{ID: "route-1", Enabled: true, MatchType: "url_regex", MatchValue: "^https://", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
	}, nil)
	if err == nil || !strings.Contains(err.Error(), "unsupported_match_type") {
		t.Fatalf("err = %v", err)
	}
}

func TestCompileAllowsNodeSupportedMatchTypes(t *testing.T) {
	nodes := []domain.Node{{ID: "node-a", ScopeKey: "scope-a", Enabled: true}}
	rules := []proxy.RouteRule{
		{ID: "domain", Enabled: true, MatchType: domain.MatchTypeDomain, MatchValue: "example.com", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
		{ID: "suffix", Enabled: true, MatchType: domain.MatchTypeDomainSuffix, MatchValue: ".example.com", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
		{ID: "ip", Enabled: true, MatchType: domain.MatchTypeIP, MatchValue: "127.0.0.1", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
		{ID: "cidr", Enabled: true, MatchType: domain.MatchTypeIPCIDR, MatchValue: "127.0.0.0/24", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
		{ID: "protocol", Enabled: true, MatchType: domain.MatchTypeProtocol, MatchValue: "https", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
		{ID: "default", Enabled: true, MatchType: domain.MatchTypeDefault, MatchValue: "*", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a"},
	}
	if _, err := Compile(nodes, nil, nil, rules, nil); err != nil {
		t.Fatal(err)
	}
}

func TestCompileRejectsIsolatedCandidate(t *testing.T) {
	nodes := []domain.Node{
		{ID: "node-a", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-e", Mode: domain.NodeModeEdge, Enabled: true},
		{ID: "node-b", Enabled: true},
		{ID: "node-c", ScopeKey: "scope-c", Enabled: true},
	}
	links := []domain.NodeLink{
		{SourceNodeID: "node-a", TargetNodeID: "node-b"},
		{SourceNodeID: "node-b", TargetNodeID: "node-c"},
	}
	chains := []proxy.Chain{{
		ID:               "chain",
		DestinationScope: "scope-c",
		Enabled:          true,
		HopGroups: []proxy.ChainHopGroup{
			{Candidates: []string{"node-a", "node-e"}},
			{Candidates: []string{"node-b"}},
			{Candidates: []string{"node-c"}},
		},
	}}

	_, err := Compile(nodes, links, chains, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "isolated_outbound_candidate node-e") {
		t.Fatalf("err = %v", err)
	}
}

func TestCompileRejectsNonEdgeEntryCandidate(t *testing.T) {
	nodes := []domain.Node{{ID: "node-a", Mode: "relay", ScopeKey: "scope-a", Enabled: true}}
	chains := []proxy.Chain{{ID: "chain", DestinationScope: "scope-a", Enabled: true, HopGroups: []proxy.ChainHopGroup{{Candidates: []string{"node-a"}}}}}

	_, err := Compile(nodes, nil, chains, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "entry_candidate_not_edge node-a") {
		t.Fatalf("err = %v", err)
	}
}
