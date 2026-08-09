package proxy

type Chain struct {
	ID               string          `json:"id"`
	CreateID         string          `json:"createId"`
	OwnerID          string          `json:"ownerId"`
	Name             string          `json:"name"`
	DestinationScope string          `json:"destinationScope"`
	Enabled          bool            `json:"enabled"`
	HopGroups        []ChainHopGroup `json:"hopGroups"`
	Permission       string          `json:"permission,omitempty"`
}

type ChainWithDetails struct {
	ID               string                `json:"id"`
	CreateID         string                `json:"createId"`
	OwnerID          string                `json:"ownerId"`
	Name             string                `json:"name"`
	DestinationScope string                `json:"destinationScope"`
	Enabled          bool                  `json:"enabled"`
	HopGroups        []ChainHopGroup       `json:"hopGroups"`
	HopGroupDetails  []ChainHopGroupDetail `json:"hopGroupDetails"`
	Permission       string                `json:"permission,omitempty"`
}

type ChainHopGroup struct {
	Candidates []string `json:"candidates"`
}

type ChainHopGroupDetail struct {
	Candidates []ChainHopDetail `json:"candidates"`
}

type ChainHopDetail struct {
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`
	Mode     string `json:"mode"`
}

type CreateChainInput struct {
	Name             string          `json:"name"`
	DestinationScope string          `json:"destinationScope"`
	HopGroups        []ChainHopGroup `json:"hopGroups"`
}

type UpdateChainInput struct {
	Name             string          `json:"name"`
	DestinationScope string          `json:"destinationScope"`
	HopGroups        []ChainHopGroup `json:"hopGroups"`
	Enabled          bool            `json:"enabled"`
}

type ValidateChainInput struct {
	Name             string          `json:"name"`
	HopGroups        []ChainHopGroup `json:"hopGroups"`
	DestinationScope string          `json:"destinationScope"`
}

type HopConnectivity struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Reachable bool   `json:"reachable"`
}

type ScopeOwnership struct {
	Scope       string `json:"scope"`
	OwnerNodeID string `json:"ownerNodeId"`
	Valid       bool   `json:"valid"`
}

type ChainValidationResult struct {
	Valid           bool              `json:"valid"`
	Errors          []string          `json:"errors"`
	Warnings        []string          `json:"warnings"`
	HopConnectivity []HopConnectivity `json:"hopConnectivity"`
	ScopeOwnership  ScopeOwnership    `json:"scopeOwnership"`
}

type PreviewChainInput struct {
	Name             string          `json:"name"`
	HopGroups        []ChainHopGroup `json:"hopGroups"`
	DestinationScope string          `json:"destinationScope"`
}

type CompiledChainConfig struct {
	ChainID          string                `json:"chainId"`
	Name             string                `json:"name"`
	HopGroups        []ChainHopGroupDetail `json:"hopGroups"`
	DestinationScope string                `json:"destinationScope"`
	RoutingPath      string                `json:"routingPath"`
}

type ChainPreviewResult struct {
	CompiledConfig CompiledChainConfig `json:"compiledConfig"`
}

type ChainProbeHop struct {
	NodeID        string `json:"nodeId"`
	NodeName      string `json:"nodeName"`
	TransportType string `json:"transportType"`
	Address       string `json:"address"`
	Status        string `json:"status"`
}

type ChainProbeResult struct {
	ChainID            string          `json:"chainId"`
	Status             string          `json:"status"`
	Message            string          `json:"message"`
	ResolvedHops       []ChainProbeHop `json:"resolvedHops"`
	BlockingGroupIndex int             `json:"blockingGroupIndex"`
	BlockingNodeID     string          `json:"blockingNodeId"`
	BlockingReason     string          `json:"blockingReason"`
	TargetHost         string          `json:"targetHost"`
	TargetPort         int             `json:"targetPort"`
	ProbedAt           string          `json:"probedAt"`
}

type SaveChainProbeResultInput struct {
	ChainID            string          `json:"chainId"`
	Status             string          `json:"status"`
	Message            string          `json:"message"`
	ResolvedHops       []ChainProbeHop `json:"resolvedHops"`
	BlockingGroupIndex int             `json:"blockingGroupIndex"`
	BlockingNodeID     string          `json:"blockingNodeId"`
	BlockingReason     string          `json:"blockingReason"`
	TargetHost         string          `json:"targetHost"`
	TargetPort         int             `json:"targetPort"`
	ProbedAt           string          `json:"probedAt"`
}
