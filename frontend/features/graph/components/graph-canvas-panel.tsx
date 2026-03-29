'use client'

import { Background, Controls, MiniMap, ReactFlow } from '@xyflow/react'
import type { GraphCanvas } from '../types/graph-canvas'

type GraphCanvasPanelProps = {
  canvas: GraphCanvas | null
  isLoading: boolean
  error: string | null
}

export function GraphCanvasPanel(props: GraphCanvasPanelProps) {
  const { canvas, isLoading, error } = props

  return (
    <section className="overflow-hidden rounded-[2rem] bg-slate-950/95 p-3 text-slate-50 shadow-panel ring-1 ring-slate-900/60 lg:col-span-2">
      <div className="flex items-center justify-between gap-4 px-4 py-4 md:px-6">
        <div>
          <p className="text-sm uppercase tracking-[0.32em] text-teal-200">Canvas</p>
          <h2 className="mt-2 text-2xl font-semibold tracking-tight">Document graph</h2>
        </div>
        <div className="rounded-full border border-white/10 bg-white/5 px-4 py-2 text-sm text-slate-200">
          {isLoading ? 'Loading graph' : canvas ? `${canvas.nodes.length} nodes / ${canvas.edges.length} edges` : 'No graph'}
        </div>
      </div>

      {error ? (
        <div className="px-4 pb-4 md:px-6">
          <p className="rounded-2xl border border-rose-400/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-100">
            Graph canvas failed: {error}
          </p>
        </div>
      ) : null}

      <div className="h-[34rem] rounded-[1.5rem] border border-white/10 bg-[radial-gradient(circle_at_top_left,_rgba(45,212,191,0.12),_transparent_24%),linear-gradient(180deg,_rgba(15,23,42,0.85),_rgba(2,6,23,1))]">
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
              nodeBorderRadius={12}
              nodeColor={(node) => String(node.style?.border ?? '#e2e8f0').replace('1px solid ', '')}
              maskColor="rgba(15, 23, 42, 0.45)"
            />
            <Controls showInteractive={false} />
            <Background color="rgba(148, 163, 184, 0.28)" gap={24} size={1.2} />
          </ReactFlow>
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-slate-300">
            {isLoading ? 'Loading graph…' : 'No graph data available.'}
          </div>
        )}
      </div>
    </section>
  )
}
