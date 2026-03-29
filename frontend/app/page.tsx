'use client'

import { HeroPanel } from '@/components/hero-panel'
import { GraphCanvasPanel } from '@/features/graph/components/graph-canvas-panel'
import { GraphSummaryPanel } from '@/features/graph/components/graph-summary-panel'
import { useGetGraph } from '@/features/graph/hooks/use-get-graph'

export default function Page() {
  const { graph, summary, isLoading, error } = useGetGraph()
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'not set'

  return (
    <main className="min-h-screen px-6 py-10 text-ink md:px-10">
      <div className="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <HeroPanel />
        <GraphSummaryPanel apiBaseUrl={apiBaseUrl} graph={summary} isLoading={isLoading} error={error} />
        <GraphCanvasPanel canvas={graph?.canvas ?? null} isLoading={isLoading} error={error} />
      </div>
    </main>
  )
}
