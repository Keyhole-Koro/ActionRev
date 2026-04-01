# Paper-in-Paper Migration Notes

Date: 2026-04-01

## Goal

Replace the current workspace graph UI with `@keyhole-koro/paper-in-paper` and stop maintaining backward compatibility with the existing React Flow-based renderer.

The user plans to bring `paper-in-paper` into this repo as a submodule and edit it directly:

- `https://github.com/Keyhole-Koro/paper-in-paper.git`

## Confirmed Package Surface

From the published package `@keyhole-koro/paper-in-paper@0.1.0`:

- Main component: `PaperCanvas`
- Main data type: `Paper`
  - `id`
  - `title`
  - `description`
  - `content`
  - `parentId`
  - `childIds`
- Utilities:
  - `buildPaperMap`
  - `findRootId`
  - expansion helpers:
    - `openNode`
    - `closeNode`
    - `setPrimaryNode`
    - `expansionReducer`

## Important Constraints

- The package is tree-first, not graph-first.
- It expects exactly one root node.
- The public API does not expose edge rendering.
- The public API does not currently expose custom node renderers.
- The public API does not currently expose:
  - search integration
  - filter integration
  - source document cards
  - neighbor expansion actions
  - custom badges like `New`
  - minimap / edge overlay
- Package peer dependency:
  - `framer-motion >= 11`

## Current App Dependencies On React Flow

Frontend currently depends on `@xyflow/react` in these places:

- `frontend/app/globals.css`
  - imports React Flow CSS
- `frontend/features/graph/components/graph-canvas-panel.tsx`
  - `ReactFlow`
  - `MiniMap`
  - `Controls`
  - `Handle`
  - node click handling
  - edge display
- `frontend/features/graph/types/graph-canvas.ts`
  - React Flow `Node`/`Edge`-based UI model
- `frontend/features/graph/model/to-graph-canvas.ts`
  - layout
  - highlight/dim logic
  - search highlight
  - edge styling
- `frontend/features/workspaces/components/workspace-graph-page.tsx`
  - expanded node state
  - search/filter integration
  - `ExpandNeighbors`
  - source documents injection
  - `isNew` / `isExpanding` annotations

## Backend Shape Today

Backend still returns a graph schema, not a tree schema:

- `proto/synthify/graph/v1/graph.proto`
- `proto/synthify/graph/v1/graph_types.proto`
- `backend/internal/repository/mock/graph.go`

Current graph behavior includes:

- canonical/document scope nodes
- multiple edge types
  - `HIERARCHICAL`
  - `RELATED_TO`
  - `SUPPORTS`
  - `CONTRADICTS`
  - `CAUSES`
  - `MEASURED_BY`
  - `MENTIONS`
- `ExpandNeighbors`
- cross-document graph composition

## Migration Implications

The current UI can be replaced, but only if we define a graph-to-tree projection.

### Can stay

- graph fetching in `useGetGraph()`
- workspace/document upload/rename flows
- backend graph API for now

### Can be deleted

- `GraphCanvas`
- `toGraphCanvas()`
- `GraphCanvasPanel`
- direct `@xyflow/react` usage
- React Flow CSS import

## Required New Adapter

Need a new adapter, probably:

- `frontend/features/graph/model/to-paper-tree-model.ts`

It should convert the current `Graph` into `Paper[]` or `PaperMap`.

## Decisions Still Needed During Implementation

### 1. Parent/child rule

Recommended priority:

- prefer `HIERARCHICAL` edges
- otherwise infer hierarchy from `level`
- use canonical/document scope as a tiebreaker
- allow fixed root such as `cn_workspace_strategy` if needed

### 2. Non-tree relations

`PaperCanvas` cannot directly show arbitrary edges.

Need to decide how to represent:

- `RELATED_TO`
- `SUPPORTS`
- `CONTRADICTS`
- `CAUSES`
- `MEASURED_BY`
- `MENTIONS`

Likely options:

- inline metadata in node content
- related node list
- separate side panel
- badges / annotations

### 3. Source document rendering

`Paper` only has `title`, `description`, and `content`.

So source documents must be either:

- embedded into `content`, or
- stored in a parallel metadata map in the app, or
- supported by direct edits to the `paper-in-paper` package

### 4. Expand neighbors

Need to define how `ExpandNeighbors` inserts newly fetched nodes into the tree:

- attach as direct child of the expanded node
- attach according to inferred hierarchy
- or show in a separate related-node region

## Recommended Migration Order

1. Add `paper-in-paper` to the repo as a submodule.
2. Inspect and edit package internals directly.
3. Add `framer-motion` to frontend if still required at integration time.
4. Create a first `toPaperTreeModel` adapter in app code.
5. Replace `GraphCanvasPanel` with a new `PaperTreePanel`.
6. Move search/filter/neighbor-expansion behavior outside the old React Flow model.
7. Remove `@xyflow/react` and its CSS import.

## Expected Package Edits

Direct edits to `paper-in-paper` will probably be needed for at least one of these:

- custom node body rendering
- richer node metadata access
- external callbacks for node selection / expansion
- integration points for app-managed search/filter state
- a way to surface non-tree relations

## Replacement Boundary

Target boundary after migration:

- app keeps:
  - graph fetching
  - workspace/document state
  - filters/search/upload actions
  - neighbor expansion requests
- `paper-in-paper` handles:
  - tree rendering
  - node expansion model
  - motion / hierarchy presentation

## Current Recommendation

Do not preserve backward compatibility.

Treat this as a full replacement of the current graph renderer, with a new tree-oriented UI model and direct package edits once the submodule is added.
