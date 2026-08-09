package proxy

type TopologyProjection struct {
	GeneratedAt string         `json:"generatedAt"`
	Paths       []TopologyPath `json:"paths"`
}

type TopologyPath struct {
	AccessPathID   string          `json:"accessPathId"`
	AccessPathName string          `json:"accessPathName"`
	ChainID        string          `json:"chainId"`
	ChainName      string          `json:"chainName"`
	TargetScope    string          `json:"targetScope"`
	Enabled        bool            `json:"enabled"`
	Status         string          `json:"status"`
	BlockingReason string          `json:"blockingReason"`
	Groups         []TopologyGroup `json:"groups"`
	Edges          []TopologyEdge  `json:"edges"`
}

type TopologyGroup struct {
	Index      int                 `json:"index"`
	Candidates []TopologyCandidate `json:"candidates"`
}

type TopologyCandidate struct {
	NodeID     string `json:"nodeId"`
	NodeName   string `json:"nodeName"`
	Mode       string `json:"mode"`
	ScopeKey   string `json:"scopeKey"`
	PublicHost string `json:"publicHost,omitempty"`
	PublicPort int    `json:"publicPort,omitempty"`
	Priority   int    `json:"priority"`
	Role       string `json:"role"`
	Health     string `json:"health"`
	Selected   bool   `json:"selected"`
	Inbound    int    `json:"inbound"`
	Outbound   int    `json:"outbound"`
}

type TopologyEdge struct {
	SourceNodeID string `json:"sourceNodeId"`
	TargetNodeID string `json:"targetNodeId"`
	Kind         string `json:"kind"`
	Active       bool   `json:"active"`
	Status       string `json:"status"`
}
