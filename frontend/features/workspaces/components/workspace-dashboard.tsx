import Link from 'next/link'
import type { WorkspaceCard } from '../types/workspace-card'

type WorkspaceDashboardProps = {
  workspaces: WorkspaceCard[]
}

export function WorkspaceDashboard({ workspaces }: WorkspaceDashboardProps) {
  return (
    <div className="min-h-screen bg-slate-50">
      {/* Nav */}
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex h-14 max-w-5xl items-center gap-3 px-6 md:px-10">
          <span className="h-2 w-2 rounded-full bg-teal-500" />
          <span className="text-sm font-semibold text-slate-900">Synthify</span>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-12 md:px-10">
        {/* Page title */}
        <div className="mb-10">
          <h1 className="text-2xl font-bold text-slate-900">Workspaces</h1>
          <p className="mt-1 text-sm text-slate-500">
            ドキュメントから抽出したグラフを確認・探索する
          </p>
        </div>

        {/* Workspace list */}
        <div className="space-y-2">
          {workspaces.map((ws) => (
            <Link
              key={ws.id}
              href={`/workspace/${ws.id}`}
              className="group flex items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white px-5 py-4 shadow-sm transition-colors hover:border-slate-300 hover:bg-slate-50"
            >
              {/* Left */}
              <div className="min-w-0">
                <div className="flex items-center gap-2.5">
                  <h2 className="truncate text-sm font-semibold text-slate-900">{ws.name}</h2>
                  <span
                    className={`shrink-0 rounded-md px-2 py-0.5 text-xs font-medium ${
                      ws.badge === 'Live stub'
                        ? 'bg-teal-50 text-teal-700'
                        : 'bg-slate-100 text-slate-500'
                    }`}
                  >
                    {ws.badge}
                  </span>
                </div>
                <p className="mt-1 truncate text-sm text-slate-500">{ws.summary}</p>
              </div>

              {/* Right */}
              <div className="flex shrink-0 items-center gap-6 text-sm">
                <div className="hidden items-center gap-4 text-slate-400 sm:flex">
                  <span>
                    <span className="font-medium text-slate-700">{ws.graphNodeCount}</span> nodes
                  </span>
                  <span>
                    <span className="font-medium text-slate-700">{ws.graphEdgeCount}</span> edges
                  </span>
                  <span className="font-mono text-xs text-slate-300">{ws.documentId}</span>
                </div>
                <span className="text-slate-300 transition-colors group-hover:text-teal-500">→</span>
              </div>
            </Link>
          ))}
        </div>
      </main>
    </div>
  )
}
