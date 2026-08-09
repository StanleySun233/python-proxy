'use client';

import 'reactflow/dist/style.css';

import {Activity, ArrowDownToLine, ArrowUpFromLine, CircleDot, Route} from 'lucide-react';
import {useTranslations} from 'next-intl';
import {useMemo, useState} from 'react';
import {Background, Controls, MarkerType, ReactFlow} from 'reactflow';

import type {TopologyPath, TopologyProjection} from '@/lib/types';

function candidateNodeID(nodeID: string) {
  return `candidate:${nodeID}`;
}

function OperationalNode({candidate}: {candidate: TopologyPath['groups'][number]['candidates'][number]}) {
  return (
    <div className={`operational-node is-${candidate.health}${candidate.selected ? ' is-selected' : ''}`}>
      <div className="operational-node-head">
        <span className="operational-health" />
        <strong>{candidate.nodeName || candidate.nodeId}</strong>
        <span className={`topology-role is-${candidate.role}`}>{candidate.role}</span>
      </div>
      <span>{candidate.mode} · {candidate.scopeKey || 'global'}</span>
      <div className="operational-dependencies">
        <span title="Inbound dependencies"><ArrowDownToLine size={12} />{candidate.inbound}</span>
        <span title="Outbound dependencies"><ArrowUpFromLine size={12} />{candidate.outbound}</span>
      </div>
    </div>
  );
}

function topologyFlow(path: TopologyPath) {
  const laneWidth = 265;
  const rowHeight = 108;
  const maxCandidates = Math.max(1, ...path.groups.map((group) => group.candidates.length));
  const centerY = Math.max(55, ((maxCandidates - 1) * rowHeight) / 2 + 35);
  const nodes = [
    {
      id: 'client',
      position: {x: 15, y: centerY},
      type: 'input',
      data: {label: <div className="topology-terminal"><CircleDot size={15} /><strong>CLIENT</strong><span>request ingress</span></div>},
      className: 'topology-terminal-node'
    },
    ...path.groups.flatMap((group) => group.candidates.map((candidate, candidateIndex) => ({
      id: candidateNodeID(candidate.nodeId),
      position: {x: 220 + group.index * laneWidth, y: 18 + candidateIndex * rowHeight},
      data: {label: <OperationalNode candidate={candidate} />},
      className: 'operational-flow-node'
    }))),
    {
      id: 'scope',
      position: {x: 220 + path.groups.length * laneWidth, y: centerY},
      type: 'output',
      data: {label: <div className="topology-terminal"><Route size={15} /><strong>{path.targetScope || 'TARGET'}</strong><span>destination scope</span></div>},
      className: 'topology-terminal-node'
    }
  ];
  const firstGroup = path.groups[0];
  const lastGroup = path.groups[path.groups.length - 1];
  const terminalEdges = [
    ...(firstGroup?.candidates || []).map((candidate) => ({
      id: `client:${candidate.nodeId}`,
      source: 'client',
      target: candidateNodeID(candidate.nodeId),
      kind: candidate.role,
      active: candidate.selected,
      status: candidate.health === 'down' ? 'blocked' : 'available'
    })),
    ...(lastGroup?.candidates || []).map((candidate) => ({
      id: `${candidate.nodeId}:scope`,
      source: candidateNodeID(candidate.nodeId),
      target: 'scope',
      kind: candidate.role,
      active: candidate.selected,
      status: candidate.health === 'down' ? 'blocked' : 'available'
    }))
  ];
  const configuredEdges = path.edges.map((edge) => ({
    id: `${edge.sourceNodeId}:${edge.targetNodeId}`,
    source: candidateNodeID(edge.sourceNodeId),
    target: candidateNodeID(edge.targetNodeId),
    kind: edge.kind,
    active: edge.active,
    status: edge.status
  }));
  const edges = [...terminalEdges, ...configuredEdges].map((edge) => ({
    id: edge.id,
    source: edge.source,
    target: edge.target,
    animated: edge.active,
    markerEnd: {type: MarkerType.ArrowClosed, color: edge.status === 'blocked' ? 'var(--danger)' : edge.active ? 'var(--success)' : 'var(--muted)'},
    className: `topology-flow-edge is-${edge.kind}${edge.active ? ' is-active' : ''}${edge.status === 'blocked' ? ' is-blocked' : ''}`,
    style: {strokeDasharray: edge.kind === 'standby' ? '7 6' : undefined}
  }));
  return {nodes, edges};
}

export function TopologyPreview({projection, compact = false}: {projection?: TopologyProjection; compact?: boolean}) {
  const t = useTranslations();
  const [selectedPathID, setSelectedPathID] = useState('');
  const paths = projection?.paths || [];
  const selectedPath = paths.find((path) => path.accessPathId === selectedPathID) || paths[0];
  const flow = useMemo(() => selectedPath ? topologyFlow(selectedPath) : {nodes: [], edges: []}, [selectedPath]);

  if (!selectedPath) {
    return <div className="topology-empty">{t('overview.noOperationalPaths')}</div>;
  }

  return (
    <div className={`operational-topology${compact ? ' is-compact' : ''}`}>
      <div className="topology-command-bar">
        <div>
          <span className={`topology-path-status is-${selectedPath.status}`}><Activity size={13} />{selectedPath.status}</span>
          <strong>{selectedPath.accessPathName}</strong>
          <span>{selectedPath.chainName} · {selectedPath.groups.length} {t('overview.chainPositions')}</span>
        </div>
        {paths.length > 1 ? (
          <select className="field-select" onChange={(event) => setSelectedPathID(event.target.value)} value={selectedPath.accessPathId}>
            {paths.map((path) => <option key={path.accessPathId} value={path.accessPathId}>{path.accessPathName}</option>)}
          </select>
        ) : null}
      </div>
      {selectedPath.blockingReason ? <div className="topology-blocking-banner">{selectedPath.blockingReason}</div> : null}
      <div className="flow-card operational-flow-card">
        <ReactFlow
          edges={flow.edges}
          elementsSelectable={false}
          fitView
          nodes={flow.nodes}
          nodesConnectable={false}
          nodesDraggable={false}
          proOptions={{hideAttribution: true}}
        >
          <Controls showInteractive={false} />
          <Background gap={24} size={1} />
        </ReactFlow>
      </div>
      <div className="topology-legend">
        <span><i className="is-active" />{t('overview.activeDependency')}</span>
        <span><i />{t('overview.primaryDependency')}</span>
        <span><i className="is-standby" />{t('overview.standbyDependency')}</span>
      </div>
    </div>
  );
}
