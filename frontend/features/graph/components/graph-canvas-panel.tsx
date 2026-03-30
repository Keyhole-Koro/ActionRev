'use client'

import type { Node, NodeMouseHandler, NodeProps } from '@xyflow/react'
import { Background, Controls, Handle, MiniMap, Position, ReactFlow } from '@xyflow/react'
import type { GraphCanvas, GraphCanvasNodeData } from '../types/graph-canvas'

function GraphCanvasNode({ data }: NodeProps<Node<GraphCanvasNodeData>>) {
  return (
    <div className="relative rounded-[24px] bg-white/92 p-4">
      <Handle
        type="target"
        position={Position.Left}
        className="!h-3 !w-3 !border-2 !border-slate-200 !bg-white"
        style={{ left: -7 }}
      />
      <Handle
        type="source"
        position={Position.Right}
        className="!h-3 !w-3 !border-2 !border-slate-200 !bg-white"
        style={{ right: -7 }}
      />
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="truncate text-sm font-semibold text-slate-900">{data.label}</p>
          <p className="mt-1 text-xs text-slate-400">
            {data.category} · {data.scope} · L{data.level}
          </p>
        </div>
      </div>

      {data.expanded && (
        <div className="mt-4 space-y-4 border-t border-slate-100 pt-4">
          <div>
            <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-400">
              Description
            </p>
            <p className="mt-2 text-xs leading-6 text-slate-600">{data.description || 'No description.'}</p>
          </div>

          <div className="grid grid-cols-2 gap-2 text-xs">
            <div className="rounded-xl bg-slate-50 px-3 py-2">
              <p className="text-slate-400">Chunks</p>
              <p className="mt-1 font-medium text-slate-700">{data.sourceChunkIds.length}</p>
            </div>
            <div className="rounded-xl bg-slate-50 px-3 py-2">
              <p className="text-slate-400">Sources</p>
              <p className="mt-1 font-medium text-slate-700">{data.sourceDocuments?.length ?? 0}</p>
            </div>
          </div>

          {data.onExpandNeighbors && (
            <button
              type="button"
              onClick={(e) => {
                e.stopPropagation()
                data.onExpandNeighbors!()
              }}
              disabled={data.isExpanding}
              className="w-full rounded-xl border border-slate-200 bg-slate-50 px-3 py-2 text-xs font-medium text-slate-600 transition-colors hover:border-slate-300 hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
            >
              {data.isExpanding ? 'Loading neighbors…' : 'Expand neighbors'}
            </button>
          )}

          <div>
            <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-slate-400">
              Source Documents
            </p>
            <div className="mt-2 space-y-2">
              {data.sourceDocuments?.length ? (
                data.sourceDocuments.map((document) => (
                  <div key={document.id} className="rounded-xl border border-slate-100 bg-slate-50 px-3 py-2">
                    <div className="flex items-center justify-between gap-2">
                      <p className="truncate text-xs font-medium text-slate-700">{document.filename}</p>
                      <span className="text-[11px] text-slate-400">{document.status}</span>
                    </div>
                    <p className="mt-1 text-[11px] text-slate-400">{document.id}</p>
                  </div>
                ))
              ) : (
                <p className="text-xs text-slate-400">No source documents linked.</p>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

const nodeTypes = {
  graphNode: GraphCanvasNode,
}

type GraphCanvasPanelProps = {
  canvas: GraphCanvas | null
  isLoading: boolean
  error: string | null
  emptyMessage?: string
  onNodeSelect?: NodeMouseHandler
}

export function GraphCanvasPanel({
  canvas,
  isLoading,
  error,
  emptyMessage = 'No graph data.',
  onNodeSelect,
}: GraphCanvasPanelProps) {
  return (
    <div className="absolute inset-0">
      {canvas ? (
        <ReactFlow
          fitView
          nodes={canvas.nodes}
          edges={canvas.edges}
          nodeTypes={nodeTypes}
          proOptions={{ hideAttribution: true }}
          defaultEdgeOptions={{ zIndex: 1 }}
          onNodeClick={onNodeSelect}
        >
          <MiniMap
            pannable
            zoomable
            nodeBorderRadius={6}
            nodeColor={(node) =>
              String(node.style?.border ?? '#e2e8f0').replace('1px solid ', '')
            }
            maskColor="rgba(248,250,252,0.7)"
            style={{
              background: '#f8fafc',
              border: '1px solid #e2e8f0',
              borderRadius: 8,
            }}
          />
          <Controls showInteractive={false} />
          <Background color="#e2e8f0" gap={24} size={1} />
        </ReactFlow>
      ) : (
        <div className="flex h-full items-center justify-center text-sm text-slate-300">
          {isLoading ? 'Loading…' : emptyMessage}
        </div>
      )}

      {/* Floating status badge — bottom left */}
      {canvas && (
        <div className="absolute bottom-4 left-4 z-10 rounded-lg border border-slate-200 bg-white/80 px-3 py-1.5 font-mono text-xs text-slate-400 shadow-sm backdrop-blur-sm">
          {canvas.nodes.length} nodes · {canvas.edges.length} edges
        </div>
      )}

      {error && (
        <div className="absolute bottom-4 left-1/2 z-10 -translate-x-1/2">
          <p className="rounded-lg border border-red-200 bg-white/90 px-4 py-2 text-xs text-red-500 shadow-sm backdrop-blur-sm">
            {error}
          </p>
        </div>
      )}
    </div>
  )
}
