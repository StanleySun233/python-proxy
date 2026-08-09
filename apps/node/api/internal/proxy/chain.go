package proxy

import (
	"net/http"

	"github.com/StanleySun233/python-proxy/apps/node/api/internal/domain"
	"github.com/StanleySun233/python-proxy/apps/node/api/internal/policystore"
)

func (s *Server) forwardChain(w http.ResponseWriter, req *http.Request, snapshot policystore.Snapshot, rule domain.RouteRule, tracker *proxySessionTracker) {
	hops, ok := s.resolveChainCandidates(snapshot, rule.ChainID)
	if !ok || len(hops) == 0 {
		tracker.finish(0, 0, domain.ProxySessionStatusError, proxyErrorInvalidChainRoute, proxyErrorInvalidChainRoute)
		writeProxyError(w, req, proxyErrorInvalidChainRoute, http.StatusBadGateway)
		return
	}
	hop := hops[0]
	if hop.isLast {
		if isWebSocketUpgrade(req) {
			s.upgradeDirect(w, req, tracker)
			return
		}
		if req.Method == http.MethodConnect {
			s.tunnelDirect(w, req, tracker)
			return
		}
		s.forwardDirect(w, req, tracker)
		return
	}
	if isWebSocketUpgrade(req) {
		s.upgradeViaCandidates(w, req, hops, tracker)
		return
	}
	if req.Method == http.MethodConnect {
		s.tunnelViaCandidates(w, req, hops, tracker)
		return
	}
	s.forwardViaCandidates(w, req, hops, tracker)
}

func (s *Server) shouldUseStream(nextHop domain.Node) bool {
	privateNextHop := nextHop.PublicHost == "" || nextHop.PublicPort <= 0
	return s.shouldUseTunnel(nextHop) || (privateNextHop && (s.hasDirectPeer(nextHop.ID) || s.directStream != nil))
}

func (s *Server) shouldUseTunnel(nextHop domain.Node) bool {
	return s.tunnelRegistry != nil && s.tunnelRegistry.HasChild(nextHop.ID)
}

func (s *Server) hasDirectPeer(nodeID string) bool {
	available, ok := s.directStream.(directPeerAvailability)
	return ok && available.HasDirectPeer(nodeID)
}

func (s *Server) resolveChainCandidates(snapshot policystore.Snapshot, chainID string) ([]chainHop, bool) {
	var chain domain.Chain
	found := false
	for _, item := range snapshot.Chains {
		if item.ID == chainID {
			chain = item
			found = true
			break
		}
	}
	if !found || len(chain.HopGroups) == 0 {
		return nil, false
	}
	index := -1
	nodeID := s.nodeIDGetter()
	for i, group := range chain.HopGroups {
		for _, candidateID := range group.Candidates {
			if candidateID == nodeID {
				index = i
				break
			}
		}
		if index >= 0 {
			break
		}
	}
	if index == -1 {
		return nil, false
	}
	if index == len(chain.HopGroups)-1 {
		return []chainHop{{isLast: true}}, true
	}
	nodes := make(map[string]domain.Node, len(snapshot.Nodes))
	for _, node := range snapshot.Nodes {
		if node.Enabled {
			nodes[node.ID] = node
		}
	}
	candidates := make([]chainHop, 0, len(chain.HopGroups[index+1].Candidates))
	for _, candidateID := range chain.HopGroups[index+1].Candidates {
		node, exists := nodes[candidateID]
		if !exists || !chainLinkExists(snapshot.Links, nodeID, candidateID) || !s.candidateUsable(node) {
			continue
		}
		remaining, connected := resolveRemainingChainPath(chain.HopGroups, snapshot.Links, nodes, index+1, candidateID)
		if connected {
			candidates = append(candidates, chainHop{node: node, remainingHops: remaining})
		}
	}
	return candidates, len(candidates) > 0
}

func (s *Server) candidateUsable(candidate domain.Node) bool {
	return s.shouldUseTunnel(candidate) || s.hasDirectPeer(candidate.ID) || (candidate.PublicHost != "" && candidate.PublicPort > 0)
}

func resolveRemainingChainPath(groups []domain.ChainHopGroup, links []domain.NodeLink, nodes map[string]domain.Node, groupIndex int, currentNodeID string) ([]string, bool) {
	if groupIndex == len(groups)-1 {
		return nil, true
	}
	for _, candidateID := range groups[groupIndex+1].Candidates {
		if _, exists := nodes[candidateID]; !exists || !chainLinkExists(links, currentNodeID, candidateID) {
			continue
		}
		tail, connected := resolveRemainingChainPath(groups, links, nodes, groupIndex+1, candidateID)
		if connected {
			return append([]string{candidateID}, tail...), true
		}
	}
	return nil, false
}

func chainLinkExists(links []domain.NodeLink, sourceNodeID string, targetNodeID string) bool {
	for _, link := range links {
		if link.SourceNodeID == sourceNodeID && link.TargetNodeID == targetNodeID {
			return true
		}
	}
	return false
}
