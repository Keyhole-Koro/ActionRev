import type { GraphSummary } from '../types/graph-summary'

type GraphSummaryPanelProps = {
  apiBaseUrl: string
  graph: GraphSummary | null
  isLoading: boolean
  error: string | null
}

export function GraphSummaryPanel(props: GraphSummaryPanelProps) {
  const { apiBaseUrl, graph, isLoading, error } = props

  return (
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
          <dd className="mt-2 break-all text-slate-900">{apiBaseUrl}</dd>
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
  )
}
