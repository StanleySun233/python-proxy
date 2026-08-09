package store

import (
	"reflect"
	"testing"

	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func TestGroupChainHopModelsPreservesPositionAndPriority(t *testing.T) {
	models := []ChainHopModel{
		{HopIndex: 0, CandidateIndex: 0, NodeID: "node-a"},
		{HopIndex: 0, CandidateIndex: 1, NodeID: "node-e"},
		{HopIndex: 1, CandidateIndex: 0, NodeID: "node-b"},
	}

	got := groupChainHopModels(models)
	want := []proxy.ChainHopGroup{
		{Candidates: []string{"node-a", "node-e"}},
		{Candidates: []string{"node-b"}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupChainHopModels() = %#v, want %#v", got, want)
	}
}
