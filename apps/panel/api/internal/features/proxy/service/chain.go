package proxyservice

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/controlrelay"
	"github.com/StanleySun233/python-proxy/apps/panel/api/internal/domain"
	proxy "github.com/StanleySun233/python-proxy/apps/panel/api/internal/features/proxy/domain"
)

func (s *Service) Chains(tenantCtx domain.TenantAuthContext) []proxy.Chain {
	return s.store.ListChainsForTenant(tenantCtx)
}

func (s *Service) ChainsWithDetails(tenantCtx domain.TenantAuthContext) []proxy.ChainWithDetails {
	chains := s.store.ListChainsForTenant(tenantCtx)
	nodes := s.store.ListNodesForTenant(tenantCtx)
	result := make([]proxy.ChainWithDetails, 0, len(chains))

	for _, chain := range chains {
		result = append(result, proxy.ChainWithDetails{
			ID:               chain.ID,
			CreateID:         chain.CreateID,
			OwnerID:          chain.OwnerID,
			Name:             chain.Name,
			DestinationScope: chain.DestinationScope,
			Enabled:          chain.Enabled,
			HopGroups:        chain.HopGroups,
			HopGroupDetails:  chainHopGroupDetails(chain.HopGroups, nodes),
			Permission:       chain.Permission,
		})
	}

	return result
}

func (s *Service) GetChain(tenantCtx domain.TenantAuthContext, chainID string) (proxy.ChainWithDetails, error) {
	if chainID == "" {
		return proxy.ChainWithDetails{}, invalidInput("missing_chain_id")
	}

	chains := s.store.ListChainsForTenant(tenantCtx)
	chain, ok := chainByID(chains, chainID)
	if !ok {
		return proxy.ChainWithDetails{}, invalidInput("chain_not_found")
	}

	nodes := s.store.ListNodesForTenant(tenantCtx)
	return proxy.ChainWithDetails{
		ID:               chain.ID,
		CreateID:         chain.CreateID,
		OwnerID:          chain.OwnerID,
		Name:             chain.Name,
		DestinationScope: chain.DestinationScope,
		Enabled:          chain.Enabled,
		HopGroups:        chain.HopGroups,
		HopGroupDetails:  chainHopGroupDetails(chain.HopGroups, nodes),
		Permission:       chain.Permission,
	}, nil
}

func (s *Service) LatestChainProbe(tenantCtx domain.TenantAuthContext, chainID string) (proxy.ChainProbeResult, bool) {
	if chainID == "" {
		return proxy.ChainProbeResult{}, false
	}
	if _, ok := s.store.ChainBindingPermission(tenantCtx, chainID); !ok {
		return proxy.ChainProbeResult{}, false
	}
	return s.store.GetChainProbeResult(chainID)
}

func (s *Service) ProbeChain(tenantCtx domain.TenantAuthContext, chainID string) (proxy.ChainProbeResult, error) {
	if chainID == "" {
		return proxy.ChainProbeResult{}, invalidInput("missing_chain_id")
	}
	chain, ok := chainByID(s.store.ListChainsForTenant(tenantCtx), chainID)
	if !ok {
		return proxy.ChainProbeResult{}, invalidInput("invalid_chain_id")
	}
	nodes := s.store.ListNodesForTenant(tenantCtx)
	links := s.store.ListNodeLinksForTenant(tenantCtx)
	transports := s.store.ListNodeTransports()
	result := proxy.ChainProbeResult{
		ChainID:            chainID,
		Status:             domain.ProbeResultStatusConnected,
		Message:            "chain_transport_ready",
		ResolvedHops:       make([]proxy.ChainProbeHop, 0, len(chain.HopGroups)),
		BlockingGroupIndex: -1,
		ProbedAt:           time.Now().UTC().Format(time.RFC3339),
	}
	resolvedHops, blockingGroupIndex, blockingReason := resolveProbeCandidatePath(chain.HopGroups, nodes, links, transports)
	result.ResolvedHops = resolvedHops
	if blockingGroupIndex >= 0 {
		result.Status = domain.ProbeResultStatusFailed
		result.Message = "chain_blocked"
		result.BlockingGroupIndex = blockingGroupIndex
		result.BlockingReason = blockingReason
		return s.store.SaveChainProbeResult(toChainProbeInput(result))
	}
	resolvedNodeIDs := make([]string, 0, len(result.ResolvedHops))
	for _, hop := range result.ResolvedHops {
		resolvedNodeIDs = append(resolvedNodeIDs, hop.NodeID)
	}
	if len(result.ResolvedHops) > 0 && (result.ResolvedHops[0].TransportType == domain.TransportTypePublicHTTP || result.ResolvedHops[0].TransportType == domain.TransportTypePublicHTTPS) {
		probeResult, err := controlrelay.Execute(result.ResolvedHops[0].Address, controlrelay.ProbeRequest{
			RemainingHopNodeIDs: resolvedNodeIDs[1:],
		})
		if err != nil {
			result.Status = domain.ProbeResultStatusFailed
			result.Message = "chain_probe_failed"
			result.BlockingGroupIndex = 0
			result.BlockingNodeID = resolvedNodeIDs[0]
			result.BlockingReason = "probe_dispatch_failed"
			return s.store.SaveChainProbeResult(toChainProbeInput(result))
		}
		result.Status = probeResult.Status
		result.Message = probeResult.Message
		if probeResult.Status != domain.ProbeResultStatusConnected && result.BlockingReason == "" && len(resolvedNodeIDs) > 0 {
			result.BlockingGroupIndex = len(resolvedNodeIDs) - 1
			result.BlockingNodeID = resolvedNodeIDs[len(resolvedNodeIDs)-1]
			result.BlockingReason = probeResult.Message
		}
	}
	return s.store.SaveChainProbeResult(toChainProbeInput(result))
}

func (s *Service) CreateChain(tenantCtx domain.TenantAuthContext, input proxy.CreateChainInput) (proxy.Chain, error) {
	if err := requireActiveTenant(tenantCtx); err != nil {
		return proxy.Chain{}, err
	}
	if !tenantCtx.SuperAdmin && tenantCtx.ActiveTenant.Role != domain.TenantRoleAdmin {
		return proxy.Chain{}, newError(http.StatusForbidden, "tenant_role_forbidden")
	}
	if input.Name == "" || input.DestinationScope == "" || len(input.HopGroups) == 0 {
		return proxy.Chain{}, invalidInput("invalid_chain_payload")
	}
	if !s.tenantScopeExists(tenantCtx, input.DestinationScope) {
		return proxy.Chain{}, invalidInput("scope_not_found")
	}
	if !s.tenantNodesExist(tenantCtx, chainNodeIDs(input.HopGroups)) {
		return proxy.Chain{}, invalidInput("node_not_found")
	}
	if len(validateChainGraph(input.HopGroups, s.store.ListNodesForTenant(tenantCtx), s.store.ListNodeLinksForTenant(tenantCtx), input.DestinationScope)) > 0 {
		return proxy.Chain{}, invalidInput("invalid_chain_graph")
	}
	return s.store.CreateChainForTenant(tenantCtx, input)
}

func (s *Service) UpdateChain(tenantCtx domain.TenantAuthContext, chainID string, input proxy.UpdateChainInput) (proxy.Chain, error) {
	if chainID == "" || input.Name == "" || input.DestinationScope == "" || len(input.HopGroups) == 0 {
		return proxy.Chain{}, invalidInput("invalid_chain_payload")
	}
	if err := s.requireTenantResourceManage(tenantCtx, func() (domain.BindingPermission, bool) {
		return s.store.ChainBindingPermission(tenantCtx, chainID)
	}); err != nil {
		return proxy.Chain{}, err
	}
	if !s.tenantScopeExists(tenantCtx, input.DestinationScope) {
		return proxy.Chain{}, invalidInput("scope_not_found")
	}
	if !s.tenantNodesExist(tenantCtx, chainNodeIDs(input.HopGroups)) {
		return proxy.Chain{}, invalidInput("node_not_found")
	}
	if len(validateChainGraph(input.HopGroups, s.store.ListNodesForTenant(tenantCtx), s.store.ListNodeLinksForTenant(tenantCtx), input.DestinationScope)) > 0 {
		return proxy.Chain{}, invalidInput("invalid_chain_graph")
	}
	return s.store.UpdateChain(chainID, input)
}

func (s *Service) DeleteChain(tenantCtx domain.TenantAuthContext, chainID string) error {
	if err := s.requireTenantResourceManage(tenantCtx, func() (domain.BindingPermission, bool) {
		return s.store.ChainBindingPermission(tenantCtx, chainID)
	}); err != nil {
		return err
	}
	if !tenantCtx.SuperAdmin && s.store.CountChainBindings(chainID) > 1 {
		return newError(http.StatusConflict, "shared_resource_delete_forbidden")
	}
	return s.store.DeleteChain(chainID)
}

func (s *Service) ChainDeleteImpact(tenantCtx domain.TenantAuthContext, chainID string) (proxy.ChainDeleteImpact, error) {
	if chainID == "" {
		return proxy.ChainDeleteImpact{}, invalidInput("missing_chain_id")
	}
	if err := s.requireTenantResourceManage(tenantCtx, func() (domain.BindingPermission, bool) {
		return s.store.ChainBindingPermission(tenantCtx, chainID)
	}); err != nil {
		return proxy.ChainDeleteImpact{}, err
	}
	if !tenantCtx.SuperAdmin && s.store.CountChainBindings(chainID) > 1 {
		return proxy.ChainDeleteImpact{}, newError(http.StatusConflict, "shared_resource_delete_forbidden")
	}
	return s.store.GetChainDeleteImpact(chainID)
}

func (s *Service) ValidateChain(tenantCtx domain.TenantAuthContext, input proxy.ValidateChainInput) (proxy.ChainValidationResult, error) {
	nodes := s.store.ListNodesForTenant(tenantCtx)
	links := s.store.ListNodeLinksForTenant(tenantCtx)
	errors := validateChainGraph(input.HopGroups, nodes, links, input.DestinationScope)
	connectivity := make([]proxy.HopConnectivity, 0)
	for groupIndex := 0; groupIndex+1 < len(input.HopGroups); groupIndex++ {
		for _, sourceID := range input.HopGroups[groupIndex].Candidates {
			for _, targetID := range input.HopGroups[groupIndex+1].Candidates {
				connectivity = append(connectivity, proxy.HopConnectivity{From: sourceID, To: targetID, Reachable: hasNodeLink(sourceID, targetID, links)})
			}
		}
	}
	ownerNodeIDs := []string{}
	if len(input.HopGroups) > 0 {
		ownerNodeIDs = append(ownerNodeIDs, input.HopGroups[len(input.HopGroups)-1].Candidates...)
	}
	result := proxy.ChainValidationResult{
		Valid:           len(errors) == 0,
		Errors:          errors,
		Warnings:        []string{},
		HopConnectivity: connectivity,
		ScopeOwnership:  proxy.ScopeOwnership{Scope: input.DestinationScope, OwnerNodeIDs: ownerNodeIDs, Valid: len(errors) == 0},
	}
	return result, nil
}

func (s *Service) PreviewChain(tenantCtx domain.TenantAuthContext, input proxy.PreviewChainInput) (proxy.ChainPreviewResult, error) {
	nodes := s.store.ListNodesForTenant(tenantCtx)
	hopGroupDetails := make([]proxy.ChainHopGroupDetail, 0, len(input.HopGroups))
	routingPath := "user"

	for _, group := range input.HopGroups {
		candidates := make([]proxy.ChainHopDetail, 0, len(group.Candidates))
		candidateNames := make([]string, 0, len(group.Candidates))
		for _, nodeID := range group.Candidates {
			node, ok := nodeByID(nodes, nodeID)
			if !ok {
				return proxy.ChainPreviewResult{}, invalidInput(fmt.Sprintf("node %s not found", nodeID))
			}
			candidates = append(candidates, proxy.ChainHopDetail{NodeID: node.ID, NodeName: node.Name, Mode: node.Mode})
			candidateNames = append(candidateNames, node.Name)
		}
		hopGroupDetails = append(hopGroupDetails, proxy.ChainHopGroupDetail{Candidates: candidates})
		routingPath += " → [" + strings.Join(candidateNames, " | ") + "]"
	}

	routingPath += fmt.Sprintf(" → target(%s)", input.DestinationScope)

	return proxy.ChainPreviewResult{
		CompiledConfig: proxy.CompiledChainConfig{
			ChainID:          "preview",
			Name:             input.Name,
			HopGroups:        hopGroupDetails,
			DestinationScope: input.DestinationScope,
			RoutingPath:      routingPath,
		},
	}, nil
}

func resolveProbeCandidatePath(groups []proxy.ChainHopGroup, nodes []domain.Node, links []domain.NodeLink, transports []domain.NodeTransport) ([]proxy.ChainProbeHop, int, string) {
	resolved := make([]proxy.ChainProbeHop, 0, len(groups))
	previousNodeID := ""
	for groupIndex, group := range groups {
		var selectedNode domain.Node
		var selectedTransport domain.NodeTransport
		for _, candidateID := range group.Candidates {
			node, ok := nodeByID(nodes, candidateID)
			if !ok || !node.Enabled || node.Status == domain.NodeStatusPending {
				continue
			}
			if previousNodeID != "" && !hasNodeLink(previousNodeID, candidateID, links) {
				continue
			}
			transport, ok := resolveProbeTransport(node, previousNodeID, transports)
			if !ok {
				continue
			}
			selectedNode = node
			selectedTransport = transport
			break
		}
		if selectedNode.ID == "" {
			if previousNodeID == "" {
				return resolved, groupIndex, "entry_candidates_unavailable"
			}
			return resolved, groupIndex, "next_hop_candidates_unavailable"
		}
		resolved = append(resolved, proxy.ChainProbeHop{
			NodeID:        selectedNode.ID,
			NodeName:      selectedNode.Name,
			TransportType: selectedTransport.TransportType,
			Address:       selectedTransport.Address,
			Status:        selectedTransport.Status,
		})
		previousNodeID = selectedNode.ID
	}
	return resolved, -1, ""
}

func resolveProbeTransport(node domain.Node, prevHopID string, transports []domain.NodeTransport) (domain.NodeTransport, bool) {
	if prevHopID != "" {
		for _, transport := range transports {
			if transport.NodeID != node.ID || transport.ParentNodeID != prevHopID {
				continue
			}
			if transport.Status != domain.TransportStatusConnected {
				continue
			}
			if strings.HasPrefix(transport.TransportType, domain.TransportTypeReverseWS) || strings.HasPrefix(transport.TransportType, domain.TransportTypeChildWS) {
				return transport, true
			}
		}
	}
	for _, transport := range transports {
		if transport.NodeID != node.ID {
			continue
		}
		if transport.Status != domain.TransportStatusAvailable && transport.Status != domain.TransportStatusConnected {
			continue
		}
		if transport.TransportType == domain.TransportTypePublicHTTPS || transport.TransportType == domain.TransportTypePublicHTTP {
			return transport, true
		}
	}
	return domain.NodeTransport{}, false
}

func toChainProbeInput(result proxy.ChainProbeResult) proxy.SaveChainProbeResultInput {
	return proxy.SaveChainProbeResultInput{
		ChainID:            result.ChainID,
		Status:             result.Status,
		Message:            result.Message,
		ResolvedHops:       result.ResolvedHops,
		BlockingGroupIndex: result.BlockingGroupIndex,
		BlockingNodeID:     result.BlockingNodeID,
		BlockingReason:     result.BlockingReason,
		TargetHost:         result.TargetHost,
		TargetPort:         result.TargetPort,
		ProbedAt:           result.ProbedAt,
	}
}

func chainByID(items []proxy.Chain, chainID string) (proxy.Chain, bool) {
	for _, item := range items {
		if item.ID == chainID {
			return item, true
		}
	}
	return proxy.Chain{}, false
}

func validateChainGraph(groups []proxy.ChainHopGroup, nodes []domain.Node, links []domain.NodeLink, destinationScope string) []string {
	errors := make([]string, 0)
	if len(groups) == 0 {
		return []string{"missing_hop_groups"}
	}
	nodesByID := make(map[string]domain.Node, len(nodes))
	for _, node := range nodes {
		nodesByID[node.ID] = node
	}
	seen := make(map[string]struct{})
	for groupIndex, group := range groups {
		if len(group.Candidates) == 0 {
			errors = append(errors, fmt.Sprintf("empty_hop_group:%d", groupIndex))
			continue
		}
		for _, candidateID := range group.Candidates {
			node, ok := nodesByID[candidateID]
			if !ok || !node.Enabled || node.Status == domain.NodeStatusPending {
				errors = append(errors, "unknown_or_disabled_candidate:"+candidateID)
				continue
			}
			if _, exists := seen[candidateID]; exists {
				errors = append(errors, "duplicate_candidate:"+candidateID)
			} else {
				seen[candidateID] = struct{}{}
			}
			if groupIndex == 0 && node.Mode != domain.NodeModeEdge {
				errors = append(errors, "entry_candidate_not_edge:"+candidateID)
			}
			if groupIndex == len(groups)-1 && node.ScopeKey != destinationScope {
				errors = append(errors, "destination_scope_mismatch:"+candidateID)
			}
		}
	}
	for groupIndex := 0; groupIndex+1 < len(groups); groupIndex++ {
		current := groups[groupIndex]
		next := groups[groupIndex+1]
		for _, sourceID := range current.Candidates {
			reachable := false
			for _, targetID := range next.Candidates {
				if hasNodeLink(sourceID, targetID, links) {
					reachable = true
					break
				}
			}
			if !reachable {
				errors = append(errors, fmt.Sprintf("isolated_outbound_candidate:%s:%d", sourceID, groupIndex))
			}
		}
		for _, targetID := range next.Candidates {
			reachable := false
			for _, sourceID := range current.Candidates {
				if hasNodeLink(sourceID, targetID, links) {
					reachable = true
					break
				}
			}
			if !reachable {
				errors = append(errors, fmt.Sprintf("isolated_inbound_candidate:%s:%d", targetID, groupIndex+1))
			}
		}
	}
	return errors
}

func hasNodeLink(sourceID string, targetID string, links []domain.NodeLink) bool {
	for _, link := range links {
		if link.SourceNodeID == sourceID && link.TargetNodeID == targetID {
			return true
		}
	}
	return false
}

func chainHopGroupDetails(groups []proxy.ChainHopGroup, nodes []domain.Node) []proxy.ChainHopGroupDetail {
	details := make([]proxy.ChainHopGroupDetail, 0, len(groups))
	for _, group := range groups {
		candidates := make([]proxy.ChainHopDetail, 0, len(group.Candidates))
		for _, nodeID := range group.Candidates {
			if node, ok := nodeByID(nodes, nodeID); ok {
				candidates = append(candidates, proxy.ChainHopDetail{NodeID: node.ID, NodeName: node.Name, Mode: node.Mode})
			}
		}
		details = append(details, proxy.ChainHopGroupDetail{Candidates: candidates})
	}
	return details
}

func chainNodeIDs(groups []proxy.ChainHopGroup) []string {
	nodeIDs := make([]string, 0)
	for _, group := range groups {
		nodeIDs = append(nodeIDs, group.Candidates...)
	}
	return nodeIDs
}
