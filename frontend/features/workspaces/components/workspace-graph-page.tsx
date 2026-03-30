'use client'

import Link from 'next/link'
import { GraphCanvasPanel } from '@/features/graph/components/graph-canvas-panel'
import { useGetGraph } from '@/features/graph/hooks/use-get-graph'
import type { WorkspaceCard } from '../types/workspace-card'

type WorkspaceGraphPageProps = {
  workspace: WorkspaceCard
}

export function WorkspaceGraphPage({ workspace }: WorkspaceGraphPageProps) {
  const { graph, isLoading, error } = useGetGraph({
    workspaceId: workspace.id,
    documentId: workspace.documentId,
  })

  return (
    <div className="relative h-screen bg-white">
      {/* Floating nav */}
      <header className="absolute inset-x-0 top-0 z-20">
        <div className="mx-auto flex h-12 max-w-none items-center justify-between gap-4 px-5">
          <nav className="flex items-center gap-2 text-sm">
            <span className="h-1.5 w-1.5 rounded-full bg-teal-500" />
            <Link href="/dashboard" className="text-slate-400 transition-colors hover:text-slate-600">
              Workspaces
            </Link>
            <span className="text-slate-300">/</span>
            <span className="text-slate-600">{workspace.name}</span>
          </nav>
          <Link
            href="/dashboard"
            className="rounded-lg border border-slate-200/70 bg-white/70 px-3 py-1.5 text-xs text-slate-400 shadow-sm backdrop-blur-sm transition-colors hover:text-slate-600"
          >
            ← Back
          </Link>
        </div>
      </header>

      {/* Full-screen canvas */}
      <main className="absolute inset-0">
        <GraphCanvasPanel canvas={graph?.canvas ?? null} isLoading={isLoading} error={error} />

      </main>
    </div>
  )
}
