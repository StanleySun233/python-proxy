package store

import (
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestDefaultAccessPathInputUsesPublicEntryNode(t *testing.T) {
	nodes := map[string]domain.Node{
		"edge":  {ID: "edge", PublicHost: "103.214.172.211", PublicPort: 2988},
		"relay": {ID: "relay"},
	}
	chain := proxy.Chain{ID: "chain", Name: "hk2astar", HopGroups: []proxy.ChainHopGroup{{Candidates: []string{"edge"}}, {Candidates: []string{"relay"}}}}

	got, ok := defaultAccessPathInput(chain, nodes)
	if !ok {
		t.Fatalf("defaultAccessPathInput ok = false")
	}

	if got.ChainID != "chain" || got.Name != "hk2astar default" {
		t.Fatalf("identity = %+v", got)
	}
	if got.Mode != domain.PathModeForward || got.Protocol != domain.AccessProtocolHTTP || got.ServiceType != domain.AccessServiceHTTPForwardProxy {
		t.Fatalf("mode/protocol/service = %+v", got)
	}
	if got.ListenHost != "103.214.172.211" || got.ListenPort != 2988 {
		t.Fatalf("listen = %+v", got)
	}
	if got.TargetProtocol != domain.AccessProtocolHTTP || got.TargetPort != 2988 || got.AuthMode != domain.AccessAuthProxyToken {
		t.Fatalf("target/auth = %+v", got)
	}
}

func TestDefaultAccessPathInputRejectsMissingHopNode(t *testing.T) {
	nodes := map[string]domain.Node{
		"edge": {ID: "edge", PublicHost: "203.0.113.10", PublicPort: 2988},
	}
	chain := proxy.Chain{ID: "chain", Name: "broken", HopGroups: []proxy.ChainHopGroup{{Candidates: []string{"edge"}}, {Candidates: []string{"missing"}}}}

	if _, ok := defaultAccessPathInput(chain, nodes); ok {
		t.Fatalf("defaultAccessPathInput ok = true")
	}
}

func TestPolicyRouteRulesWithAccessPathsAssignsChainPath(t *testing.T) {
	rules := []proxy.RouteRule{
		{ID: "route-1", ActionType: domain.ActionTypeChain, ChainID: "chain-1", Enabled: true},
		{ID: "route-2", ActionType: domain.ActionTypeDirect, DestinationScope: "scope-a", Enabled: true},
	}
	paths := []domain.NodeAccessPath{
		{ID: "path-disabled", ChainID: "chain-1", Enabled: false},
		{ID: "path-1", ChainID: "chain-1", Enabled: true},
	}

	got := policyRouteRulesWithAccessPaths(rules, paths)

	if got[0].AccessPathID != "path-1" {
		t.Fatalf("chain route accessPathId = %q", got[0].AccessPathID)
	}
	if got[1].AccessPathID != "" {
		t.Fatalf("direct route accessPathId = %q", got[1].AccessPathID)
	}
}
