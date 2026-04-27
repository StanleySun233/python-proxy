# Routes Management UI

## Background

Current `/routes` page requires operators to manually input match values, chain IDs, and destination scopes, leading to errors and inefficiency. The page needs intelligent form assistance that automatically loads available options, validates input, and reduces manual typing.

## Goal

Enhance `/routes` page to automatically load available chains, scopes, and match patterns, provide autocomplete and validation, and streamline route rule creation with minimal manual input.

## Functional Requirements

### FR1. Route Rule List View

- Display all route rules in a table with columns:
  - Priority
  - Match Type
  - Match Value
  - Action Type
  - Chain (visual: `1 → 2 → 3`)
  - Destination Scope
  - Enabled Status
  - Actions (Edit/Delete/Duplicate)
- Support sorting by priority
- Support filtering by enabled status, match type, action type
- Support search by match value
- Display rule validation status

### FR2. Smart Route Rule Editor

- Auto-load available chains from `/api/v1/chains`
- Auto-load available scopes from `/api/v1/nodes/scopes`
- Provide match type dropdown with predefined options
- Provide match value suggestions based on match type
- Validate match value format in real-time
- Show chain preview when chain is selected
- Show scope ownership information when scope is selected

### FR3. Match Type and Value Intelligence

Match types and their input assistance:

- `domain`: Autocomplete from recent domains, validate domain format
- `domain_suffix`: Autocomplete from recent suffixes, validate format
- `ip_cidr`: Validate CIDR notation, provide format hints
- `ip_range`: Validate IP range format, provide format hints
- `port`: Validate port number (1-65535)
- `url_regex`: Validate regex syntax, provide regex tester
- `default`: No validation, catch-all rule

### FR4. Chain Selection with Preview

- Dropdown shows all available chains
- Each option displays: `Chain Name (1 → 2 → 3)`
- Selecting chain shows preview card:
  - Chain name
  - Hop sequence with node names
  - Destination scope
  - Validation status
- Warn if selected chain is disabled
- Warn if selected chain has validation errors

### FR5. Scope Selection with Validation

- Dropdown shows all available scopes
- Each option displays: `scope-key (owned by node-name)`
- Selecting scope shows ownership information
- Warn if scope does not match chain's final hop node
- Suggest compatible scopes based on selected chain

### FR6. Priority Management

- Auto-suggest next available priority
- Warn if priority conflicts with existing rule
- Allow manual priority adjustment
- Show priority order preview in rule list

### FR7. Rule Validation

Validation rules:
- Match value must be valid for match type
- Chain must exist and be enabled
- Destination scope must exist
- Scope should match chain's final hop node (warning, not error)
- Priority must be unique (warning if duplicate)
- Enabled rules must pass all validation

## User Interaction Flow

### Flow 1: Creating a Domain-Based Route Rule

1. Operator navigates to `/routes`
2. Clicks "Create Route Rule" button
3. Route rule editor drawer opens
4. System auto-suggests priority: `100` (next available)
5. Operator selects match type: `domain`
6. Match value input shows placeholder: `example.com`
7. Operator types: `api.internal.com`
8. System validates domain format: ✓ Valid
9. Operator selects action type: `proxy`
10. Chain dropdown loads available chains:
    - `prod-k8s-path (1 → 2 → 3)`
    - `staging-path (1 → 4)`
11. Operator selects: `prod-k8s-path (1 → 2 → 3)`
12. System shows chain preview card:
    ```
    Chain: prod-k8s-path
    Hops: edge-node-a → relay-node-b → k8s-gateway-node
    Destination Scope: k8s-prod
    Status: ✓ Valid
    ```
13. Scope dropdown loads available scopes:
    - `k8s-prod (owned by k8s-gateway-node)` ← Recommended
    - `k8s-staging (owned by k8s-staging-node)`
14. System pre-selects recommended scope: `k8s-prod`
15. Operator reviews configuration
16. Operator clicks "Save"
17. System creates route rule and displays success message

### Flow 2: Creating an IP CIDR Route Rule

1. Operator clicks "Create Route Rule"
2. Operator selects match type: `ip_cidr`
3. Match value input shows placeholder: `10.0.0.0/24`
4. Match value input shows format hint: "Enter CIDR notation (e.g., 192.168.1.0/24)"
5. Operator types: `10.0.0.0/16`
6. System validates CIDR format: ✓ Valid
7. Operator continues with chain and scope selection
8. System validates complete rule configuration
9. Operator saves rule

### Flow 3: Editing Existing Route Rule

1. Operator clicks "Edit" on existing rule
2. Route rule editor opens with current configuration
3. All fields are populated with current values
4. Chain preview shows current chain details
5. Operator changes match value from `api.internal.com` to `*.internal.com`
6. System validates new match value: ✓ Valid
7. Operator changes chain to different chain
8. System shows warning: "New chain's scope (k8s-staging) differs from current scope (k8s-prod)"
9. Operator updates scope to match new chain
10. Operator saves changes

### Flow 4: Handling Validation Errors

1. Operator creates route rule with invalid CIDR: `10.0.0.0/33`
2. System shows error: "Invalid CIDR notation. Prefix must be 0-32."
3. Operator corrects to: `10.0.0.0/24`
4. Error clears
5. Operator selects disabled chain
6. System shows warning: "Selected chain is disabled. Enable chain before publishing policy."
7. Operator acknowledges warning and saves rule

## Backend API Requirements

### Enhanced API Endpoints

#### `GET /api/v1/route-rules`

Add query parameters:
- `?enabled=true` - filter by enabled status
- `?match_type=domain` - filter by match type
- `?action_type=proxy` - filter by action type
- `?search=api` - search by match value

Response includes validation status and chain details:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "priority": 100,
      "match_type": "domain",
      "match_value": "api.internal.com",
      "action_type": "proxy",
      "chain_id": 1,
      "chain": {
        "id": 1,
        "name": "prod-k8s-path",
        "hops": [1, 2, 3],
        "enabled": true
      },
      "destination_scope": "k8s-prod",
      "enabled": true,
      "validation_status": "valid",
      "validation_errors": [],
      "validation_warnings": [],
      "created_at": "2026-04-27T10:00:00Z"
    }
  ]
}
```

#### `POST /api/v1/route-rules/validate`

Purpose: Validate route rule configuration before saving

Auth: Control plane Bearer token required

Request:
```json
{
  "priority": 100,
  "match_type": "domain",
  "match_value": "api.internal.com",
  "action_type": "proxy",
  "chain_id": 1,
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
      "Priority 100 conflicts with existing rule #5"
    ],
    "match_value_validation": {
      "valid": true,
      "format": "domain",
      "message": "Valid domain format"
    },
    "chain_validation": {
      "valid": true,
      "chain_enabled": true,
      "chain_hops": [1, 2, 3]
    },
    "scope_validation": {
      "valid": true,
      "scope_exists": true,
      "scope_owner_node_id": 3,
      "matches_chain_final_hop": true
    }
  }
}
```

#### `GET /api/v1/route-rules/match-types`

Purpose: List available match types with descriptions

Auth: Control plane Bearer token required

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "type": "domain",
      "label": "Domain",
      "description": "Match exact domain name",
      "placeholder": "example.com",
      "validation_regex": "^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$"
    },
    {
      "type": "domain_suffix",
      "label": "Domain Suffix",
      "description": "Match domain suffix (e.g., *.example.com)",
      "placeholder": "*.example.com",
      "validation_regex": "^\\*\\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$"
    },
    {
      "type": "ip_cidr",
      "label": "IP CIDR",
      "description": "Match IP address range in CIDR notation",
      "placeholder": "10.0.0.0/24",
      "validation_regex": "^([0-9]{1,3}\\.){3}[0-9]{1,3}/[0-9]{1,2}$"
    },
    {
      "type": "ip_range",
      "label": "IP Range",
      "description": "Match IP address range",
      "placeholder": "10.0.0.1-10.0.0.255",
      "validation_regex": "^([0-9]{1,3}\\.){3}[0-9]{1,3}-([0-9]{1,3}\\.){3}[0-9]{1,3}$"
    },
    {
      "type": "port",
      "label": "Port",
      "description": "Match port number",
      "placeholder": "8080",
      "validation_regex": "^[0-9]{1,5}$"
    },
    {
      "type": "url_regex",
      "label": "URL Regex",
      "description": "Match URL using regular expression",
      "placeholder": "^https://api\\..*",
      "validation_regex": null
    },
    {
      "type": "default",
      "label": "Default",
      "description": "Catch-all rule (lowest priority)",
      "placeholder": "*",
      "validation_regex": null
    }
  ]
}
```

#### `GET /api/v1/route-rules/suggestions`

Purpose: Get match value suggestions based on match type

Auth: Control plane Bearer token required

Query parameters:
- `match_type=domain`
- `query=api` (optional, for filtering)

Response:
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "match_type": "domain",
    "suggestions": [
      "api.internal.com",
      "api.staging.com",
      "api.prod.com"
    ],
    "source": "recent_rules"
  }
}
```

### Validation Logic

Backend must implement match value validation:

```go
func validateMatchValue(matchType, matchValue string) (bool, string) {
  switch matchType {
  case "domain":
    return validateDomain(matchValue)
  case "domain_suffix":
    return validateDomainSuffix(matchValue)
  case "ip_cidr":
    return validateIPCIDR(matchValue)
  case "ip_range":
    return validateIPRange(matchValue)
  case "port":
    return validatePort(matchValue)
  case "url_regex":
    return validateRegex(matchValue)
  case "default":
    return true, ""
  default:
    return false, "Unknown match type"
  }
}

func validateIPCIDR(cidr string) (bool, string) {
  _, _, err := net.ParseCIDR(cidr)
  if err != nil {
    return false, "Invalid CIDR notation"
  }
  return true, ""
}

func validatePort(port string) (bool, string) {
  p, err := strconv.Atoi(port)
  if err != nil || p < 1 || p > 65535 {
    return false, "Port must be between 1 and 65535"
  }
  return true, ""
}

func validateRegex(pattern string) (bool, string) {
  _, err := regexp.Compile(pattern)
  if err != nil {
    return false, fmt.Sprintf("Invalid regex: %s", err.Error())
  }
  return true, ""
}
```

## Frontend Requirements

### Component Structure

```
/routes
├── RouteRuleListPage
│   ├── RouteRuleListHeader
│   │   ├── SearchInput
│   │   ├── FilterDropdown
│   │   └── CreateButton
│   ├── RouteRuleTable
│   │   ├── RouteRuleRow
│   │   │   ├── PriorityCell
│   │   │   ├── MatchTypeCell
│   │   │   ├── MatchValueCell
│   │   │   ├── ActionTypeCell
│   │   │   ├── ChainVisualization
│   │   │   ├── ScopeCell
│   │   │   ├── StatusBadge
│   │   │   └── ActionButtons
│   │   └── EmptyState
│   └── RouteRuleEditorDrawer
│       ├── PriorityInput
│       ├── MatchTypeSelector
│       ├── MatchValueInput (smart input based on match type)
│       ├── ActionTypeSelector
│       ├── ChainSelector
│       ├── ChainPreviewCard
│       ├── ScopeSelector
│       ├── ScopeInfoCard
│       ├── ValidationPanel
│       └── SaveButton
└── RegexTesterModal
    ├── RegexInput
    ├── TestStringInput
    └── MatchResult
```

### Smart Match Value Input

Implement match-type-specific input components:

```tsx
function MatchValueInput({ matchType, value, onChange }) {
  const { data: matchTypes } = useSWR('/api/v1/route-rules/match-types');
  const { data: suggestions } = useSWR(
    matchType ? `/api/v1/route-rules/suggestions?match_type=${matchType}` : null
  );

  const matchTypeConfig = matchTypes?.find(mt => mt.type === matchType);

  switch (matchType) {
    case 'domain':
    case 'domain_suffix':
      return (
        <Autocomplete
          value={value}
          onChange={onChange}
          options={suggestions?.suggestions || []}
          placeholder={matchTypeConfig?.placeholder}
          helperText={matchTypeConfig?.description}
        />
      );
    
    case 'ip_cidr':
      return (
        <Input
          value={value}
          onChange={onChange}
          placeholder={matchTypeConfig?.placeholder}
          helperText="Enter CIDR notation (e.g., 192.168.1.0/24)"
          pattern={matchTypeConfig?.validation_regex}
        />
      );
    
    case 'url_regex':
      return (
        <div>
          <Input
            value={value}
            onChange={onChange}
            placeholder={matchTypeConfig?.placeholder}
            helperText="Enter regular expression pattern"
          />
          <Button onClick={() => openRegexTester(value)}>
            Test Regex
          </Button>
        </div>
      );
    
    default:
      return (
        <Input
          value={value}
          onChange={onChange}
          placeholder={matchTypeConfig?.placeholder}
        />
      );
  }
}
```

### Chain Selector with Preview

```tsx
function ChainSelector({ value, onChange }) {
  const { data: chains } = useSWR('/api/v1/chains');
  const selectedChain = chains?.find(c => c.id === value);

  return (
    <div>
      <Select value={value} onChange={onChange}>
        {chains?.map(chain => (
          <option key={chain.id} value={chain.id}>
            {chain.name} ({chain.hops.join(' → ')})
          </option>
        ))}
      </Select>
      
      {selectedChain && (
        <ChainPreviewCard chain={selectedChain} />
      )}
    </div>
  );
}

function ChainPreviewCard({ chain }) {
  return (
    <Card>
      <CardHeader>
        <h4>{chain.name}</h4>
        {!chain.enabled && <Badge variant="warning">Disabled</Badge>}
      </CardHeader>
      <CardBody>
        <div>
          <strong>Hops:</strong>
          <HopVisualization hops={chain.hops} />
        </div>
        <div>
          <strong>Destination Scope:</strong> {chain.destination_scope}
        </div>
        <div>
          <strong>Status:</strong>
          {chain.validation_status === 'valid' ? (
            <Badge variant="success">Valid</Badge>
          ) : (
            <Badge variant="error">Invalid</Badge>
          )}
        </div>
      </CardBody>
    </Card>
  );
}
```

### Real-Time Validation

```tsx
const [validationResult, setValidationResult] = useState(null);

const validateRule = useMemo(
  () => debounce(async (ruleData) => {
    const result = await api.post('/api/v1/route-rules/validate', ruleData);
    setValidationResult(result.data);
  }, 500),
  []
);

useEffect(() => {
  if (matchType && matchValue && chainId && destinationScope) {
    validateRule({
      priority,
      match_type: matchType,
      match_value: matchValue,
      action_type: actionType,
      chain_id: chainId,
      destination_scope: destinationScope
    });
  }
}, [priority, matchType, matchValue, actionType, chainId, destinationScope]);
```

### Scope Selector with Recommendations

```tsx
function ScopeSelector({ value, onChange, chainId }) {
  const { data: scopes } = useSWR('/api/v1/nodes/scopes');
  const { data: chains } = useSWR('/api/v1/chains');
  
  const selectedChain = chains?.find(c => c.id === chainId);
  const finalHopNodeId = selectedChain?.hops[selectedChain.hops.length - 1];
  
  const recommendedScopes = scopes?.filter(
    s => s.owner_node_id === finalHopNodeId
  );

  return (
    <div>
      <Select value={value} onChange={onChange}>
        {recommendedScopes?.length > 0 && (
          <optgroup label="Recommended">
            {recommendedScopes.map(scope => (
              <option key={scope.scope_key} value={scope.scope_key}>
                {scope.scope_key} (owned by {scope.owner_node_name})
              </option>
            ))}
          </optgroup>
        )}
        <optgroup label="All Scopes">
          {scopes?.map(scope => (
            <option key={scope.scope_key} value={scope.scope_key}>
              {scope.scope_key} (owned by {scope.owner_node_name})
            </option>
          ))}
        </optgroup>
      </Select>
      
      {value && !recommendedScopes?.find(s => s.scope_key === value) && (
        <Alert variant="warning">
          Selected scope is not owned by chain's final hop node.
          This may cause routing errors.
        </Alert>
      )}
    </div>
  );
}
```

## Acceptance Criteria

### AC1. Route Rule List Display

- [ ] Route rule list displays all rules with correct data
- [ ] Chains are visualized as `1 → 2 → 3`
- [ ] Validation status badge shows valid/invalid state
- [ ] Search, filter, and sort work correctly

### AC2. Smart Route Rule Editor

- [ ] Match type selector loads available types
- [ ] Match value input adapts to selected match type
- [ ] Chain selector loads available chains
- [ ] Scope selector loads available scopes
- [ ] All dropdowns display correctly formatted options

### AC3. Match Value Validation

- [ ] Domain format validation works correctly
- [ ] CIDR notation validation works correctly
- [ ] Port number validation works correctly
- [ ] Regex syntax validation works correctly
- [ ] Validation errors display in real-time

### AC4. Chain and Scope Intelligence

- [ ] Chain preview card displays when chain is selected
- [ ] Recommended scopes are highlighted
- [ ] Warning displays when scope doesn't match chain
- [ ] Disabled chain warning displays correctly

### AC5. Rule Validation

- [ ] Complete rule validation runs on input change
- [ ] Validation errors prevent saving
- [ ] Validation warnings allow saving with acknowledgment
- [ ] Validation panel displays all errors and warnings

### AC6. Rule Operations

- [ ] Creating new rule saves correctly
- [ ] Editing existing rule updates correctly
- [ ] Duplicating rule creates copy with modified priority
- [ ] Deleting rule requires confirmation

## Accessibility

- All form inputs have proper labels
- Validation errors are announced to screen readers
- Dropdown options are keyboard navigable
- All interactive elements have proper ARIA attributes

## Performance

- Debounced validation to reduce API calls
- Cached chain and scope data with SWR
- Optimistic UI updates for better perceived performance
- Lazy loading for large rule lists (pagination)

## Future Enhancements

- Bulk rule import from CSV/JSON
- Rule templates for common patterns
- Rule testing tool (simulate request matching)
- Rule conflict detection and resolution
- Rule usage analytics (which rules are most matched)
- Visual rule priority editor (drag to reorder)
