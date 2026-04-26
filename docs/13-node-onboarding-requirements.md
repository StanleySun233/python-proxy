# 13 Node Onboarding Requirements

## Goal

Support node onboarding in multi-node and partially unreachable networks without assuming every node can directly reach the control plane.

## Problem Statement

Current enrollment only covers the direct model:

- node reaches control plane directly
- control plane approves enrollment
- node exchanges secret and comes online

That is not enough for these real cases:

- `node2` can only reach `node1`
- control plane can reach `node1`, but not `node2`
- control plane knows `node2` address and wants to initiate onboarding from the panel
- control plane must use a relay chain such as `panel -> node1 -> nodeX -> node2`

## Functional Scope

### F1. Dual Onboarding Modes

The product must support two onboarding modes:

- `agent_pull`: target node initiates onboarding through an upstream node or direct control-plane address
- `control_push`: control plane initiates onboarding toward the target node directly or through a relay path

### F2. First-Class Access Path

The product must model onboarding reachability as a first-class object instead of scattering relay settings across node records.

An access path must describe:

- path name
- mode
- target node
- optional entry node
- ordered relay hops
- target host and port
- enabled state

### F3. Panel-Driven Onboarding Task

The panel must be able to create an onboarding task instead of performing an implicit one-shot connection.

An onboarding task must capture:

- onboarding mode
- target node or target address
- selected access path
- requested operator
- current status
- latest status message

### F4. Direct And Multi-Hop Variants

The panel must support these variants:

- direct control-plane to target node
- control-plane to target through a saved relay path
- target node to control-plane through a saved upstream entry node

### F5. Unified Enrollment State Machine

Regardless of path choice, the node lifecycle must converge to the same state machine:

- discovered
- pending
- approved
- active
- degraded
- offline

### F6. Manual Operator Control

The first implementation must prefer explicit operator control over auto-routing.

The operator should be able to:

- define access paths manually
- choose which path to use for a task
- decide between direct and relay onboarding

Automatic path planning can come later.

## Non-Goals For This Iteration

- no full overlay network
- no automatic route search
- no mesh gossip
- no traffic audit trail
- no background auto-onboard based on passive discovery

## User Flows

### Flow A. Pull Through Existing Node

- `node1` is already online
- operator deploys `node2`
- `node2` is configured to use `node1` as upstream onboarding entry
- `node1` forwards onboarding traffic to control plane
- control plane shows the pending node and approval state

### Flow B. Push Directly From Panel

- operator creates a target node draft in panel
- operator enters `node2` host and port
- control plane creates a direct onboarding task
- task attempts control connection to the target node
- successful task moves node to pending or active onboarding state

### Flow C. Push Through Relay Chain

- operator creates relay path `panel -> node1 -> nodeX -> node2`
- operator selects that path when creating onboarding task
- control plane dispatches onboarding over the selected path
- task state tracks each attempt until success or failure

## Acceptance Criteria

- operator can save and list onboarding access paths
- operator can create onboarding tasks from panel
- node-agent can run without direct control-plane URL
- API model distinguishes onboarding mode from route policy
- future relay execution can be added without changing the top-level control-plane objects
