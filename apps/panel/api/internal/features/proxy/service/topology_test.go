package proxyservice

import (
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestBuildTopologyProjectionMarksPrimaryAndStandbyEdges(t *testing.T) {
	chain := proxy.Chain{ID: "chain-1", Name: "resilient", DestinationScope: "scope-c", Enabled: true, HopGroups: []proxy.ChainHopGroup{
		{Candidates: []string{"a", "e"}},
		{Candidates: []string{"b", "f"}},
		{Candidates: []string{"c"}},
	}}
	nodes := []domain.Node{
		{ID: "a", Name: "A", Enabled: true, Status: domain.NodeStatusHealthy},
		{ID: "e", Name: "E", Enabled: true, Status: domain.NodeStatusHealthy},
		{ID: "b", Name: "B", Enabled: true, Status: domain.NodeStatusHealthy},
		{ID: "f", Name: "F", Enabled: true, Status: domain.NodeStatusHealthy},
		{ID: "c", Name: "C", Enabled: true, Status: domain.NodeStatusHealthy},
	}
	links := []domain.NodeLink{
		{SourceNodeID: "a", TargetNodeID: "b"},
		{SourceNodeID: "a", TargetNodeID: "f"},
		{SourceNodeID: "e", TargetNodeID: "b"},
		{SourceNodeID: "b", TargetNodeID: "c"},
		{SourceNodeID: "f", TargetNodeID: "c"},
	}
	paths := []domain.NodeAccessPath{{ID: "path-1", ChainID: "chain-1", Name: "SSH path", Enabled: true}}
	probes := map[string]proxy.ChainProbeResult{"chain-1": {ChainID: "chain-1", Status: "connected", ResolvedHops: []proxy.ChainProbeHop{{NodeID: "a"}, {NodeID: "b"}, {NodeID: "c"}}}}

	projection := buildTopologyProjection(paths, []proxy.Chain{chain}, nodes, links, nil, probes)
	if len(projection.Paths) != 1 {
		t.Fatalf("paths = %d", len(projection.Paths))
	}
	edges := projection.Paths[0].Edges
	assertTopologyEdge(t, edges, "a", "b", "primary", true)
	assertTopologyEdge(t, edges, "a", "f", "standby", false)
	assertTopologyEdge(t, edges, "e", "b", "standby", false)
}

func assertTopologyEdge(t *testing.T, edges []proxy.TopologyEdge, source string, target string, kind string, active bool) {
	t.Helper()
	for _, edge := range edges {
		if edge.SourceNodeID == source && edge.TargetNodeID == target {
			if edge.Kind != kind || edge.Active != active {
				t.Fatalf("edge %s -> %s = %+v", source, target, edge)
			}
			return
		}
	}
	t.Fatalf("edge %s -> %s not found", source, target)
}
