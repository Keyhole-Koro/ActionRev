# Paper-in-Paper Change Plan

Date: 2026-04-01

## Purpose

This document separates the migration work into:

- changes to make inside `vendor/paper-in-paper`
- changes to make inside Synthify

Backward compatibility is not required.

The plan assumes we are free to reshape `paper-in-paper` into a Synthify-oriented tree renderer instead of preserving its current generic library API.

## Current Decision

We are **not** treating `paper-in-paper` as a fixed external adapter target.

We are treating it as:

- an internal submodule
- a reusable base
- something we can redesign directly for Synthify

That means:

- `adapter-first` design is no longer the priority
- `renderer and data model redesign` is the priority

## What Paper-in-Paper Should Own

`paper-in-paper` should become responsible for:

- hierarchical tree rendering
- node expansion and collapse behavior
- primary child / focused branch behavior
- drag and reorder interactions
- tree motion and layout transitions
- node-level presentation primitives
- tree-level selection / highlight / dim behavior

It should become the dedicated visual engine for Synthify's workspace tree.

## What Synthify Should Own

Synthify should continue owning:

- graph fetching
- workspace/document fetching
- upload flow
- rename flow
- search query state
- filter state
- backend `ExpandNeighbors` requests
- graph-to-tree transformation rules
- domain-specific metadata
- side panels / top bar / workspace page shell

## Paper-in-Paper: Required Changes

### 1. Redesign the core node type

Current `Paper` shape is too small:

- `id`
- `title`
- `description`
- `content`
- `parentId`
- `childIds`

Synthify needs a richer node model.

Recommended additions:

- `category`
- `scope`
- `documentId`
- `documentIds`
- `sourceDocuments`
- `sourceChunkIds`
- `relatedNodes`
- `badges`
- `isNew`
- `isLoading`
- `actions`

This can be done either by:

- expanding `Paper` directly, or
- introducing a generic metadata field such as `meta`

Direct expansion is probably simpler for now.

### 2. Add external interaction callbacks

`PaperCanvas` needs explicit extension points.

Required callbacks:

- `onNodeClick`
- `onNodeOpen`
- `onNodeClose`
- `onNodeReorder`
- `onNodeAction`
- `onRequestExpandNeighbors`

Without these, Synthify cannot coordinate backend-driven behavior cleanly.

### 3. Add custom node rendering hooks

Synthify needs to control node body rendering.

Required extension points:

- `renderNode`
- `renderNodeHeader`
- `renderNodeBody`
- `renderNodeFooter`

Minimum requirement:

- one `renderNode` callback that receives the full node data and interaction state

Why this is needed:

- source documents need custom cards
- category/scope labels need custom styling
- `Expand neighbors` button must appear conditionally
- `New` / `Loading` badges must be app-aware

### 4. Add externally controlled visual state

Synthify already has search/filter/highlight logic.

`PaperCanvas` should accept externally controlled state such as:

- `expandedNodeIds`
- `selectedNodeId`
- `highlightedNodeIds`
- `dimmedNodeIds`
- `hiddenNodeIds`

This avoids pushing app logic into the renderer internals.

### 5. Support node-level actions

Tree nodes need action slots or structured actions.

At minimum:

- per-node action buttons
- loading state on an action
- disabled state on an action

This is needed for:

- `Expand neighbors`
- future node actions

### 6. Support non-tree relation display

Synthify still has non-hierarchical graph semantics.

The renderer does not need full edge drawing, but it should support:

- related node summaries
- relation badges
- relation sections inside expanded nodes

Otherwise important graph semantics disappear.

### 7. Separate generic tree engine from current demo assumptions

The current package still carries demo-oriented structure and internal assumptions.

It should be cleaned into:

- stable tree primitives
- stable public types
- stable render hooks
- minimal demo coupling

### 8. Revisit drag-and-drop semantics

`REORDER` already exists, which is good.

Still needs review:

- should drag be optional via prop
- should reorder be opt-in
- should app be able to veto a reorder
- should reorder work only within allowed parent scopes

Synthify may not want arbitrary tree mutation everywhere.

### 9. Clarify styling contract

The package should expose a clear styling boundary.

Needed:

- stable class names or slot props
- theme variables
- a way to override spacing, width, and typography

Right now we should assume the default look is not yet final for Synthify.

## Synthify: Required Changes

### 1. Replace React Flow types and renderer

These should be removed:

- `frontend/features/graph/components/graph-canvas-panel.tsx`
- `frontend/features/graph/types/graph-canvas.ts`
- `frontend/features/graph/model/to-graph-canvas.ts`
- `@xyflow/react`
- React Flow CSS import in `frontend/app/globals.css`

### 2. Introduce a tree model transformation

Need a new transformation layer, probably:

- `frontend/features/graph/model/to-paper-tree-model.ts`

It should:

- convert backend graph data into tree node data for `paper-in-paper`
- preserve domain metadata needed by the renderer
- decide root/parent/child structure
- decide how non-tree relations are carried

### 3. Define graph-to-tree projection rules

Synthify must decide:

- which edge types define hierarchy
- how to infer parentage when no hierarchical edge exists
- what the canonical root should be
- where document-scoped nodes belong
- how `ExpandNeighbors` results are inserted

Suggested starting rule:

- prefer `HIERARCHICAL`
- otherwise infer from `level`
- prefer canonical nodes above document nodes

### 4. Move workspace page to the new renderer contract

`frontend/features/workspaces/components/workspace-graph-page.tsx` will need to be simplified.

It should stop injecting React Flow node style state and instead provide:

- search state
- filter state
- selected / expanded state
- source document data
- expand-neighbors action wiring

### 5. Keep backend graph API for now

No immediate proto rewrite is required.

The backend graph service can continue returning:

- `Graph`
- `Node`
- `Edge`

The transformation into a tree can remain frontend-side initially.

### 6. Re-evaluate search and filters for tree UI

Current search/filter behavior is graph-oriented.

For the tree UI, Synthify should decide:

- whether unmatched nodes are hidden or dimmed
- whether parent chain should remain visible for matched descendants
- whether filters prune subtrees or annotate them

### 7. Define how non-tree relations appear in Synthify UI

Needed decisions:

- inline on each node
- separate right-side panel
- relation badges in expanded body
- future overlay or secondary view

This is a product-level decision, not just a renderer concern.

### 8. Replace graph-specific tests

Current and future tests should move away from React Flow assumptions.

Need:

- renderer-level tests for the new tree UI
- workspace E2E tests against the new tree interactions

## Sequence Of Work

Recommended order:

1. Redesign `paper-in-paper` types
2. Add `PaperCanvas` extension points
3. Add custom node rendering
4. Add external visual state control
5. Create Synthify tree model transformation
6. Replace `GraphCanvasPanel`
7. Remove `@xyflow/react`
8. Update tests

## Immediate First Tasks

If work starts now, the first concrete tasks should be:

1. Modify `vendor/paper-in-paper/src/lib/core/types.ts`
   - expand the node shape for Synthify metadata
2. Modify `vendor/paper-in-paper/src/lib/react/PaperCanvas.tsx`
   - accept custom renderers and callbacks
3. Inspect internal node components under `vendor/paper-in-paper/src/lib/react/internal`
   - route richer node data into visible UI
4. Create `frontend/features/graph/model/to-paper-tree-model.ts`
   - convert current backend graph into the new tree input model

## Non-Goals For The First Pass

Do not optimize for these yet:

- preserving existing npm API shape
- preserving React Flow behavior
- perfect general-purpose library design
- backend schema redesign

The first pass should optimize for:

- getting Synthify off React Flow
- making the tree renderer expressive enough
- keeping the data flow understandable
