package proxy

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/StanleySun233/python-proxy/apps/node/api/internal/domain"
	"github.com/StanleySun233/python-proxy/apps/node/api/internal/policystore"
)

type candidateStreamOpener struct {
	attempted []string
}

func (o *candidateStreamOpener) OpenDirectStream(_ context.Context, nextHop domain.Node, _ []string, _ string, _ int) (net.Conn, error) {
	o.attempted = append(o.attempted, nextHop.ID)
	if nextHop.ID == "b" {
		return nil, errors.New("primary down")
	}
	client, server := net.Pipe()
	go func() {
		defer server.Close()
		request, err := http.ReadRequest(bufio.NewReader(server))
		if err != nil {
			return
		}
		request.Body.Close()
		_, _ = fmt.Fprint(server, "HTTP/1.1 200 OK\r\nContent-Length: 7\r\n\r\nstandby")
	}()
	return client, nil
}

func (o *candidateStreamOpener) HasDirectPeer(string) bool {
	return true
}

func TestResolveChainCandidatesPreservesPriorityAndBuildsConnectedPaths(t *testing.T) {
	snapshot := policystore.Snapshot{
		Nodes: []domain.Node{
			{ID: "a", Enabled: true},
			{ID: "e", Enabled: true},
			{ID: "b", Enabled: true, PublicHost: "b.example", PublicPort: 8080},
			{ID: "f", Enabled: true, PublicHost: "f.example", PublicPort: 8080},
			{ID: "c", Enabled: true},
		},
		Links: []domain.NodeLink{
			{SourceNodeID: "a", TargetNodeID: "b"},
			{SourceNodeID: "a", TargetNodeID: "f"},
			{SourceNodeID: "e", TargetNodeID: "b"},
			{SourceNodeID: "b", TargetNodeID: "c"},
			{SourceNodeID: "f", TargetNodeID: "c"},
		},
		Chains: []domain.Chain{{
			ID: "chain-1",
			HopGroups: []domain.ChainHopGroup{
				{Candidates: []string{"a", "e"}},
				{Candidates: []string{"b", "f"}},
				{Candidates: []string{"c"}},
			},
		}},
	}
	server := NewServer(nil, func() string { return "a" }, nil)

	candidates, ok := server.resolveChainCandidates(snapshot, "chain-1")
	if !ok {
		t.Fatal("chain candidates were not resolved")
	}
	if len(candidates) != 2 {
		t.Fatalf("candidate count = %d, want 2", len(candidates))
	}
	if candidates[0].node.ID != "b" || !reflect.DeepEqual(candidates[0].remainingHops, []string{"c"}) {
		t.Fatalf("first candidate = %+v", candidates[0])
	}
	if candidates[1].node.ID != "f" || !reflect.DeepEqual(candidates[1].remainingHops, []string{"c"}) {
		t.Fatalf("second candidate = %+v", candidates[1])
	}
}

func TestResolveChainCandidatesSkipsDisconnectedPath(t *testing.T) {
	snapshot := policystore.Snapshot{
		Nodes: []domain.Node{
			{ID: "a", Enabled: true},
			{ID: "b", Enabled: true, PublicHost: "b.example", PublicPort: 8080},
			{ID: "f", Enabled: true, PublicHost: "f.example", PublicPort: 8080},
			{ID: "c", Enabled: true},
		},
		Links: []domain.NodeLink{
			{SourceNodeID: "a", TargetNodeID: "b"},
			{SourceNodeID: "a", TargetNodeID: "f"},
			{SourceNodeID: "f", TargetNodeID: "c"},
		},
		Chains: []domain.Chain{{
			ID: "chain-1",
			HopGroups: []domain.ChainHopGroup{
				{Candidates: []string{"a"}},
				{Candidates: []string{"b", "f"}},
				{Candidates: []string{"c"}},
			},
		}},
	}
	server := NewServer(nil, func() string { return "a" }, nil)

	candidates, ok := server.resolveChainCandidates(snapshot, "chain-1")
	if !ok || len(candidates) != 1 || candidates[0].node.ID != "f" {
		t.Fatalf("candidates = %+v, ok = %v", candidates, ok)
	}
}

func TestOpenCandidateStreamTriesStandbyAfterPrimaryFailure(t *testing.T) {
	opener := &candidateStreamOpener{}
	hops := []chainHop{{node: domain.Node{ID: "b"}}, {node: domain.Node{ID: "f"}}}

	conn, selected, err := openCandidateStream(context.Background(), opener, nil, hops, "target.internal", 443)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if selected.node.ID != "f" {
		t.Fatalf("selected = %s", selected.node.ID)
	}
	if !reflect.DeepEqual(opener.attempted, []string{"b", "f"}) {
		t.Fatalf("attempted = %v", opener.attempted)
	}
}

func TestForwardChainUsesStandbyAfterPrimaryOpenFailure(t *testing.T) {
	store := policystore.New("")
	payload, err := json.Marshal(policystore.Snapshot{
		Nodes: []domain.Node{
			{ID: "b", Enabled: true},
			{ID: "f", Enabled: true},
			{ID: "c", Enabled: true},
		},
		Links: []domain.NodeLink{
			{SourceNodeID: "a", TargetNodeID: "b"},
			{SourceNodeID: "a", TargetNodeID: "f"},
			{SourceNodeID: "b", TargetNodeID: "c"},
			{SourceNodeID: "f", TargetNodeID: "c"},
		},
		Chains: []domain.Chain{{
			ID: "chain-1",
			HopGroups: []domain.ChainHopGroup{
				{Candidates: []string{"a"}},
				{Candidates: []string{"b", "f"}},
				{Candidates: []string{"c"}},
			},
		}},
		RouteRules: []domain.RouteRule{{ID: "default", MatchType: domain.MatchTypeDefault, ActionType: domain.ActionTypeChain, ChainID: "chain-1", Enabled: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Update("test", string(payload)); err != nil {
		t.Fatal(err)
	}
	opener := &candidateStreamOpener{}
	server := newAuthenticatedForwardServer(t, store)
	server.nodeIDGetter = func() string { return "a" }
	server.SetDirectStreamOpener(opener)
	req := httptest.NewRequest(http.MethodGet, "http://target.internal/data", nil)
	setForwardProxyToken(req)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, req)

	if response.Code != http.StatusOK || response.Body.String() != "standby" {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	if !reflect.DeepEqual(opener.attempted, []string{"b", "f"}) {
		t.Fatalf("attempted = %v", opener.attempted)
	}
}
