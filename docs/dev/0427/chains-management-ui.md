# Chains Management UI

## Background

Current `/chains` page provides basic chain listing but lacks comprehensive management capabilities. Operators need visual chain editing, hop reordering via drag-and-drop, destination scope validation, and real-time compilation preview to ensure chains are correctly configured before policy publication.

## Goal

Enhance `/chains` page to support visual chain management with drag-and-drop hop editing, scope validation, and compilation preview, improving operator efficiency and reducing configuration errors.

## Functional Requirements

### FR1. Chain List View

- Display all chains in a table with columns:
  - ID
  - Name
  - Hops (visual representation: `1 → 2 → 3`)
  - Destination Scope
  - Enabled Status
  - Actions (Edit/Delete/Duplicate)
- Support filtering by enabled status
- Support search by chain name
- Display chain validation status (valid/invalid)

### FR2. Visual Chain Editor

- Modal or drawer-based chain editor
- Drag-and-drop interface for hop reordering
- Node selection from available nodes
- Destination scope input with validation
- Real-time validation feedback
- Preview compiled chain configuration

### FR3. Hop Management

- Add node to chain by selecting from dropdown
- Remove node from chain by clicking remove button
- Reorder nodes by dragging hop cards
- Visual feedback during drag operation
- Prevent duplicate nodes in same chain
- Validate node connectivity (each hop must be reachable from previous)

### FR4. Scope Validation

- Validate destination scope format
- Check scope exists in node scope registry
- Warn if scope is not owned by final hop node
- Display scope ownership information
- Suggest valid scopes based on final hop node

### FR5. Compilation Preview

- Show compiled chain configuration in JSON format
- Display routing path: `user → node1 → node2 → node3 → target`
- Show scope resolution: `target resolved in scope: k8s-prod`
- Highlight validation errors or warnings
- Allow copy compiled configuration to clipboard

### FR6. Chain Validation Rules

- Chain must have at least one hop
- First hop must be an edge node
- Each hop must be reachable from previous hop (based on node_links)
- Destination scope must be owned by final hop node
- Chain name must be unique
- Enabled chains must pass all validation rules

## User Interaction Flow

### Flow 1: Creating a New Chain

1. Operator navigates to `/chains`
2. Clicks "Create Chain" button
3. Chain editor drawer opens
4. Operator enters chain name: `prod-k8s-path`
5. Operator selects destination scope: `k8s-prod`
6. Operator adds first hop:
   - Clicks "Add Hop"
   - Selects node from dropdown: `1 - edge-node-a`
   - Node card appears in hop list
7. Operator adds second hop:
   - Clicks "Add Hop"
   - Selects node: `2 - relay-node-b`
   - Node card appears below first hop
8. Operator adds third hop:
   - Selects node: `3 - k8s-gateway-node`
9. System validates chain:
   - ✓ All hops are reachable
   - ✓ Destination scope `k8s-prod` is owned by node 3
   - ✓ Chain name is unique
10. Operator clicks "Preview Compilation"
11. System shows compiled configuration:
    ```json
    {
      "chain_id": "new",
      "name": "prod-k8s-path",
      "hops": [1, 2, 3],
      "destination_scope": "k8s-prod",
      "routing_path": "user → 1 → 2 → 3 → target(k8s-prod)"
    }
    ```
12. Operator clicks "Save"
13. System creates chain and displays success message
14. Chain appears in chain list

### Flow 2: Editing Existing Chain

1. Operator clicks "Edit" on existing chain
2. Chain editor opens with current configuration
3. Operator drags hop 2 to position 3
4. System reorders hops: `1 → 3 → 2`
5. System validates new order:
   - ✗ Node 3 is not reachable from node 1
   - Error message: "Node 3 cannot be reached from node 1. Check node links."
6. Operator drags hop back to original position
7. Operator changes destination scope to `k8s-staging`
8. System validates scope:
   - ⚠ Warning: "Scope k8s-staging is not owned by final hop node 3"
9. Operator reviews warning and decides to proceed
10. Operator clicks "Save"
11. System updates chain with warning acknowledgment

### Flow 3: Duplicating Chain

1. Operator clicks "Duplicate" on existing chain
2. Chain editor opens with copied configuration
3. System appends " (copy)" to chain name
4. Operator modifies hops or scope as needed
5. Operator saves new chain

## Backend API Requirements

### Enhanced API Endpoints

#### `GET /api/v1/chains`

Add query parameters for filtering:
- `?enabled=true` - filter by enabled status
- `?search=prod` - search by name

Response includes validation status:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "prod-k8s-path",
      "hops": [1, 2, 3],
      "destination_scope": "k8s-prod",
      "enabled": true,
      "validation_status": "valid",
      "validation_errors": [],
      "created_at": "2026-04-27T10:00:00Z"
    }
  ]
}
```

#### `POST /api/v1/chains/validate`

Purpose: Validate chain configuration before saving

Auth: Control plane Bearer token required

Request:
```json
{
  "name": "prod-k8s-path",
  "hops": [1, 2, 3],
  "destination_scope": "k8s-prod"
}
```

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "valid": true,
    "errors": [],
    "warnings": [
      "Scope k8s-staging is not owned by final hop node 3"
    ],
    "hop_connectivity": [
      {"from": 1, "to": 2, "reachable": true},
      {"from": 2, "to": 3, "reachable": true}
    ],
    "scope_ownership": {
      "scope": "k8s-prod",
      "owner_node_id": 3,
      "valid": true
    }
  }
}
```

#### `POST /api/v1/chains/preview`

Purpose: Preview compiled chain configuration

Auth: Control plane Bearer token required

Request:
```json
{
  "name": "prod-k8s-path",
  "hops": [1, 2, 3],
  "destination_scope": "k8s-prod"
}
```

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "compiled_config": {
      "chain_id": "preview",
      "name": "prod-k8s-path",
      "hops": [
        {"node_id": 1, "node_name": "edge-node-a", "mode": "edge"},
        {"node_id": 2, "node_name": "relay-node-b", "mode": "relay"},
        {"node_id": 3, "node_name": "k8s-gateway-node", "mode": "relay"}
      ],
      "destination_scope": "k8s-prod",
      "routing_path": "user → edge-node-a → relay-node-b → k8s-gateway-node → target(k8s-prod)"
    }
  }
}
```

#### `GET /api/v1/nodes/scopes`

Purpose: List available destination scopes

Auth: Control plane Bearer token required

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "scope_key": "k8s-prod",
      "owner_node_id": 3,
      "owner_node_name": "k8s-gateway-node",
      "description": "Production Kubernetes cluster"
    },
    {
      "scope_key": "k8s-staging",
      "owner_node_id": 4,
      "owner_node_name": "k8s-staging-node",
      "description": "Staging Kubernetes cluster"
    }
  ]
}
```

### Validation Logic

Backend must implement validation rules:

1. **Hop Connectivity Validation**
   ```go
   func validateHopConnectivity(hops []int) (bool, []string) {
     errors := []string{}
     for i := 0; i < len(hops)-1; i++ {
       link := getNodeLink(hops[i], hops[i+1])
       if link == nil {
         errors = append(errors, fmt.Sprintf("Node %d cannot reach node %d", hops[i], hops[i+1]))
       }
     }
     return len(errors) == 0, errors
   }
   ```

2. **Scope Ownership Validation**
   ```go
   func validateScopeOwnership(scope string, finalHopNodeID int) (bool, string) {
     node := getNodeByID(finalHopNodeID)
     if node.ScopeKey != scope {
       return false, fmt.Sprintf("Scope %s is not owned by node %d", scope, finalHopNodeID)
     }
     return true, ""
   }
   ```

3. **First Hop Edge Node Validation**
   ```go
   func validateFirstHopIsEdge(firstHopNodeID int) (bool, string) {
     node := getNodeByID(firstHopNodeID)
     if node.Mode != "edge" {
       return false, "First hop must be an edge node"
     }
     return true, ""
   }
   ```

## Frontend Requirements

### Component Structure

```
/chains
├── ChainListPage
│   ├── ChainListHeader
│   │   ├── SearchInput
│   │   ├── FilterDropdown
│   │   └── CreateButton
│   ├── ChainTable
│   │   ├── ChainRow
│   │   │   ├── ChainIDCell
│   │   │   ├── ChainNameCell
│   │   │   ├── HopsVisualization
│   │   │   ├── ScopeCell
│   │   │   ├── StatusBadge
│   │   │   └── ActionButtons
│   │   └── EmptyState
│   └── ChainEditorDrawer
│       ├── ChainNameInput
│       ├── ScopeSelector
│       ├── HopList
│       │   ├── HopCard (draggable)
│       │   └── AddHopButton
│       ├── ValidationPanel
│       ├── PreviewButton
│       └── SaveButton
└── CompilationPreviewModal
    ├── JSONViewer
    ├── RoutingPathVisualization
    └── CopyButton
```

### Drag-and-Drop Implementation

Use `@dnd-kit/core` for drag-and-drop:

```tsx
import { DndContext, closestCenter } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';

function HopList({ hops, onReorder }) {
  const handleDragEnd = (event) => {
    const { active, over } = event;
    if (active.id !== over.id) {
      const oldIndex = hops.findIndex(h => h.id === active.id);
      const newIndex = hops.findIndex(h => h.id === over.id);
      onReorder(arrayMove(hops, oldIndex, newIndex));
    }
  };

  return (
    <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
      <SortableContext items={hops} strategy={verticalListSortingStrategy}>
        {hops.map((hop, index) => (
          <SortableHopCard key={hop.id} hop={hop} index={index} />
        ))}
      </SortableContext>
    </DndContext>
  );
}
```

### Real-Time Validation

Implement debounced validation:

```tsx
const [validationResult, setValidationResult] = useState(null);

const validateChain = useMemo(
  () => debounce(async (chainData) => {
    const result = await api.post('/api/v1/chains/validate', chainData);
    setValidationResult(result.data);
  }, 500),
  []
);

useEffect(() => {
  if (chainName && hops.length > 0 && destinationScope) {
    validateChain({ name: chainName, hops, destination_scope: destinationScope });
  }
}, [chainName, hops, destinationScope]);
```

### Scope Selector with Autocomplete

```tsx
function ScopeSelector({ value, onChange, finalHopNodeId }) {
  const { data: scopes } = useSWR('/api/v1/nodes/scopes');
  
  const suggestedScopes = scopes?.filter(
    s => s.owner_node_id === finalHopNodeId
  );

  return (
    <Autocomplete
      value={value}
      onChange={onChange}
      options={scopes?.map(s => s.scope_key) || []}
      renderOption={(props, option) => {
        const scope = scopes.find(s => s.scope_key === option);
        const isSuggested = scope.owner_node_id === finalHopNodeId;
        return (
          <li {...props}>
            {option}
            {isSuggested && <Badge>Recommended</Badge>}
          </li>
        );
      }}
    />
  );
}
```

## Acceptance Criteria

### AC1. Chain List Display

- [ ] Chain list displays all chains with correct data
- [ ] Hops are visualized as `1 → 2 → 3`
- [ ] Validation status badge shows valid/invalid state
- [ ] Search and filter work correctly

### AC2. Chain Editor

- [ ] Chain editor opens when clicking "Create" or "Edit"
- [ ] All form fields are populated correctly
- [ ] Validation runs automatically on input change
- [ ] Validation errors and warnings display correctly

### AC3. Drag-and-Drop

- [ ] Hop cards can be dragged and reordered
- [ ] Visual feedback during drag operation
- [ ] Hop order updates correctly after drop
- [ ] Validation re-runs after reorder

### AC4. Scope Validation

- [ ] Scope selector shows available scopes
- [ ] Recommended scopes are highlighted
- [ ] Invalid scope shows warning message
- [ ] Scope ownership information displays correctly

### AC5. Compilation Preview

- [ ] Preview button opens modal with compiled configuration
- [ ] JSON is formatted and syntax-highlighted
- [ ] Routing path visualization is correct
- [ ] Copy button copies JSON to clipboard

### AC6. Chain Operations

- [ ] Creating new chain saves correctly
- [ ] Editing existing chain updates correctly
- [ ] Duplicating chain creates copy with modified name
- [ ] Deleting chain requires confirmation

## Accessibility

- Drag-and-drop must support keyboard navigation
- Screen reader announces hop reorder actions
- Validation errors are announced to screen readers
- All interactive elements have proper ARIA labels

## Performance

- Chain list pagination for large datasets (>100 chains)
- Debounced validation to reduce API calls
- Optimistic UI updates for better perceived performance
- Lazy loading for compilation preview modal

## Future Enhancements

- Visual topology graph showing all chains
- Chain templates for common patterns
- Bulk chain operations (enable/disable multiple)
- Chain usage analytics (which chains are most used in routes)
- Chain testing tool (simulate request through chain)
