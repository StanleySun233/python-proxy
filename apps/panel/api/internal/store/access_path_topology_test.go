package store

import (
	"reflect"
	"testing"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestDeriveNodeAccessPathUsesChainTopology(t *testing.T) {
	path := domain.NodeAccessPath{ID: "path-1", ChainID: "chain-1", ListenPort: 2200}
	chain := proxy.Chain{ID: "chain-1", HopGroups: []proxy.ChainHopGroup{
		{Candidates: []string{"a", "e"}},
		{Candidates: []string{"b", "f"}},
		{Candidates: []string{"c"}},
	}}
	nodes := map[string]domain.Node{
		"a": {ID: "a", PublicHost: "a.example", PublicPort: 2988, Status: "healthy"},
		"e": {ID: "e", PublicHost: "e.example", PublicPort: 3988, Status: "degraded"},
		"b": {ID: "b"},
		"f": {ID: "f"},
		"c": {ID: "c"},
	}

	got := deriveNodeAccessPath(path, chain, nodes)

	if got.EntryNodeID != "a" || got.TargetNodeID != "c" || !reflect.DeepEqual(got.RelayNodeIDs, []string{"b"}) {
		t.Fatalf("primary path = entry:%s relay:%v target:%s", got.EntryNodeID, got.RelayNodeIDs, got.TargetNodeID)
	}
	if len(got.Entrypoints) != 2 || got.Entrypoints[0].NodeID != "a" || got.Entrypoints[1].NodeID != "e" {
		t.Fatalf("entrypoints = %+v", got.Entrypoints)
	}
	if len(got.TopologyGroups) != 3 || !reflect.DeepEqual(got.TopologyGroups[1].Candidates, []string{"b", "f"}) {
		t.Fatalf("topology groups = %+v", got.TopologyGroups)
	}
}
