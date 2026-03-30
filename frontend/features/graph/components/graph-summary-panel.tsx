import type { GraphSummary } from '../types/graph-summary'

type GraphSummaryPanelProps = {
  apiBaseUrl: string
  graph: GraphSummary | null
  isLoading: boolean
  error: string | null
}

export function GraphSummaryPanel({ apiBaseUrl, graph, isLoading, error }: GraphSummaryPanelProps) {
  return (
    <aside className="overflow-hidden rounded-xl border border-slate-200/70 bg-white/70 shadow-sm backdrop-blur-md text-sm">
      {/* Status */}
      <div className="flex items-center justify-between gap-3 border-b border-slate-100/80 px-4 py-2.5">
        <span className="text-xs font-medium text-slate-400">Overview</span>
        <span className={`text-xs font-medium ${isLoading ? 'text-amber-400' : 'text-teal-500'}`}>
          {isLoading ? 'Loading…' : 'Live'}
        </span>
      </div>

      {/* Connection */}
      <div className="px-4 py-3">
        <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-slate-300">
          Connection
        </p>
        <dl className="space-y-2">
          {[
            { label: 'API', value: apiBaseUrl },
            { label: 'Workspace', value: 'ws_demo' },
            { label: 'Document', value: 'doc_demo' },
          ].map(({ label, value }) => (
            <div key={label} className="flex items-start justify-between gap-2">
              <dt className="shrink-0 text-xs text-slate-400">{label}</dt>
              <dd className="truncate text-right font-mono text-xs text-slate-500">{value}</dd>
            </div>
          ))}
        </dl>
      </div>

      {error && (
        <div className="px-4 pb-3">
          <p className="rounded-lg border border-red-200 bg-red-50/80 px-3 py-2 text-xs text-red-500">
            {error}
          </p>
        </div>
      )}

      {graph && (
        <>
          <div className="border-t border-slate-100/80 px-4 py-3">
            <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-slate-300">
              Graph
            </p>
            <div className="space-y-1.5">
              <div className="rounded-lg bg-slate-50/80 px-3 py-2">
                <p className="text-xs text-slate-400">Document</p>
                <p className="mt-0.5 truncate font-mono text-xs text-slate-600">{graph.documentId}</p>
              </div>
              <div className="grid grid-cols-2 gap-1.5">
                <div className="rounded-lg bg-slate-50/80 px-3 py-2">
                  <p className="text-xs text-slate-400">Nodes</p>
                  <p className="mt-0.5 text-lg font-bold text-slate-700">{graph.nodeCount}</p>
                </div>
                <div className="rounded-lg bg-slate-50/80 px-3 py-2">
                  <p className="text-xs text-slate-400">Edges</p>
                  <p className="mt-0.5 text-lg font-bold text-slate-700">{graph.edgeCount}</p>
                </div>
              </div>
            </div>
          </div>

          <div className="border-t border-slate-100/80 px-4 py-3">
            <p className="mb-2 text-xs font-semibold uppercase tracking-widest text-slate-300">
              Labels
            </p>
            <div className="flex flex-wrap gap-1">
              {graph.nodeLabels.map((label) => (
                <span
                  key={label}
                  className="rounded-md border border-slate-200/70 bg-white/60 px-2 py-0.5 text-xs text-slate-500"
                >
                  {label}
                </span>
              ))}
            </div>
          </div>
        </>
      )}
    </aside>
  )
}
