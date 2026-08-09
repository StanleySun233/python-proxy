'use client';

import {DndContext, closestCenter, DragEndEvent, PointerSensor, useSensor, useSensors} from '@dnd-kit/core';
import {SortableContext, arrayMove, verticalListSortingStrategy, useSortable} from '@dnd-kit/sortable';
import {CSS} from '@dnd-kit/utilities';
import {ArrowDown, ArrowUp, GripVertical, Plus, ShieldCheck, X} from 'lucide-react';
import {useTranslations} from 'next-intl';
import {useCallback, useEffect, useMemo, useRef, useState} from 'react';

import {NameTag} from '@/components/common/name-tag';
import {validateChain} from '@/lib/api';
import type {ChainHopGroup, ChainValidationResult, Node, Scope} from '@/lib/types';

type ChainEditorProps = {
  accessToken: string;
  activeTenantId: string | null;
  chainName: string;
  destinationScope: string;
  hopGroups: ChainHopGroup[];
  nodes: Node[];
  scopes: Scope[];
  onNameChange: (name: string) => void;
  onScopeChange: (scope: string) => void;
  onHopGroupsChange: (groups: ChainHopGroup[]) => void;
  onSave: () => void;
  onCancel: () => void;
  onPreview: () => void;
  saving: boolean;
  previewing: boolean;
};

type SortableHopGroupProps = {
  group: ChainHopGroup;
  index: number;
  nodes: Node[];
  unavailableNodeIDs: Set<string>;
  onAddCandidate: (nodeID: string) => void;
  onMoveCandidate: (candidateIndex: number, direction: -1 | 1) => void;
  onRemoveCandidate: (candidateIndex: number) => void;
  onRemoveGroup: () => void;
};

function SortableHopGroup({
  group,
  index,
  nodes,
  unavailableNodeIDs,
  onAddCandidate,
  onMoveCandidate,
  onRemoveCandidate,
  onRemoveGroup
}: SortableHopGroupProps) {
  const chainsT = useTranslations('proxyChains');
  const t = useTranslations();
  const [selectedNodeID, setSelectedNodeID] = useState('');
  const {attributes, listeners, setNodeRef, transform, transition, isDragging} = useSortable({id: `group-${index}`});
  const nodeByID = useMemo(() => new Map(nodes.map((node) => [node.id, node])), [nodes]);
  const availableNodes = nodes.filter((node) => !unavailableNodeIDs.has(node.id));

  return (
    <article
      className={`hop-group-card${isDragging ? ' is-dragging' : ''}`}
      ref={setNodeRef}
      style={{transform: CSS.Transform.toString(transform), transition}}
    >
      <header className="hop-group-head">
        <button className="hop-group-grip" type="button" {...attributes} {...listeners}>
          <GripVertical size={17} />
        </button>
        <span className="hop-group-index">{String(index + 1).padStart(2, '0')}</span>
        <div>
          <strong>{chainsT('hopGroupTitle', {index: index + 1})}</strong>
          <span>{chainsT('candidateCount', {count: group.candidates.length})}</span>
        </div>
        <button className="hop-card-remove" onClick={onRemoveGroup} title={chainsT('removeGroup')} type="button">
          <X size={16} />
        </button>
      </header>

      <div className="candidate-stack">
        {group.candidates.map((nodeID, candidateIndex) => {
          const node = nodeByID.get(nodeID);
          return (
            <div className={`candidate-row${candidateIndex === 0 ? ' is-primary' : ''}`} key={nodeID}>
              <span className="candidate-rank">{candidateIndex === 0 ? <ShieldCheck size={14} /> : candidateIndex + 1}</span>
              <div className="candidate-identity">
                <NameTag kind="node">{node?.name || t('common.unknown')}</NameTag>
                <span>{candidateIndex === 0 ? chainsT('primaryCandidate') : chainsT('standbyCandidate')} · {node?.mode || 'unknown'}</span>
              </div>
              <div className="candidate-actions">
                <button disabled={candidateIndex === 0} onClick={() => onMoveCandidate(candidateIndex, -1)} title={chainsT('raisePriority')} type="button">
                  <ArrowUp size={14} />
                </button>
                <button disabled={candidateIndex === group.candidates.length - 1} onClick={() => onMoveCandidate(candidateIndex, 1)} title={chainsT('lowerPriority')} type="button">
                  <ArrowDown size={14} />
                </button>
                <button onClick={() => onRemoveCandidate(candidateIndex)} title={chainsT('removeCandidate')} type="button">
                  <X size={14} />
                </button>
              </div>
            </div>
          );
        })}
        {group.candidates.length === 0 ? <div className="candidate-empty">{chainsT('emptyGroup')}</div> : null}
      </div>

      <div className="candidate-add-row">
        <select className="field-select" onChange={(event) => setSelectedNodeID(event.target.value)} value={selectedNodeID}>
          <option value="">{chainsT('selectCandidate')}</option>
          {availableNodes.map((node) => <option key={node.id} value={node.id}>{node.name} · {node.mode}</option>)}
        </select>
        <button
          className="secondary-button"
          disabled={!selectedNodeID}
          onClick={() => {
            onAddCandidate(selectedNodeID);
            setSelectedNodeID('');
          }}
          type="button"
        >
          <Plus size={15} />
          {chainsT('addCandidate')}
        </button>
      </div>
    </article>
  );
}

export function ChainEditor({
  accessToken,
  activeTenantId,
  chainName,
  destinationScope,
  hopGroups,
  nodes,
  scopes,
  onNameChange,
  onScopeChange,
  onHopGroupsChange,
  onSave,
  onCancel,
  onPreview,
  saving,
  previewing
}: ChainEditorProps) {
  const t = useTranslations();
  const chainsT = useTranslations('proxyChains');
  const [validationResult, setValidationResult] = useState<ChainValidationResult | null>(null);
  const [validationPending, setValidationPending] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const sensors = useSensors(useSensor(PointerSensor, {activationConstraint: {distance: 8}}));
  const usedNodeIDs = useMemo(() => new Set(hopGroups.flatMap((group) => group.candidates)), [hopGroups]);

  const runValidation = useCallback(async (name: string, scope: string, groups: ChainHopGroup[]) => {
    if (!name.trim() || !scope.trim() || groups.length === 0) {
      setValidationResult(null);
      return;
    }
    setValidationPending(true);
    try {
      setValidationResult(await validateChain(accessToken, activeTenantId, {name, destinationScope: scope, hopGroups: groups}));
    } catch {
      setValidationResult(null);
    } finally {
      setValidationPending(false);
    }
  }, [accessToken, activeTenantId]);

  useEffect(() => {
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }
    debounceRef.current = setTimeout(() => void runValidation(chainName, destinationScope, hopGroups), 500);
    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
    };
  }, [chainName, destinationScope, hopGroups, runValidation]);

  const updateGroup = (index: number, candidates: string[]) => {
    onHopGroupsChange(hopGroups.map((group, groupIndex) => groupIndex === index ? {candidates} : group));
  };

  const handleDragEnd = ({active, over}: DragEndEvent) => {
    if (!over || active.id === over.id) {
      return;
    }
    const oldIndex = Number(String(active.id).replace('group-', ''));
    const newIndex = Number(String(over.id).replace('group-', ''));
    onHopGroupsChange(arrayMove(hopGroups, oldIndex, newIndex));
  };

  const canSave = chainName.trim() && destinationScope && hopGroups.length > 0 && hopGroups.every((group) => group.candidates.length > 0);

  return (
    <div className="chain-editor">
      <div className="chain-editor-banner">
        <div>
          <p className="section-kicker">{chainsT('chainEditor')}</p>
          <h3>{chainName || chainsT('newChain')}</h3>
          <p>{chainsT('candidateGroupHint')}</p>
        </div>
        <span className={`chain-validation-light${validationResult?.valid ? ' is-valid' : validationResult ? ' is-invalid' : ''}`}>
          {validationPending ? t('common.validating') : validationResult?.valid ? t('common.valid') : validationResult ? t('common.invalid') : chainsT('draft')}
        </span>
      </div>

      <div className="forms-grid">
        <label className="field-stack">
          <span>{chainsT('chainName')}</span>
          <input className="field-input" onChange={(event) => onNameChange(event.target.value)} placeholder={chainsT('chainNamePlaceholder')} value={chainName} />
        </label>
        <label className="field-stack">
          <span>{chainsT('destinationScope')}</span>
          <select className="field-select" onChange={(event) => onScopeChange(event.target.value)} value={destinationScope}>
            <option value="">{chainsT('destinationScopePlaceholder')}</option>
            {scopes.map((scope) => <option key={scope.id} value={scope.id}>{scope.name}</option>)}
          </select>
        </label>
      </div>

      <section className="hop-editor-section">
        <div className="section-header">
          <div>
            <h4>{chainsT('hopGroups')}</h4>
            <span className="muted-text">{chainsT('priorityHint')}</span>
          </div>
          <button className="secondary-button" onClick={() => onHopGroupsChange([...hopGroups, {candidates: []}])} type="button">
            <Plus size={16} />
            {chainsT('addGroup')}
          </button>
        </div>

        <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd} sensors={sensors}>
          <SortableContext items={hopGroups.map((_, index) => `group-${index}`)} strategy={verticalListSortingStrategy}>
            <div className="hop-group-list">
              {hopGroups.map((group, index) => (
                <SortableHopGroup
                  group={group}
                  index={index}
                  key={`group-${index}`}
                  nodes={nodes}
                  onAddCandidate={(nodeID) => updateGroup(index, [...group.candidates, nodeID])}
                  onMoveCandidate={(candidateIndex, direction) => updateGroup(index, arrayMove(group.candidates, candidateIndex, candidateIndex + direction))}
                  onRemoveCandidate={(candidateIndex) => updateGroup(index, group.candidates.filter((_, itemIndex) => itemIndex !== candidateIndex))}
                  onRemoveGroup={() => onHopGroupsChange(hopGroups.filter((_, groupIndex) => groupIndex !== index))}
                  unavailableNodeIDs={usedNodeIDs}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
        {hopGroups.length === 0 ? <div className="empty-hops">{chainsT('noGroups')}</div> : null}
      </section>

      {validationResult && (validationResult.errors.length > 0 || validationResult.warnings.length > 0) ? (
        <div className="probe-results-section">
          {validationResult.errors.map((message) => <div className="chain-validation-message is-error" key={message}>{message}</div>)}
          {validationResult.warnings.map((message) => <div className="chain-validation-message is-warning" key={message}>{message}</div>)}
        </div>
      ) : null}

      <div className="submit-row">
        <button className="primary-button" disabled={saving || !canSave} onClick={onSave} type="button">{saving ? t('common.saving') : chainsT('saveChain')}</button>
        <button className="secondary-button" disabled={previewing || !canSave} onClick={onPreview} type="button">{previewing ? t('common.compiling') : chainsT('preview')}</button>
        <button className="secondary-button" onClick={onCancel} type="button">{t('common.cancel')}</button>
      </div>
    </div>
  );
}
