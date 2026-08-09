package proxyservice

import (
	"time"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func (s *Service) Topology(tenantCtx domain.TenantAuthContext) proxy.TopologyProjection {
	paths := s.store.ListNodeAccessPathsForTenant(tenantCtx)
	chains := s.store.ListChainsForTenant(tenantCtx)
	nodes := s.store.ListNodesForTenant(tenantCtx)
	links := s.store.ListNodeLinksForTenant(tenantCtx)
	transports := s.store.ListNodeTransports()
	probes := make(map[string]proxy.ChainProbeResult, len(chains))
	for _, chain := range chains {
		if result, ok := s.store.GetChainProbeResult(chain.ID); ok {
			probes[chain.ID] = result
		}
	}
	projection := buildTopologyProjection(paths, chains, nodes, links, transports, probes)
	projection.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	return projection
}

func buildTopologyProjection(paths []domain.NodeAccessPath, chains []proxy.Chain, nodes []domain.Node, links []domain.NodeLink, transports []domain.NodeTransport, probes map[string]proxy.ChainProbeResult) proxy.TopologyProjection {
	chainsByID := make(map[string]proxy.Chain, len(chains))
	for _, chain := range chains {
		chainsByID[chain.ID] = chain
	}
	nodesByID := make(map[string]domain.Node, len(nodes))
	for _, node := range nodes {
		nodesByID[node.ID] = node
	}
	projection := proxy.TopologyProjection{Paths: make([]proxy.TopologyPath, 0, len(paths))}
	for _, path := range paths {
		chain, ok := chainsByID[path.ChainID]
		if !ok {
			continue
		}
		probe := probes[chain.ID]
		selectedPairs := make(map[string]bool)
		selectedNodes := make(map[string]bool)
		for index, hop := range probe.ResolvedHops {
			selectedNodes[hop.NodeID] = true
			if index > 0 {
				selectedPairs[probe.ResolvedHops[index-1].NodeID+"\x00"+hop.NodeID] = true
			}
		}
		edges := topologyEdges(chain.HopGroups, links, selectedPairs, nodesByID, transports)
		inbound := make(map[string]int)
		outbound := make(map[string]int)
		for _, edge := range edges {
			outbound[edge.SourceNodeID]++
			inbound[edge.TargetNodeID]++
		}
		groups := make([]proxy.TopologyGroup, 0, len(chain.HopGroups))
		for groupIndex, group := range chain.HopGroups {
			candidates := make([]proxy.TopologyCandidate, 0, len(group.Candidates))
			for priority, nodeID := range group.Candidates {
				node := nodesByID[nodeID]
				role := "standby"
				if priority == 0 {
					role = "primary"
				}
				candidates = append(candidates, proxy.TopologyCandidate{
					NodeID: node.ID, NodeName: node.Name, Mode: node.Mode, ScopeKey: node.ScopeKey,
					PublicHost: node.PublicHost, PublicPort: node.PublicPort, Priority: priority,
					Role: role, Health: topologyNodeHealth(node, transports), Selected: selectedNodes[nodeID],
					Inbound: inbound[nodeID], Outbound: outbound[nodeID],
				})
			}
			groups = append(groups, proxy.TopologyGroup{Index: groupIndex, Candidates: candidates})
		}
		status := probe.Status
		if status == "" {
			status = "unknown"
		}
		blockingReason := probe.BlockingReason
		if blockingReason == "" && status != "connected" {
			blockingReason = probe.Message
		}
		projection.Paths = append(projection.Paths, proxy.TopologyPath{
			AccessPathID: path.ID, AccessPathName: path.Name, ChainID: chain.ID, ChainName: chain.Name,
			TargetScope: chain.DestinationScope, Enabled: path.Enabled && chain.Enabled, Status: status,
			BlockingReason: blockingReason, Groups: groups, Edges: edges,
		})
	}
	return projection
}

func topologyEdges(groups []proxy.ChainHopGroup, links []domain.NodeLink, selectedPairs map[string]bool, nodes map[string]domain.Node, transports []domain.NodeTransport) []proxy.TopologyEdge {
	priority := make(map[string]int)
	groupIndex := make(map[string]int)
	for index, group := range groups {
		for candidateIndex, nodeID := range group.Candidates {
			priority[nodeID] = candidateIndex
			groupIndex[nodeID] = index
		}
	}
	edges := make([]proxy.TopologyEdge, 0)
	for _, link := range links {
		if groupIndex[link.TargetNodeID] != groupIndex[link.SourceNodeID]+1 {
			continue
		}
		kind := "standby"
		if priority[link.SourceNodeID] == 0 && priority[link.TargetNodeID] == 0 {
			kind = "primary"
		}
		status := "available"
		if topologyNodeHealth(nodes[link.SourceNodeID], transports) == "down" || topologyNodeHealth(nodes[link.TargetNodeID], transports) == "down" {
			status = "blocked"
		}
		edges = append(edges, proxy.TopologyEdge{
			SourceNodeID: link.SourceNodeID, TargetNodeID: link.TargetNodeID, Kind: kind,
			Active: selectedPairs[link.SourceNodeID+"\x00"+link.TargetNodeID], Status: status,
		})
	}
	return edges
}

func topologyNodeHealth(node domain.Node, transports []domain.NodeTransport) string {
	if node.ID == "" || !node.Enabled {
		return "down"
	}
	health := node.Status
	for _, transport := range transports {
		if transport.NodeID != node.ID {
			continue
		}
		if transport.Status == domain.TransportStatusConnected || transport.Status == domain.TransportStatusAvailable {
			return "healthy"
		}
		if transport.Status == domain.TransportStatusFailed {
			health = "down"
		}
	}
	if health == domain.NodeStatusHealthy {
		return "healthy"
	}
	if health == "" {
		return "unknown"
	}
	return health
}
