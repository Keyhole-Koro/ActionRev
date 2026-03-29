'use client'

import { useEffect, useState } from 'react'
import { graphClient } from '../lib/graph-client'

type GraphSummary = {
  documentId: string
  nodeCount: number
  edgeCount: number
  nodeLabels: string[]
}

export default function Page() {
  const [graph, setGraph] = useState<GraphSummary | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    async function load() {
      try {
        setIsLoading(true)
        const response = await graphClient.getGraph({
          workspaceId: 'ws_demo',
          documentId: 'doc_demo',
          categoryFilters: [],
          levelFilters: [],
          edgeTypeFilters: [],
          limit: 50,
          sourceFilename: '',
          resolveAliases: false,
        })

        if (cancelled) {
          return
        }

        setGraph({
          documentId: response.documentId,
          nodeCount: response.graph?.nodes.length ?? 0,
          edgeCount: response.graph?.edges.length ?? 0,
          nodeLabels: response.graph?.nodes.map((node: { label: string }) => node.label) ?? [],
        })
        setError(null)
      } catch (err) {
        if (cancelled) {
          return
        }
        setError(err instanceof Error ? err.message : 'failed to load graph')
      } finally {
        if (!cancelled) {
          setIsLoading(false)
        }
      }
    }

    void load()

    return () => {
      cancelled = true
    }
  }, [])

  return (
    <main className="min-h-screen px-6 py-10 text-ink md:px-10">
      <div className="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[1.15fr_0.85fr]">
        <section className="overflow-hidden rounded-[2rem] bg-slate-950 px-8 py-10 text-slate-50 shadow-panel">
          <p className="text-sm uppercase tracking-[0.35em] text-teal-200">Synthify</p>
          <h1 className="mt-4 max-w-2xl text-4xl font-semibold tracking-tight md:text-6xl">
            Graph extraction now runs on Next.js with a Tailwind-driven shell.
          </h1>
          <p className="mt-6 max-w-xl text-base leading-7 text-slate-300 md:text-lg">
            This page already calls the generated Connect client and renders the current GetGraph stub from the backend.
          </p>
          <div className="mt-8 flex flex-wrap gap-3 text-sm text-slate-200">
            <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Next.js App Router</span>
            <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Tailwind CSS</span>
            <span className="rounded-full border border-white/15 bg-white/5 px-4 py-2">Connect RPC</span>
          </div>
        </section>

        <section className="rounded-[2rem] bg-white/80 p-6 shadow-panel ring-1 ring-slate-200 backdrop-blur">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm font-medium uppercase tracking-[0.28em] text-slate-500">Runtime</p>
              <h2 className="mt-2 text-2xl font-semibold text-slate-900">GetGraph Summary</h2>
            </div>
            <div className="rounded-full bg-accentSoft px-4 py-2 text-sm font-medium text-accent">
              {isLoading ? 'Loading' : 'Connected'}
            </div>
          </div>

          <dl className="mt-6 grid gap-3 rounded-3xl bg-panel p-5 text-sm text-slate-600 sm:grid-cols-3">
            <div>
              <dt className="uppercase tracking-[0.2em] text-slate-400">API</dt>
              <dd className="mt-2 break-all text-slate-900">{process.env.NEXT_PUBLIC_API_BASE_URL ?? 'not set'}</dd>
            </div>
            <div>
              <dt className="uppercase tracking-[0.2em] text-slate-400">Workspace</dt>
              <dd className="mt-2 text-slate-900">ws_demo</dd>
            </div>
            <div>
              <dt className="uppercase tracking-[0.2em] text-slate-400">Document</dt>
              <dd className="mt-2 text-slate-900">doc_demo</dd>
            </div>
          </dl>

          {error ? (
            <p className="mt-6 rounded-2xl border border-rose-200 bg-rose-50 px-4 py-3 text-sm text-rose-700">
              Graph load failed: {error}
            </p>
          ) : null}

          {graph ? (
            <div className="mt-6 space-y-6">
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="rounded-3xl bg-slate-950 px-5 py-4 text-white">
                  <p className="text-xs uppercase tracking-[0.24em] text-slate-400">document_id</p>
                  <p className="mt-3 text-lg font-semibold">{graph.documentId}</p>
                </div>
                <div className="rounded-3xl bg-white px-5 py-4 ring-1 ring-slate-200">
                  <p className="text-xs uppercase tracking-[0.24em] text-slate-400">nodes</p>
                  <p className="mt-3 text-3xl font-semibold text-slate-950">{graph.nodeCount}</p>
                </div>
                <div className="rounded-3xl bg-white px-5 py-4 ring-1 ring-slate-200">
                  <p className="text-xs uppercase tracking-[0.24em] text-slate-400">edges</p>
                  <p className="mt-3 text-3xl font-semibold text-slate-950">{graph.edgeCount}</p>
                </div>
              </div>

              <div>
                <p className="text-xs uppercase tracking-[0.24em] text-slate-400">Node labels</p>
                <div className="mt-4 flex flex-wrap gap-3">
                  {graph.nodeLabels.map((label) => (
                    <span key={label} className="rounded-full bg-slate-900 px-4 py-2 text-sm font-medium text-white">
                      {label}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          ) : null}
        </section>
      </div>
    </main>
  )
}
