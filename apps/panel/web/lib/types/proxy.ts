import type {ActionType, MatchType, NodeMode, ProbeResultStatus, TransportStatus, TransportType} from './common';
import type {ResourcePermissionMetadata} from './grants';

export type Chain = ResourcePermissionMetadata & {
  id: string;
  name: string;
  destinationScope: string;
  enabled: boolean;
  hopGroups: ChainHopGroup[];
};

export type ChainHopGroup = {
  candidates: string[];
};

export type DeleteImpactItem = {
  id: string;
  name: string;
  detail?: string;
};

export type RouteRuleGroup = ResourcePermissionMetadata & {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  createId: string;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
  ruleCount: number;
};

export type RouteRuleGroupDeleteImpact = {
  groupId: string;
  delete: {
    group: DeleteImpactItem[];
    routeRules: DeleteImpactItem[];
    tenantBindings: DeleteImpactItem[];
  };
};

export type ChainDeleteImpact = {
  chainId: string;
  delete: {
    chain: DeleteImpactItem[];
    chainHops: DeleteImpactItem[];
    routeRules: DeleteImpactItem[];
    accessPaths: DeleteImpactItem[];
    onboardingTasks: DeleteImpactItem[];
    chainProbeResults: DeleteImpactItem[];
    tenantBindings: DeleteImpactItem[];
  };
};

export type ChainProbeHop = {
  nodeId: string;
  nodeName: string;
  transportType: TransportType;
  address: string;
  status: TransportStatus;
};

export type ChainProbeResult = {
  chainId: string;
  status: ProbeResultStatus;
  message: string;
  resolvedHops: ChainProbeHop[];
  blockingGroupIndex: number;
  blockingNodeId: string;
  blockingReason: string;
  targetHost: string;
  targetPort: number;
  probedAt: string;
};

export type ChainValidationResult = {
  valid: boolean;
  errors: string[];
  warnings: string[];
  hopConnectivity: { from: string; to: string; reachable: boolean }[];
  scopeOwnership: { scope: string; ownerNodeIds: string[]; valid: boolean };
};

export type CompiledChainHop = {
  nodeId: string;
  nodeName: string;
  mode: NodeMode;
};

export type CompiledChainConfig = {
  chainId: string;
  name: string;
  hopGroups: {candidates: CompiledChainHop[]}[];
  destinationScope: string;
  routingPath: string;
};

export type ChainPreviewResult = {
  compiledConfig: CompiledChainConfig;
};

export type RouteRule = ResourcePermissionMetadata & {
  id: string;
  groupId: string;
  priority: number;
  matchType: MatchType;
  matchValue: string;
  actionType: ActionType;
  chainId?: string;
  destinationScope?: string;
  enabled: boolean;
};

export type MatchValueValidation = {
  valid: boolean;
  format: string;
  message: string;
};

export type ChainValidation = {
  valid: boolean;
  chainEnabled: boolean;
  chainHops: string[];
};

export type ScopeValidation = {
  valid: boolean;
  scopeExists: boolean;
  scopeOwnerNodeId: string;
  matchesChainFinalHop: boolean;
};

export type RouteRuleValidationResult = {
  valid: boolean;
  errors: string[];
  warnings: string[];
  matchValueValidation: MatchValueValidation;
  chainValidation: ChainValidation;
  scopeValidation: ScopeValidation;
};

export type Scope = ResourcePermissionMetadata & {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

export type TopologyCandidate = {
  nodeId: string;
  nodeName: string;
  mode: string;
  scopeKey: string;
  publicHost?: string;
  publicPort?: number;
  priority: number;
  role: 'primary' | 'standby';
  health: string;
  selected: boolean;
  inbound: number;
  outbound: number;
};

export type TopologyEdge = {
  sourceNodeId: string;
  targetNodeId: string;
  kind: 'primary' | 'standby';
  active: boolean;
  status: string;
};

export type TopologyPath = {
  accessPathId: string;
  accessPathName: string;
  chainId: string;
  chainName: string;
  targetScope: string;
  enabled: boolean;
  status: string;
  blockingReason: string;
  groups: {index: number; candidates: TopologyCandidate[]}[];
  edges: TopologyEdge[];
};

export type TopologyProjection = {
  generatedAt: string;
  paths: TopologyPath[];
};
