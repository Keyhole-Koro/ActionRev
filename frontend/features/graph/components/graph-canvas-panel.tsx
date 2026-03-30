'use client'

import { Background, Controls, MiniMap, ReactFlow } from '@xyflow/react'
import type { GraphCanvas } from '../types/graph-canvas'

type GraphCanvasPanelProps = {
  canvas: GraphCanvas | null
  isLoading: boolean
  error: string | null
  emptyMessage?: string
}

export function GraphCanvasPanel({
  canvas,
  isLoading,
  error,
  emptyMessage = 'No graph data.',
}: GraphCanvasPanelProps) {
  return (
    <div className="absolute inset-0">
      {canvas ? (
        <ReactFlow
          fitView
          nodes={canvas.nodes}
          edges={canvas.edges}
          proOptions={{ hideAttribution: true }}
          defaultEdgeOptions={{ zIndex: 1 }}
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
